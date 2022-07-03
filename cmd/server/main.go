package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yosadchyi/go-client-server/pkg/message"
	"github.com/yosadchyi/go-client-server/pkg/server"
	"github.com/yosadchyi/go-client-server/pkg/util"
)

func main() {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")
	resolver := util.LocalResolver(awsEndpoint, awsRegion)

	parallelismDegree := flag.Int(
		"paralellism-degree",
		runtime.NumCPU(),
		"number of processors to be run concurrently, by default equal to system's number of CPU",
	)
	queueUrl := flag.String(
		"queue-url",
		os.Getenv("QUEUE_URL"),
		"SQS queue",
	)
	waitTimeSeconds := flag.Int(
		"wait-time-seconds",
		1,
		"number of seconds to wait for SQS messages, bigger value decreases CPU load",
	)
	flag.Parse()

	ctx, cancelFn := context.WithCancel(context.Background())

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
	)
	if err != nil {
		log.Fatalf("failed to load default config %s", err)
	}

	sqsSvc := sqs.NewFromConfig(cfg)
	messages := make(chan *message.Any, 128)
	reader := server.NewReader(sqsSvc, *queueUrl, messages)
	storage := server.NewRWLockedStorage(server.NewMemoryStorage())
	processor := server.NewProcessor(messages)

	for i := 1; i <= *parallelismDegree; i++ {
		go processor.Run(ctx, server.NewProcessFn(i, storage))
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		s := <-sig
		log.Printf("system signal: %+v", s)
		cancelFn()
	}()

	log.Printf("waiting for messages on %s...", *queueUrl)

	reader.Run(ctx, int32(*waitTimeSeconds))
}
