package consumer

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"

	"github.com/werbersondev/token-generator-test/domain/model"
)

type GenerateTokenUseCase interface {
	GenerateToken(ctx context.Context, request model.TokenGenerationRequest) (string, error)
}

type GenerateTokenConsumer struct {
	topicSubscription *pubsub.Subscription
	useCase           GenerateTokenUseCase
	startCh, stopCh   chan struct{}
}

func NewGenerateTokenConsumer(topicSubscription *pubsub.Subscription, uc GenerateTokenUseCase) *GenerateTokenConsumer {
	return &GenerateTokenConsumer{
		topicSubscription: topicSubscription,
		useCase:           uc,
		startCh:           make(chan struct{}),
		stopCh:            make(chan struct{}),
	}
}

// Start begins consuming messages from the Pub/Sub subscription and processing them.
func (c *GenerateTokenConsumer) Start(ctx context.Context) error {
	defer close(c.startCh)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c.stopCh
		cancel()
	}()

	if err := c.topicSubscription.Receive(ctx, c.GenerateTokenHandler); err != nil && !errors.Is(err, context.Canceled) {
		log.Ctx(ctx).Error().Err(err).Msg("Error receiving messages")
		return err
	}

	return nil
}

// Stop initiates the graceful shutdown of the consumer. It signals the consumer to stop
// polling for new messages and waits for the ongoing transaction to finish.
func (c *GenerateTokenConsumer) Stop(ctx context.Context) error {
	// Stop polling for new messages.
	close(c.stopCh)

	// Wait for all processing to finish or ctx done.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.startCh:
		return nil
	}
}
