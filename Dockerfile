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

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /go/bin/tracer /go/bin/tracer

ENTRYPOINT ["/go/bin/tracer"]
