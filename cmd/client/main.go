package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yosadchyi/go-client-server/pkg/client"
	"github.com/yosadchyi/go-client-server/pkg/util"
)

func main() {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")
	resolver := util.LocalResolver(awsEndpoint, awsRegion)

	queueUrl := flag.String(
		"queue-url",
		os.Getenv("QUEUE_URL"),
		"SQS queue",
	)
	inputFile := flag.String(
		"input-file",
		"",
		"input file to read commands from, otherwise stdin will be used",
	)
	flag.Parse()

	ctx, cancelFn := context.WithCancel(context.Background())

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
	)
	if err != nil {
		log.Fatalf("failed to load default config %e", err)
	}

	svc := sqs.NewFromConfig(cfg)
	file := os.Stdin

	if *inputFile != "" {
		file, err = os.OpenFile(*inputFile, os.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("can't open input file %s", *inputFile)
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	go func() {
		s := <-sig
		fmt.Printf("system signal: %+v\n", s)
		cancelFn()
	}()

	lines := make(chan string, 1)
	isInteractive := file == os.Stdin

	go func() {
		s := bufio.NewScanner(file)
		for s.Scan() {
			lines <- s.Text()
		}
	}()

	executor := client.NewExecutor(file, svc, *queueUrl)
	var responder client.Responder

	if isInteractive {
		responder = client.NewInteractiveResponder()
	} else {
		responder = client.NewBatchResponder()
	}

	processor := client.NewProcessor(lines, executor, responder)

	processor.Run(ctx)
}
