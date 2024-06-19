package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/ardanlabs/conf/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/werbersondev/token-generator-test/cmd/worker/consumer"
	"github.com/werbersondev/token-generator-test/domain/service"
	"github.com/werbersondev/token-generator-test/extensions/loggerx"
	"github.com/werbersondev/token-generator-test/extensions/pubsubx"
	"github.com/werbersondev/token-generator-test/gateway/sonarclient"
)

type config struct {
	PubSubHost                    string        `conf:"env:PUBSUB_EMULATOR_HOST,default:localhost:8085"`
	ProjectID                     string        `conf:"env:GCP_PROJECT_ID,default:my_project_key"`
	TokenGenerationTopicID        string        `conf:"env:GCP_TOKEN_GENERATOR_TOPIC,default:token_generation_topic"`
	TokenGenerationSubscriptionID string        `conf:"env:GCP_TOKEN_GENERATOR_SUBSCRIPTION,default:token_generation_subscription"`
	SonarAPIAddress               string        `conf:"env:SONAR_API_ADDRESS,default:http://localhost:9000"`
	SonarAPITimeout               time.Duration `conf:"env:SONAR_API_TIMEOUT,default:30s"`
	SonarAuthToken                string        `conf:"env:SONAR_AUTH_TOKEN,required"`
}

func main() {
	logger := loggerx.NewDevelopment()
	zerolog.DefaultContextLogger = &logger

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	if err := runConsumer(ctx); err != nil {
		log.Ctx(ctx).Fatal().Err(err).Send()
	}
}

func runConsumer(ctx context.Context) error {
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

	subs, err := pubsubx.CreateSubscriptionIfNotExists(ctx, client, topic, cfg.TokenGenerationSubscriptionID)
	if err != nil {
		return fmt.Errorf("creating subscription %s: %w", cfg.TokenGenerationSubscriptionID, err)
	}

	httpClient := sonarclient.New(sonarclient.Config{
		Timeout:   cfg.SonarAPITimeout,
		BaseURL:   cfg.SonarAPIAddress,
		AuthToken: cfg.SonarAuthToken,
	})

	tokenService := service.NewTokenGenerationService(struct {
		*sonarclient.HTTPClient
	}{
		httpClient,
	})

	tokenGeneratorConsumer := consumer.NewGenerateTokenConsumer(subs, tokenService)

	go func() {
		log.Ctx(ctx).Info().Str("project_id", cfg.ProjectID).
			Str("topic", cfg.TokenGenerationTopicID).
			Str("subscription", cfg.TokenGenerationSubscriptionID).
			Msg("consumer started")
		if err := tokenGeneratorConsumer.Start(ctx); err != nil {
			return
		}
	}()

	// Setup signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-sigs
	log.Ctx(ctx).Info().Str("signal", sig.String()).Msg("Received signal")

	ctxStop, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	// Stop the consumer gracefully
	if err := tokenGeneratorConsumer.Stop(ctxStop); err != nil {
		return err
	}
	log.Ctx(ctx).Info().Msg("Consumer stopped gracefully")

	return nil
}
