package main

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

type MQTTConnectionArgs struct {
	ServerUrls []*url.URL
}

type MQTTTransport struct {
	clientConfig autopaho.ClientConfig
	connection   *autopaho.ConnectionManager

	pubChan chan *paho.Publish //
}

func NewMQTTTransport(args MQTTConnectionArgs) (*MQTTTransport, error) {
	transport := &MQTTTransport{}

	cliCfg := autopaho.ClientConfig{
		ServerUrls:                    args.ServerUrls,
		CleanStartOnInitialConnection: false,
		ClientConfig: paho.ClientConfig{
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				transport.received},
		},
	}

	transport.clientConfig = cliCfg

	return transport, nil
}

func (tra *MQTTTransport) Subscribe(ctx context.Context, topics []string) error {
	subscribeTopics := make([]paho.SubscribeOptions, 0)
	for _, topic := range topics {
		subscribeTopics = append(subscribeTopics, paho.SubscribeOptions{
			Topic: topic,
			QoS:   0,
		})
	}
	_, err := tra.connection.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: subscribeTopics,
	})
	return err
}

func (tra *MQTTTransport) received(pr paho.PublishReceived) (bool, error) {
	ctx := extract(context.Background(), &pr)

	SubJob1(ctx, &pr)

	return true, nil
}

func (tra *MQTTTransport) publish(ctx context.Context, conn *autopaho.ConnectionManager, p *paho.Publish) error {
	if _, err := conn.Publish(ctx, p); err != nil {
		if ctx.Err() == nil {
			slog.Error("publish failed due to context cancelled", slog.Any("error", err))
			return fmt.Errorf("publish failed due to context cancelled")
		}
	}
	slog.Info("publish success", slog.String("topic", p.Topic))
	return nil
}

func (tra *MQTTTransport) Connect(ctx context.Context) error {
	conn, err := autopaho.NewConnection(ctx, tra.clientConfig)
	if err != nil {
		return fmt.Errorf("NewConnection failed, %w", err)
	}
	// Wait for the connection to come up
	slog.Info("connect started")
	if err = conn.AwaitConnection(ctx); err != nil {
		return fmt.Errorf("AwaitConnection failed, %w", err)
	}
	tra.connection = conn
	slog.Info("connection established")
	return nil
}

func (tra *MQTTTransport) Start(ctx context.Context) error {
	for {
		select {
		case msg, ok := <-tra.pubChan:
			if !ok {
				slog.Warn("dataChan closed")
				return nil
			}
			if err := tra.publish(ctx, tra.connection, msg); err != nil {
				slog.Error("publish failed", slog.Any("error", err))
			}
		case <-tra.connection.Done():
			slog.Info("MQTT connection finished")
			return nil
		case <-ctx.Done():
			slog.Info("MQTT connection finished by context")
			return nil
		}
	}
}
