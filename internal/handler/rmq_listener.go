package handler

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/kainnsoft/migrator/internal/usecase"

	"github.com/kainnsoft/migrator/pkg/rmq"
)

type (
	IConsumer interface {
		Consume(qname string, worker func([]byte))
		CloseConsumer()
	}

	consumer struct {
		//		conn    *amqp.Connection
		ch        *amqp.Channel
		pgUsecase usecase.IPGUsecase
		done      chan struct{}
		logger    zap.Logger
	}
)

func NewListener(rmqChannel rmq.IRmq,
	queueName string,
	pgUsecase usecase.IPGUsecase,
) (c IConsumer) {
	done := make(chan struct{})
	c = &consumer{
		ch:        rmqChannel.Rmq(),
		pgUsecase: pgUsecase,
		done:      done,
	}

	go func() {
		c.Consume(queueName, pgUsecase.ConsumePayment)
	}()

	return c
}

func (c *consumer) Consume(qname string, worker func([]byte)) {
	var (
		q        amqp.Queue
		messages <-chan amqp.Delivery
		err      error
	)
	if q, err = c.ch.QueueDeclare(
		qname, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	); err != nil {
		c.logger.Error("failed to declare a consumer queue; error:", zap.Error(err))
	}

	if messages, err = c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	); err != nil {
		c.logger.Error("failed to register a consumer; error:", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for message := range messages {
			worker(message.Body)
		}
	}(ctx)

	log.Print(" [*] Waiting for messages. To exit press CTRL+C")

	select {
	case <-c.done:
		cancel()
		fmt.Println("consumer stopped")
		return
	}

}

func (c *consumer) CloseConsumer() {
	c.pgUsecase.ClosePGUseCase()
	c.done <- struct{}{}
	close(c.done)
}
