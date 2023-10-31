package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Scheduler SchedulerConf
	Logger    LoggerConf
	Database  DatabaseConf
	Server    struct {
		HTTP ServerConf
		Grpc ServerConf
	}
	Rabbitmq RabbitmqConf
}

type SchedulerConf struct {
	PeriodSec int16
	OldDate   time.Time
}

type LoggerConf struct {
	Level string
	Path  string
}

type DatabaseConf struct {
	Type string
	Dsn  string
}

type RabbitmqConf struct {
	Dsn          string
	Exchange     string
	ExchangeType string
	Queue        string
	Key          string
	ConsumerTag  string
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig(configFilePath string) *Config {
	viper.SetConfigFile(configFilePath)
	viper.SetEnvPrefix("CONFIG")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Ошибка чтения файла конфигурации:", err)
	}

	// Создаем новый экземпляр структуры Config
	config := &Config{}

	// Распаковываем данные в структуру Config
	if err := viper.Unmarshal(config); err != nil {
		fmt.Println("Ошибка при распаковке данных конфигурации:", err)
	}

	return config
}
