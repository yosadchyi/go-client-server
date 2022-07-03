package client

import "context"

// Processor processes stream of incoming commands
type Processor struct {
	lines     chan string
	executor  *Executor
	responder Responder
}

// NewProcessor creates new Processor
func NewProcessor(lines chan string, executor *Executor, responder Responder) *Processor {
	return &Processor{
		lines:     lines,
		executor:  executor,
		responder: responder,
	}
}

// Run starts processing loop
func (p *Processor) Run(ctx context.Context) {
	p.responder.Help()

	for {
		select {
		case <-ctx.Done():
			p.responder.Bye()
			return

		case line := <-p.lines:
			if line == "EOF" {
				return
			}
			switch err := p.executor.ExecuteCmd(ctx, line); err {
			case UnknownCommand:
				p.responder.Error(err)
				p.responder.Help()
			case KeyValueExpected:
				p.responder.Error(err)
				p.responder.Help()
			case nil:
				p.responder.Ok()
			default:
				p.responder.Error(err)
			}
		}
	}
}
