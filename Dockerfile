ARG GOLANG_VERSION=1.18
ARG ALPINE_VERSION=3

FROM golang:${GOLANG_VERSION}-alpine AS build
WORKDIR /build
RUN apk add --update git

COPY . .
RUN go mod download -x

RUN go build -o client cmd/client/main.go
RUN go build -o server cmd/server/main.go

FROM alpine:${ALPINE_VERSION}
COPY --from=build /build/client /client
COPY --from=build /build/server /server
ENTRYPOINT ["/server"]
