# go-client-server
Client-Service with SQS in a middle written in GO

## Prerequisites

- docker / docker-compose
- go 1.18 to run outside of docker
- Optionally you can use direnv https://direnv.net

## How to run

### With docker-compose

```shell
make up
```

This command spins up localstack, server and executes client with input from ```test/data.txt```.

### Locally

To run the project locally (both client and server) you need to define following environment variables:

```shell
AWS_ENDPOINT=http://localhost:4566
AWS_REGION=eu-central-1
QUEUE_URL=http://localhost:4566/000000000000/queue
```

If you're using direnv, you need to approve contents of .env placed within project:

```shell
direnv allow
```

After setting up you can run go server & client as usually.

Server command line flags:
```text
Usage of ./server:
  -paralellism-degree int
        number of processors to be run concurrently, by default equal to system's number of CPU (default 6)
  -queue-url string
        SQS queue
  -wait-time-seconds int
        number of seconds to wait for SQS messages, bigger value decreases CPU load (default 1)
```

Client command line flags:
```text
Usage of ./client:
  -input-file string
        input file to read commands from, otherwise stdin will be used
  -queue-url string
        SQS queue
```

## Syntax of client input lines

```text
    +ITEM
            add item with data 'ITEM'
    -INDEX
            remove item with index INDEX, where index is an integer number
    <INDEX
            get item with index INDEX, where index is an integer number
    *
            list all items
    EOF or ^C
            quit
```
