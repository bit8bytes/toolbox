package broker

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
)

type Options struct {
	Servers       []string
	MaxReconnects int
	ReconnectWait time.Duration
}

type Broker interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Subscribe(ctx context.Context, subject string, handler func(msg []byte) error) (*Subscription, error)
	SubscribeAndReply(ctx context.Context, subject string, handler func(msg *Message)) (*Subscription, error)
	Request(ctx context.Context, subject string, msg []byte) (*Message, error)
	Drain() error
	Close()
}

type Subscription = nats.Subscription

type Message = nats.Msg
