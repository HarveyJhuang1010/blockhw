package messaging

import (
	"fmt"
	"time"

	"github.com/HarveyJhuang1010/blockhw/internal/config"
	"github.com/HarveyJhuang1010/blockhw/internal/logging"
	"github.com/HarveyJhuang1010/blockhw/internal/models/bo"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var (
	defaultClient *natsClient
	log           = logging.NewZapLogger("mq")
)

type natsClient struct {
	cfg *config.NatsConfig

	natsConnect      *nats.Conn
	streamingConnect nats.JetStreamContext
}

func Initialize(cfg *config.NatsConfig) error {
	defaultClient = &natsClient{
		cfg: cfg,
	}

	if err := defaultClient.connectNats(); err != nil {
		return err
	}
	return defaultClient.connectStreaming()
}

func Finalize() error {
	defaultClient.natsConnect.Close()
	return nil
}

func GetClient() bo.NatsClient {
	return defaultClient
}

func (n *natsClient) connectNats() (err error) {

	opts := []nats.Option{
		nats.ReconnectHandler(func(conn *nats.Conn) {
			log.Error(fmt.Sprintf("[NATS] Reconnection count: %d", conn.Reconnects))
		}),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				log.Error("[Nats] unable to connect to nats. ", zap.Error(err))
			}
		}),
		nats.ErrorHandler(n.logError),
		nats.Timeout(10 * time.Second),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Warn("[nats] client reconnected")
		}),
	}

	n.natsConnect, err = nats.Connect(n.cfg.GetURL(), opts...)
	if err != nil {
		return err
	}

	log.Info("[nats] connect success")
	return nil
}

func (n *natsClient) connectStreaming() (err error) {
	jsCtx, err := n.natsConnect.JetStream()
	if err != nil {
		return err
	}

	n.streamingConnect = jsCtx
	info, err := jsCtx.AccountInfo()
	if err != nil {
		return err
	}

	log.Info("[nats]", zap.Any("info", info))

	return nil
}

func (n *natsClient) logError(conn *nats.Conn, subscription *nats.Subscription, err error) {
	log.Warn(
		fmt.Sprintf(
			"A problem happens on the topic: %s | queue name: %s | error: %s",
			subscription.Subject,
			subscription.Queue,
			err.Error(),
		),
	)
}
