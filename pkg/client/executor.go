package client

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yosadchyi/go-client-server/pkg/message"
	"github.com/yosadchyi/go-client-server/pkg/util"
)

var (
	KeyValueExpected = errors.New("key/value expected")
	UnknownCommand   = errors.New("unknown command")
)

// Executor is and executor of the commands, provided as text
type Executor struct {
	inputFile *os.File
	sqsClient *sqs.Client
	queueUrl  string
}

// NewExecutor creates new executor
func NewExecutor(inputFile *os.File, sqsClient *sqs.Client, queueUrl string) *Executor {
	return &Executor{
		inputFile: inputFile,
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
	}
}

// ExecuteCmd executes command
func (e *Executor) ExecuteCmd(ctx context.Context, line string) error {
	var msg util.JSONEr

	cmd := line[0]
	data := line[1:]

	switch cmd {
	case '+':
		if idx := strings.Index(data, ":"); idx >= 0 {
			msg = message.NewAdd(data[:idx], data[idx+1:])
		} else {
			return KeyValueExpected
		}
	case '-':
		msg = message.NewRemove(data)
	case '<':
		msg = message.NewGet(data)
	case '*':
		msg = message.NewGetAll()
	}

	if msg == nil {
		return UnknownCommand
	}

	_, err := e.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(e.queueUrl),
		MessageBody: msg.ToJSON(),
	})

	if err != nil {
		return err
	}

	return nil
}
