package listener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HarveyJhuang1010/blockhw/internal/appcontext"
	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

type workflowListener struct {
	in digIn

	queue *semaphore.Weighted
}

func newBlockListener(in digIn) bo.Listener {
	return &workflowListener{
		in:    in,
		queue: semaphore.NewWeighted(config.GetConfig().Worker.MaxWorkers),
	}
}

func (w *workflowListener) Listen(ctx context.Context) {
	logger := appcontext.GetLogger()

	logger.Info("start create stream")
	streamName := "block"
	subject := "block.*"
	err := w.in.MQ.CreateStream(ctx, streamName, subject)
	if err != nil {
		logger.Error("failed to create stream", zap.Error(err))
		panic(err)
	}

	logger.Info("subscribe workflow event")
	sub, err := w.PullSubscribeWorkflowEvent(ctx, streamName, subject, "block-syncer")
	if err != nil {
		logger.Error("`ailed to subscribe workflow event", zap.Error(err))
		panic(err)
	}

	logger.Info("start handle events")
	for {
		select {
		case <-ctx.Done():
			logger.Warn("server shutdown, workflow listener exit")
			return

		default:
			msgs, err := sub.Fetch(1)
			if errors.Is(err, nats.ErrTimeout) {
				continue
			}
			if err != nil {
				logger.Error("failed to fetch workflow event", zap.Error(err))
				panic(err)
			}
			msg := msgs[0]

			if ok := w.queue.TryAcquire(1); !ok {
				logger.Warn("queue is full, failed to acquire semaphore, redeliver message later")
				delayTime := time.Duration(config.GetConfig().Worker.DelayMinute) * time.Minute
				msg.NakWithDelay(delayTime)
				continue
			}

			go func(m bo.NatsMsg) {
				defer w.queue.Release(1)
				w.handleMessage(ctx, m)
			}(msg)
		}
	}
}

func (w *workflowListener) PullSubscribeWorkflowEvent(ctx context.Context, streamName, subject, durableName string) (bo.NatSubscription, error) {
	logger := appcontext.GetLogger()

	var sub *nats.Subscription
	consumerInfo, err := w.in.MQ.GetJetStream().ConsumerInfo(streamName, durableName)
	if err != nil {
		if !errors.Is(err, nats.ErrConsumerNotFound) {
			return nil, errors.Wrap(nats.ErrConsumerNotFound, "failed to get consumer info")
		}
		sub, err = w.in.MQ.PullSubscribe(ctx, subject, durableName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to pull subscribe")
		}
	} else {
		cfg := consumerInfo.Config
		cfg.MaxDeliver = 5
		_, err := w.in.MQ.GetJetStream().UpdateConsumer(streamName, &cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to update consumer")
		}
		sub, err = w.in.MQ.PullSubscribe(ctx, subject, durableName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to pull subscribe")
		}
	}

	logger.Info("pull subscribe workflow event", zap.String("subject", subject))
	return &bo.NatsSub{Subscription: sub}, nil
}

func (w *workflowListener) handleMessage(ctx context.Context, msg bo.NatsMsg) {
	logger := appcontext.GetLogger()

	rawMsg, ok := msg.(*nats.Msg)
	if !ok || rawMsg == nil {
		logger.Error("invalid msg type")
		rawMsg.Term()
		return
	}
	logger = logger.With(zap.String("subject", rawMsg.Subject))
	logger.Info("receive msg", zap.ByteString("data", rawMsg.Data))

	var blockNum uint64
	if err := json.Unmarshal(rawMsg.Data, &blockNum); err != nil {
		logger.Error("failed to decode event", zap.Error(err))
		rawMsg.Term() // tell nats server not to redeliver this message or msg will be redelivered infinitely
		return
	}

	if err := w.in.Block.SyncBlockByNum(ctx, blockNum); err != nil {
		logger.Error("failed to sync block", zap.Error(err))
		rawMsg.Nak()
	} else {
		rawMsg.Ack()
	}
}
