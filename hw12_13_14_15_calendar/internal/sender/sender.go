package sender

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Sender struct {
	Logger   logger.Log
	Client   rabbitmq.Client
	Exchange string
	Key      string
}

func New(logger logger.Log, client rabbitmq.Client, exchange, key string) *Sender {
	return &Sender{
		Logger:   logger,
		Client:   client,
		Exchange: exchange,
		Key:      key,
	}
}

func (s *Sender) Run(ctx context.Context, deliveries <-chan amqp.Delivery) error {
	for {
		select {
		case <-ctx.Done():
			return errors.New("sender context is done")
		case d := <-deliveries:
			var message rabbitmq.NotifyMessage
			if err := json.Unmarshal(d.Body, &message); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Can't unmarshal delivery message: [%v] %q", d.DeliveryTag, d.Body))
			}
			s.Logger.Info(fmt.Sprintf("got %dB delivery: [%v] %v", len(d.Body), d.DeliveryTag, &message))
			if err := d.Ack(false); err != nil {
				return errors.Wrap(err, fmt.Sprintf("Can't ack delivery message: [%v] %q", d.DeliveryTag, message))
			}
		}
	}
}
