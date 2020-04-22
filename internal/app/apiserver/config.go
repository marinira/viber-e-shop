package apiserver

import (
	"github.com/marinira/http-rest-api/internal/app/db"
)

//конфигурация API сервера
type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Mongo    db.DataBaseInterface
}

// создаем новую конфигурацию
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
		Mongo:    db.NewMongo(),
	}
}
