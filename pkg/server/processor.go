package server

import (
	"context"
	"log"

	"github.com/yosadchyi/go-client-server/pkg/message"
)

// Processor allows to process incoming messages with given processing function
type Processor struct {
	messages MessageChan
}

// NewProcessor creates new processor
func NewProcessor(messages MessageChan) *Processor {
	return &Processor{
		messages: messages,
	}
}

// Run runs processing, can be stopped with context's cancel function
func (s *Processor) Run(ctx context.Context, processFn func(*message.Any)) {
	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down processor")
			return
		case msg := <-s.messages:
			processFn(msg)
		}
	}
}
