package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Rabbitmq RabbitmqConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type RabbitmqConf struct {
	Dsn          string
	Exchange     string
	ExchangeType string
	Queue        string
	Key          string
	ConsumerTag  string
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
