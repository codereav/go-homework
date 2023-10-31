package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/codereav/go-homework/hw12_13_14_15_calendar/internal/sender"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var configFile string

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(errors.Wrap(err, "Unable to load .env"))
	}

	err = run(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
			defer cancel()

			config := NewConfig(configFile)
			log := logger.New(config.Logger.Level, config.Logger.Path)

			amqpclient := rabbitmq.NewClient(config.Rabbitmq.Dsn)
			if err := amqpclient.Connect(); err != nil {
				log.Error(errors.Wrap(err, "failed connect to amqp").Error())
				return
			}
			defer func() {
				if err := amqpclient.Shutdown(); err != nil {
					log.Error(errors.Wrap(err, "fail to close amqp connection").Error())
					return
				}
			}()
			if err := amqpclient.ExchangeDeclare(config.Rabbitmq.Exchange, config.Rabbitmq.ExchangeType); err != nil {
				log.Error(errors.Wrap(err, "fail to declare exchange").Error())
				return
			}
			if err := amqpclient.QueueDeclare(config.Rabbitmq.Queue); err != nil {
				log.Error(errors.Wrap(err, "fail to declare queue").Error())
				return
			}
			if err := amqpclient.QueueBind(config.Rabbitmq.Queue,
				config.Rabbitmq.Key, config.Rabbitmq.Exchange); err != nil {
				log.Error(errors.Wrap(err, "fail to bind queue").Error())
				return
			}
			log.Info(fmt.Sprintf("Queue bound to Exchange, starting Consume (consumer tag %q)", config.Rabbitmq.ConsumerTag))
			deliveries, err := amqpclient.Consume(config.Rabbitmq.Queue, config.Rabbitmq.ConsumerTag)
			if err != nil {
				log.Error(errors.Wrap(err, "queue consume").Error())
				return
			}

			send := sender.New(log, amqpclient, config.Rabbitmq.Exchange, config.Rabbitmq.Key)
			log.Info("sender is running...")
			if err := send.Run(ctx, deliveries); err != nil {
				log.Error(errors.Wrap(err, "fail when running sender").Error())
			}
		},
	}
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c",
		"/etc/calendar/sender_config.yaml", "Path to config file")
	err := rootCmd.PersistentFlags().Parse(os.Args)
	if err != nil {
		return err
	}

	return errors.Wrap(rootCmd.ExecuteContext(ctx), "run application")
}
