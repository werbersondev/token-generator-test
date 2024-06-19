package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"

	"github.com/werbersondev/token-generator-test/domain/model"
)

type RequestTokenGenerationPublisher struct {
	topic *pubsub.Topic
}

func NewRequestTokenGenerationPublisher(topic *pubsub.Topic) *RequestTokenGenerationPublisher {
	return &RequestTokenGenerationPublisher{
		topic: topic,
	}
}

func (r *RequestTokenGenerationPublisher) PublishRequestTokenGeneration(ctx context.Context, request model.TokenGenerationRequest) error {
	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshalling request data: %w", err)
	}

	result := r.topic.Publish(ctx, &pubsub.Message{
		Data: data,
	})

	if _, err := result.Get(ctx); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}
	return nil
}
