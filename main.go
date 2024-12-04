package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"log/slog"

	"github.com/eclipse/paho.golang/paho"
	"go.opentelemetry.io/otel"
)

const TOPIC = "demo/topic"

var tracer = otel.Tracer("github.com/shirou/mqttotel")

const serviceName = "mqttotel"

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [command]")
		fmt.Println("Commands: pub, sub")
		os.Exit(1)
	}

	u, err := url.Parse("tcp://localhost:1883")
	if err != nil {
		panic(err)
	}
	args := MQTTConnectionArgs{
		ServerUrls: []*url.URL{u},
	}

	transport, err := NewMQTTTransport(args)
	if err != nil {
		panic(err)
	}

	transport.pubChan = make(chan *paho.Publish)

	command := os.Args[1]

	if err := transport.Connect(ctx); err != nil {
		panic(err)
	}
	go transport.Start(context.Background())
	shutdown, err := setupOTelSDK(ctx, command)
	if err != nil {
		panic(err)
	}

	switch command {
	case "pub":
		if err := PublisherStart(ctx, transport.pubChan); err != nil {
			slog.Error("PublisherStart exit with error", slog.Any("error", err))
		}
	case "sub":
		if err := transport.Subscribe(ctx, []string{TOPIC}); err != nil {
			slog.Error("Subscribe exit with error", slog.Any("error", err))
		}
		<-ctx.Done()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Use 'help' for a list of available commands.")
	}

	if err := shutdown(ctx); err != nil {
		slog.Error("shutdown exit with error", slog.Any("error", err))
	}
}
