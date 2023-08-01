package repository

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/kainnsoft/migrator/internal/entity"
)

type (
	IRmqRepo interface {
		PushMsg(ctx context.Context, data []byte) (err error)
	}
	rmqRepo struct {
		ch     *amqp.Channel
		logger *zap.Logger
	}
)

func NewRmqRepo(ch *amqp.Channel, logger *zap.Logger) IRmqRepo {
	return &rmqRepo{
		ch:     ch,
		logger: logger,
	}
}

func (r *rmqRepo) PushMsg(ctx context.Context, data []byte) (err error) {
	var (
		cancel context.CancelFunc
	)
	ctx, cancel = context.WithTimeout(ctx, mySqlTimeout)
	defer cancel()

	if err = r.ch.PublishWithContext(ctx,
		entity.Exchange,           // exchange  // default
		entity.RoutingKeyPayments, // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			Body: data,
		}); err != nil {
		return err
	}

	return nil
}
