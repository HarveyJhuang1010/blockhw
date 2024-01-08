package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func (n *natsClient) Publish(ctx context.Context, topic string, msg interface{}) error {
	logger := appcontext.GetLogger()
	var msgByte []byte
	var err error

	if _, ok := msg.([]byte); ok {
		msgByte = msg.([]byte)
	} else {
		msgByte, err = json.Marshal(msg)
		if err != nil {
			logger.Error("error marshal msg", zap.Error(err))
			return err
		}
	}

	_, err = n.streamingConnect.Publish(topic, msgByte)

	return err
}

func (n *natsClient) Subscribe(ctx context.Context, topic string, callBack func(msg *nats.Msg)) (*nats.Subscription, error) {
	return n.streamingConnect.Subscribe(topic, callBack)
}

func (n *natsClient) SubscribeSync(ctx context.Context, topic, durable string) (*nats.Subscription, error) {
	return n.streamingConnect.SubscribeSync(topic, nats.Durable(durable))
}

func (n *natsClient) PullSubscribe(ctx context.Context, topic, durable string) (*nats.Subscription, error) {
	return n.streamingConnect.PullSubscribe(topic, durable, nats.DeliverLast(), nats.AckWait(5*time.Minute))
}

func (n *natsClient) CreateStream(ctx context.Context, streamName, subjects string) error {
	logger := appcontext.GetLogger()
	// Check if the stream already exists; if not, create it.
	stream, err := n.streamingConnect.StreamInfo(streamName)
	if err != nil && err != nats.ErrStreamNotFound {
		logger.Error("error getting stream info: %v", zap.Error(err))
		return err
	}
	if stream == nil {
		logger.Info("creating streams", zap.String("streamName", streamName), zap.String("subjects", subjects))
		_, err = n.streamingConnect.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{subjects},
		})
		if err != nil {
			logger.Error("error creating stream", zap.Error(err))
			return err
		}
	}
	return nil
}

func (n *natsClient) DeleteStream(ctx context.Context, streamName string) error {
	// Check if the ORDERS stream already exists; if not, create it.
	logger := appcontext.GetLogger()
	if _, err := n.streamingConnect.StreamInfo(streamName); err != nil {
		if err == nats.ErrStreamNotFound {
			logger.Info("stream not found not need to delete", zap.String("streamName", streamName))
			return nil
		}
		logger.Error("error getting stream info: %v", zap.Error(err))
		return err
	}

	if err := n.streamingConnect.DeleteStream(streamName); err != nil {
		logger.Error("delete streams", zap.String("streamName", streamName), zap.Error(err))
		return err
	}
	logger.Info("delete streams", zap.String("streamName", streamName))
	return nil
}

func (n *natsClient) QueueSubscribe(ctx context.Context, topic, queue, durable string, callBack func(msg *nats.Msg)) (*nats.Subscription, error) {
	return n.streamingConnect.QueueSubscribe(topic, queue, callBack, nats.Durable(durable))
}

func (n *natsClient) QueueSubscribeSync(ctx context.Context, topic, queue, durable string) (*nats.Subscription, error) {
	return n.streamingConnect.QueueSubscribeSync(topic, queue, nats.Durable(durable), nats.AckExplicit())
}

func (n *natsClient) AddConsumer(ctx context.Context, stream string, cfg *nats.ConsumerConfig) (*nats.ConsumerInfo, error) {
	return n.streamingConnect.AddConsumer(stream, cfg)
}

func (n *natsClient) DeleteConsumer(ctx context.Context, stream, consumer string) error {
	return n.streamingConnect.DeleteConsumer(stream, consumer)
}

func (n *natsClient) GetJetStream() nats.JetStreamContext {
	return n.streamingConnect
}
