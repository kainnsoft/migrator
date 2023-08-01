package rmq

import (
	"fmt"

	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/kainnsoft/migrator/config"
)

type (
	IRmq interface {
		Rmq() *amqp.Channel
		RmqGetConn() *amqp.Connection
		CloseRmqChan() error
		ExchangeDeclare(exchangeName, exchangeType string) (err error)
		QueueDeclare(queueName string) (queue amqp.Queue, err error)
		QueueBind(queue amqp.Queue, routingKey, exchangeName string) (err error)
	}
	rmq struct {
		cfg        *config.RMQ
		connection *amqp.Connection
		channel    *amqp.Channel
	}
)

func NewRmqConn(
	cfg *config.RMQ,
) (IRmq, error) {
	var r = &rmq{
		cfg: cfg,
	}

	if err := r.get(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *rmq) get() (err error) {
	// если нет текущего коннекта, то создаем
	if r.channel == nil {
		if err = r.connect(dsn(r.cfg)); err != nil {
			return err
		}
	}

	return nil
}

func dsn(cfg *config.RMQ) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
	)
}

func (r *rmq) connect(connString string) (err error) {
	var (
		conn *amqp.Connection
		ch   *amqp.Channel
	)
	conn, err = amqp.Dial(connString)
	if err != nil {
		return errors.Wrap(err, "error rmq connect")
	}
	r.connection = conn

	ch, err = conn.Channel()
	if err != nil {
		return errors.Wrap(err, "error rmq open channel")
	}
	r.channel = ch

	return nil
}

func (r *rmq) Rmq() *amqp.Channel {
	return r.channel
}

func (r *rmq) RmqGetConn() *amqp.Connection {
	return r.connection
}

func (r *rmq) CloseRmqChan() error {
	if r.channel != nil {
		_ = r.channel.Close()
	}
	if r.connection != nil {
		_ = r.connection.Close()
	}

	return nil
}

func (r *rmq) ExchangeDeclare(exchangeName, exchangeType string) (err error) {
	if err = r.channel.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("exchange declare error: %s", err)
	}

	return nil
}

func (r *rmq) QueueDeclare(queueName string) (queue amqp.Queue, err error) {
	return r.channel.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
}

func (r *rmq) QueueBind(queue amqp.Queue, routingKey, exchangeName string) (err error) {
	if err = r.channel.QueueBind(
		queue.Name,   // name of the queue
		routingKey,   // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("queue bind error: %s", err)
	}

	return nil
}
