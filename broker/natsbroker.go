package broker

import (
	"context"

	"github.com/nats-io/nats.go"
)

type natsBroker struct {
	conn *nats.Conn
	opts *Options
}

func NewNatsBroker(opts *Options) (*natsBroker, error) {
	options := []nats.Option{
		nats.MaxReconnects(opts.MaxReconnects),
		nats.ReconnectWait(opts.ReconnectWait),
	}

	nc, err := nats.Connect(opts.Servers[0], options...)
	if err != nil {
		return nil, err
	}

	return &natsBroker{conn: nc, opts: opts}, nil
}

func (n *natsBroker) Publish(ctx context.Context, subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

func (n *natsBroker) Subscribe(ctx context.Context, subject string, handler func(msg []byte) error) (*Subscription, error) {
	return n.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m.Data)
	})
}

func (n *natsBroker) SubscribeAndReply(ctx context.Context, subject string, handler func(msg *Message)) (*Subscription, error) {
	return n.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m)
	})
}

func (b *natsBroker) Request(ctx context.Context, subject string, msg []byte) (*Message, error) {
	return b.conn.RequestWithContext(ctx, subject, msg)
}

func (n *natsBroker) Drain() error {
	return n.conn.Drain()
}

func (n *natsBroker) Close() {
	n.conn.Close()
}
