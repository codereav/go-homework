package rabbitmq

import (
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Client interface {
	Connect() error
	Shutdown() error
	ExchangeDeclare(exchange, exchangeType string) error
	QueueDeclare(queueName string) error
	QueueBind(queueName, key, exchange string) error
	Publish(exchange, routingKey string, body []byte) error
	NotifyPublish(confirm chan amqp.Confirmation) (chan amqp.Confirmation, error)
	Consume(queueName, tag string) (<-chan amqp.Delivery, error)
}

type AMQPClient struct {
	Client
	dsn     string
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewClient(dsn string) Client {
	return &AMQPClient{
		dsn: dsn,
	}
}

func (c *AMQPClient) Connect() error {
	var err error
	c.conn, err = amqp.Dial(c.dsn)
	if err != nil {
		return errors.Wrap(err, "trying connect to rabbitmq")
	}
	c.channel, err = c.conn.Channel()
	if err != nil {
		return errors.Wrap(err, "trying to create channel")
	}

	return nil
}

func (c *AMQPClient) Shutdown() error {
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "AMQP connection close error")
	}

	return nil
}

func (c *AMQPClient) ExchangeDeclare(exchange, exchangeType string) error {
	if c.channel == nil {
		return errors.New("exchange declare: channel is not initialized")
	}
	err := c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return errors.Wrap(err, "exchange declare")
	}

	return nil
}

func (c *AMQPClient) QueueDeclare(queueName string) error {
	if c.channel == nil {
		return errors.New("queue declare: channel is not initialized")
	}
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return errors.Wrap(err, "queue declare")
	}
	c.queue = &queue

	return nil
}

func (c *AMQPClient) QueueBind(queueName, key, exchange string) error {
	if c.channel == nil {
		return errors.New("queue bind: channel is not initialized")
	}
	err := c.channel.QueueBind(
		queueName, // name of the queue
		key,       // bindingKey
		exchange,  // sourceExchange
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return errors.Wrap(err, "queue bind")
	}

	return nil
}

func (c *AMQPClient) Publish(exchange, routingKey string, body []byte) error {
	if c.channel == nil {
		return errors.New("publish: channel is not initialized")
	}
	err := c.channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	)
	if err != nil {
		return errors.Wrap(err, "publishing message")
	}

	return nil
}

func (c *AMQPClient) NotifyPublish(confirm chan amqp.Confirmation) (chan amqp.Confirmation, error) {
	if c.channel == nil {
		return nil, errors.New("notifyPublish: channel is not initialized")
	}
	return c.channel.NotifyPublish(confirm), nil
}

func (c *AMQPClient) Consume(queueName, tag string) (<-chan amqp.Delivery, error) {
	if c.channel == nil {
		return nil, errors.New("queue consume: channel is not initialized")
	}
	deliveries, err := c.channel.Consume(
		queueName, // name
		tag,       // consumerTag,
		false,     // noAck
		false,     // exclusive
		false,     // noLocal
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, errors.Wrap(err, "queue consume")
	}

	return deliveries, nil
}
