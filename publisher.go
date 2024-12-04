package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/eclipse/paho.golang/paho"
)

func PublisherStart(ctx context.Context, pubchan chan *paho.Publish) error {
	ctx, span := tracer.Start(ctx, "PublisherStart")
	defer span.End()

	PubJob1(ctx, pubchan)

	return nil
}

func PubJob1(ctx context.Context, pubchan chan *paho.Publish) error {
	ctx, span := tracer.Start(ctx, "PubJob1")
	defer span.End()

	slog.Info("PubJob1")
	for i := 0; i < 2; i++ {
		actualPublish(ctx, pubchan, i)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func actualPublish(ctx context.Context, pubchan chan *paho.Publish, i int) error {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("actual publish: %d", i))
	defer span.End()

	slog.Info("actualPublish")

	p := &paho.Publish{
		Topic:   TOPIC,
		Payload: []byte(fmt.Sprintf("Hello, World!: %d", i)),
	}
	inject(ctx, p)

	pubchan <- p

	return nil
}
