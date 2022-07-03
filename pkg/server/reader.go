package server

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yosadchyi/go-client-server/pkg/message"
)

// Reader is responsible for reading data from SQS and passing it to messages channel
type Reader struct {
	sqsClient *sqs.Client
	queueUrl  string
	messages  MessageChan
}

// NewReader creates new reader
func NewReader(sqsClient *sqs.Client, queueUrl string, messages MessageChan) *Reader {
	return &Reader{
		sqsClient: sqsClient,
		queueUrl:  queueUrl,
		messages:  messages,
	}
}

// Run runs reading, can be stopped with context's cancel function
func (s *Reader) Run(ctx context.Context, waitTimeSeconds int32) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping SQS reader")
			return
		default:
			s.receiveMessages(waitTimeSeconds)
		}
	}
}

func (s *Reader) receiveMessages(waitTimeSeconds int32) {
	out, err := s.sqsClient.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueUrl),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     waitTimeSeconds,
	})
	if err != nil {
		log.Printf("error receiving message %s", err)
	}

	for _, m := range out.Messages {
		if m.Body == nil {
			log.Print("received message with empty body")
			continue
		}
		msg, err := message.AnyFromJSON(*m.Body)
		if err != nil {
			log.Printf("error parsing message: %s", err)
			continue
		}
		s.messages <- msg

		_, err = s.sqsClient.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(s.queueUrl),
			ReceiptHandle: m.ReceiptHandle,
		})
		if err != nil {
			log.Printf("error deleting message: %s", err)
			continue
		}
	}
}
