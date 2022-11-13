package util

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type MsgHandler func(ctx context.Context, data []byte) error

type NatsJetstream interface {
	CreateStream(stream, subject string) error
	PublishAsyncContext(ctx context.Context, subject string, data []byte, opts ...nats.PubOpt) (string, error)
	QueueSubscribe(subject, queue string, cb MsgHandler, opts ...nats.SubOpt) error
	Close()
}

// TODO: update nats conf
type NatsConf struct {
	Url           string        `mapstructure:"URL"`
	UserName      string        `mapstructure:"USERNAME"`
	Password      string        `mapstructure:"PASSWORD"`
	MaxReconnect  int           `mapstructure:"MAX_RECONNECT"`
	ReconnectWait time.Duration `mapstructure:"RECONNECT_WAIT"`
}

type NatsJSImpl struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	subs   []*nats.Subscription
	logger *zap.Logger
}

func ConnectNats(logger *zap.Logger, conf NatsConf) (*NatsJSImpl, error) {
	conn, err := nats.Connect(
		conf.Url,
		nats.UserInfo(conf.UserName, conf.Password),
		nats.ReconnectWait(conf.ReconnectWait),
		nats.MaxReconnects(conf.MaxReconnect),
	)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Cannot connect to nats: %s", err.Error()))
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Cannot get jetstream: %s", err.Error()))
	}

	natsjs := &NatsJSImpl{
		conn:   conn,
		js:     js,
		logger: logger,
	}

	return natsjs, nil
}

func (n *NatsJSImpl) CreateStream(stream, subject string) error {
	info, _ := n.js.StreamInfo(stream)
	if info != nil {
		n.logger.Info("CreateStream: stream is exist")
		return nil
	}

	cfg := &nats.StreamConfig{
		Name:      stream,
		Retention: nats.InterestPolicy,
		Replicas:  1,
		Subjects:  []string{subject},
	}

	_, err := n.js.AddStream(cfg)
	return err
}

func (n *NatsJSImpl) PublishAsyncContext(ctx context.Context, subject string, data []byte, opts ...nats.PubOpt) (string, error) {
	opts = append(opts, nats.MsgId(fmt.Sprintf("%s", uuid.New())))
	pubAck, err := n.js.PublishAsync(subject, data, opts...)
	if err != nil {
		return "", err
	}

	return pubAck.Msg().Header.Get("Nats-Msg-Id"), err
}

func (n *NatsJSImpl) QueueSubscribe(subject, queue string, cb MsgHandler, opts ...nats.SubOpt) error {
	sub, err := n.js.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		var sequence uint64
		metadata, err := msg.Metadata()
		if err == nil {
			sequence = metadata.Sequence.Stream
		}

		msgID := msg.Header.Get("Nats-Msg-Id")

		if err := cb(ctx, msg.Data); err != nil {
			n.logger.Error("QueueSubscribe: Handle msg occur issue", zap.Uint64("sequence_id", sequence), zap.String("msg_id", msgID), zap.Error(err))
			return
		}

		if err := msg.Ack(); err != nil {
			n.logger.Error("QueueSubscribe: ack msg occur issue", zap.Uint64("sequence_id", sequence), zap.String("msg_id", msgID), zap.Error(err))
			return
		}

		n.logger.Info("QueueSubscribe: Ack msg success", zap.Uint64("sequence_id", sequence), zap.String("msg_id", msgID))
	}, opts...)

	n.subs = append(n.subs, sub)
	return err
}

func (n *NatsJSImpl) Close() {
	for _, s := range n.subs {
		info, err := s.ConsumerInfo()
		if err != nil {
			return
		}
		n.logger.Info(fmt.Sprintf("Close: Draining consumer: %s", info.Name))
		err = s.Unsubscribe()
		if err != nil {
			n.logger.Error("Close: Unsubscribe consumer occur problem", zap.Error(err))
		}

		err = s.Drain()
		if err != nil {
			n.logger.Error("Close: Draining consumer occur problem", zap.Error(err))
		}
	}

	n.conn.Close()
}
