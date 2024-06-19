package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ardanlabs/conf/v3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/werbersondev/token-generator-test/cmd/httpservice/api"
	"github.com/werbersondev/token-generator-test/domain/service"
	"github.com/werbersondev/token-generator-test/extensions/httpx"
	"github.com/werbersondev/token-generator-test/extensions/loggerx"
	"github.com/werbersondev/token-generator-test/extensions/pubsubx"
	pubsubgw "github.com/werbersondev/token-generator-test/gateway/pubsub"
)

type config struct {
	ServerAddr         string        `conf:"env:SERVER_ADDR,default:0.0.0.0:3000"`
	ServerReadTimeout  time.Duration `conf:"env:SERVER_READ_TIMEOUT,default:30s"`
	ServerWriteTimeout time.Duration `conf:"env:SERVER_WRITE_TIMEOUT,default:30s"`

	PubSubHost             string `conf:"env:PUBSUB_EMULATOR_HOST,required"`
	ProjectID              string `conf:"env:GCP_PROJECT_ID,default:my_project_key"`
	TokenGenerationTopicID string `conf:"env:GCP_TOKEN_GENERATOR_TOPIC,default:token_generation_topic"`
}

func main() {
	logger := loggerx.NewDevelopment()
	zerolog.DefaultContextLogger = &logger

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	if err := runServer(ctx); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Send()
	}
}

func runServer(ctx context.Context) error {
	var cfg config
	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return fmt.Errorf("usage: httpservice [config]\n%s", help)
		}

		return fmt.Errorf("error parsing the configuration: %w", err)
	}

	client, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("close pubsub client")
		}
	}()

	topic, err := pubsubx.CreateTopicIfNotExists(ctx, client, cfg.TokenGenerationTopicID)
	if err != nil {
		return fmt.Errorf("creating topic %s: %w", cfg.TokenGenerationTopicID, err)
	}

	publisher := pubsubgw.NewRequestTokenGenerationPublisher(topic)

	tokenService := service.NewRequestTokenGenerationService(publisher)

	server := createServer(tokenService, cfg)

	httpx.Run(ctx, &server)

	return nil
}

func createServer(tokenService *service.RequestTokenGenerationService, cfg config) http.Server {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	apiV1 := api.New(tokenService)
	apiV1.Routes(router)

	return http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      router,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}
}
