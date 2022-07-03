package client

import (
	"context"
	"errors"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yosadchyi/go-client-server/pkg/message"
	"github.com/yosadchyi/go-client-server/pkg/util"
)

var (
	IntegerValueExpected = errors.New("integer value expected")
	UnknownCommand       = errors.New("unknown command")
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
		msg = message.NewAdd(data)
	case '-':
		if value, err := strconv.Atoi(data); err != nil {
			return IntegerValueExpected
		} else {
			msg = message.NewRemove(value)
		}
	case '<':
		if value, err := strconv.Atoi(data); err != nil {
			return IntegerValueExpected
		} else {
			msg = message.NewGet(value)
		}
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
