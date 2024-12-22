package broker

import (
	"context"

	"github.com/nats-io/nats.go"
)

type ConcreteNatsBroker struct {
	conn *nats.Conn
	opts *Options
}

func NewNatsBroker(opts *Options) (*ConcreteNatsBroker, error) {
	options := []nats.Option{
		nats.MaxReconnects(opts.MaxReconnects),
		nats.ReconnectWait(opts.ReconnectWait),
	}

	nc, err := nats.Connect(opts.Servers[0], options...)
	if err != nil {
		return nil, err
	}

	return &ConcreteNatsBroker{conn: nc, opts: opts}, nil
}

func (n *ConcreteNatsBroker) Publish(ctx context.Context, subject string, data []byte) error {
	return n.conn.Publish(subject, data)
}

func (n *ConcreteNatsBroker) Subscribe(ctx context.Context, subject string, handler func(msg *Message)) (*Subscription, error) {
	return n.conn.Subscribe(subject, func(m *nats.Msg) {
		handler(m)
	})
}

func (b *ConcreteNatsBroker) Request(ctx context.Context, subject string, msg []byte) (*Message, error) {
	return b.conn.RequestWithContext(ctx, subject, msg)
}

func (n *ConcreteNatsBroker) Drain() error {
	return n.conn.Drain()
}

func (n *ConcreteNatsBroker) Close() {
	n.conn.Close()
}
