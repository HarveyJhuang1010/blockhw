package bo

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsClient interface {
	CreateStream(ctx context.Context, streamName, subjects string) error
	DeleteStream(ctx context.Context, streamName string) error
	Publish(ctx context.Context, topic string, msg interface{}) error
	Subscribe(ctx context.Context, topic string, callBack func(msg *nats.Msg)) (*nats.Subscription, error)
	SubscribeSync(ctx context.Context, topic, durable string) (*nats.Subscription, error)
	PullSubscribe(ctx context.Context, topic, durable string) (*nats.Subscription, error)
	QueueSubscribe(ctx context.Context, topic, queue, durable string, callBack func(msg *nats.Msg)) (*nats.Subscription, error)
	QueueSubscribeSync(ctx context.Context, topic, queue, durable string) (*nats.Subscription, error)
	AddConsumer(ctx context.Context, stream string, cfg *nats.ConsumerConfig) (*nats.ConsumerInfo, error)
	DeleteConsumer(ctx context.Context, stream, consumer string) error
	GetJetStream() nats.JetStreamContext
}

type NatSubscription interface {
	NextMsg(timeout time.Duration) (NatsMsg, error)
	Fetch(batch int) ([]NatsMsg, error)
}

type NatsSub struct {
	*nats.Subscription
}

// NextMsg for push mode
func (s *NatsSub) NextMsg(timeout time.Duration) (NatsMsg, error) {
	msg, err := s.Subscription.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Fetch for pull mode
func (s *NatsSub) Fetch(batch int) ([]NatsMsg, error) {
	msgs, err := s.Subscription.Fetch(batch, nats.MaxWait(5*time.Second))
	if err != nil {
		return nil, err
	}

	ms := make([]NatsMsg, len(msgs))
	for i, msg := range msgs {
		ms[i] = msg
	}
	return ms, nil
}

type NatsMsg interface {
	Ack(opts ...nats.AckOpt) error
	Nak(opts ...nats.AckOpt) error
	InProgress(opts ...nats.AckOpt) error
	NakWithDelay(delay time.Duration, opts ...nats.AckOpt) error
	Term(opts ...nats.AckOpt) error
}
