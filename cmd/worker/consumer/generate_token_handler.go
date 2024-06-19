package consumer

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog/log"

	"github.com/werbersondev/token-generator-test/domain/model"
)

func (c *GenerateTokenConsumer) GenerateTokenHandler(ctx context.Context, msg *pubsub.Message) {
	defer msg.Ack()

	var request model.TokenGenerationRequest
	if err := json.Unmarshal(msg.Data, &request); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("failed to unmarshal message")
		return
	}

	token, err := c.useCase.GenerateToken(ctx, request)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to generate token")
		return
	}

	log.Ctx(ctx).Info().Msgf("Token generated for project: %s token: %s", request.ProjectID, token)
}
