package broker

import "context"

type MuxBroker struct {
	handlers map[string]func(msg *Message)
	broker   Broker
}

func NewServeMux(b Broker) *MuxBroker {
	return &MuxBroker{
		handlers: make(map[string]func(msg *Message)),
		broker:   b,
	}
}

func (mux *MuxBroker) HandleFunc(pattern string, handler func(msg *Message)) {
	mux.handlers[pattern] = handler
}

func (mux *MuxBroker) Subscribe() error {
	ctx := context.Background()
	for subject, handler := range mux.handlers {
		_, err := mux.broker.Subscribe(ctx, subject, func(msg []byte) error {
			handler(&Message{
				Subject: subject,
				Data:    msg,
			})
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
