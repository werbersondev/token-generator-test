package pubsubx

import (
	"context"

	"cloud.google.com/go/pubsub"
)

// CreateTopicIfNotExists checks if a Pub/Sub topic exists and creates it if it does not.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the operation.
//   - client: The Pub/Sub client used to interact with the Pub/Sub service.
//   - topicName: The name of the topic to check or create.
//
// Returns:
//   - *pubsub.Topic: A reference to the existing or newly created topic.
//   - error: Any error encountered during the operation.
func CreateTopicIfNotExists(ctx context.Context, client *pubsub.Client, topicName string) (*pubsub.Topic, error) {
	topic := client.Topic(topicName)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return topic, nil
	}
	return client.CreateTopic(ctx, topicName)
}

// CreateSubscriptionIfNotExists checks if a Pub/Sub subscription exists and creates it if it does not.
//
// Parameters:
//   - ctx: The context for managing the lifecycle of the operation.
//   - client: The Pub/Sub client used to interact with the Pub/Sub service.
//   - topic: The Pub/Sub topic to which the subscription will be associated.
//   - subscriptionName: The name of the subscription to check or create.
//
// Returns:
//   - *pubsub.Subscription: A reference to the existing or newly created subscription.
//   - error: Any error encountered during the operation.
func CreateSubscriptionIfNotExists(ctx context.Context, client *pubsub.Client, topic *pubsub.Topic, subscriptionName string) (*pubsub.Subscription, error) {
	sub := client.Subscription(subscriptionName)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return sub, nil
	}
	return client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{Topic: topic})
}
