package main

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger   LoggerConf
	Database DatabaseConf
	Server   ServerConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type DatabaseConf struct {
	Type string
	Dsn  string
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

// TODO
