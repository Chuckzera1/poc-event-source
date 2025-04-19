package infrastructure

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"poc-event-source/internal/application"
)

type NatsBroker struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	stream jetstream.Stream
}

type natsSubscription struct {
	cons jetstream.ConsumeContext
}

func Nats(url string, steamName string, subjects []string, ctx context.Context, cancel context.CancelFunc) (*NatsBroker, error) {
	defer cancel()

	nc, err := nats.Connect(url)
	if err != nil {
		log.Printf("Error connecting to nats\n URL: %s \n ERR: %s", url, err.Error())
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		log.Printf("Error creating jetstream\n ERR: %s", err.Error())
		return nil, err
	}

	s, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{Name: steamName, Subjects: subjects})
	if err != nil {
		log.Printf("Error creating stream\n ERR: %s", err.Error())
		return nil, err
	}

	return &NatsBroker{nc: nc, js: js, stream: s}, nil
}

func (b *NatsBroker) Publish(ctx context.Context, topic string, data []byte) error {
	_, err := b.js.Publish(ctx, topic, data)
	return err
}

func (b *NatsBroker) Subscribe(
	ctx context.Context,
	subject string,
	handler func(context.Context, *application.Message),
) (application.Subscription, error) {
	c, err := b.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("durable-%s", subject),
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	cons, err := c.Consume(func(msg jetstream.Msg) {
		handler(ctx, &application.Message{
			Topic: msg.Subject(),
			Data:  msg.Data(),
			Ack:   msg.Ack,
		})
	})
	if err != nil {
		return nil, err
	}

	return &natsSubscription{cons}, nil
}

func (b *NatsBroker) QueueSubscribe(
	ctx context.Context,
	subject, queue string,
	handler func(context.Context, *application.Message),
) (application.Subscription, error) {
	c, err := b.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       fmt.Sprintf("durable-%s-%s", subject, queue),
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	cons, err := c.Consume(func(msg jetstream.Msg) {
		handler(ctx, &application.Message{
			Topic: msg.Subject(),
			Data:  msg.Data(),
			Ack:   msg.Ack,
		})
	})
	if err != nil {
		return nil, err
	}

	return &natsSubscription{cons}, nil
}

func (b *NatsBroker) Close() error {
	err := b.nc.Drain()
	if err != nil {
		return err
	}

	b.nc.Close()
	return nil
}

func (s *natsSubscription) Unsubscribe() error {
	s.cons.Stop()

	return nil
}
