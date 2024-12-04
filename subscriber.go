package main

import (
	"context"
	"log/slog"

	"github.com/eclipse/paho.golang/paho"
)

func SubJob1(ctx context.Context, pr *paho.PublishReceived) error {
	ctx, span := tracer.Start(ctx, "SubJob1")
	defer span.End()

	slog.Info("SubJob1", slog.String("topic", pr.Packet.Topic), slog.String("body", string(pr.Packet.Payload)))

	SubJob2(ctx)
	return nil
}

func SubJob2(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "SubJob2")
	defer span.End()

	slog.Info("SubJob2")

	return nil
}
