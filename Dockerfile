FROM golang:alpine AS build

RUN apk update && apk add --no-cache git ca-certificates

RUN mkdir -p /tracer
WORKDIR /tracer

COPY ./go.mod ./go.sum ./
RUN GOPROXY=https://proxy.golang.org go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags="-w -s" \
	-installsuffix "static" \
	-o /go/bin/tracer /tracer/cmd/tracer

FROM alpine:latest

ENV ENVIRONMENT = "PROD"
ENV DATABASE_URL "localhost:9080"
ENV BROKER_URL "localhost:5672"
ENV BROKER_USER "guest"
ENV BROKER_PASSWORD "guest"
ENV BATCH_SIZE_LIMIT "1000"
ENV BATCH_TIMEOUT "100"

COPY --from=build /go/bin/tracer /go/bin/tracer

CMD "/go/bin/tracer"
