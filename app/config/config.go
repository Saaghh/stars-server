package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Postgres   PostgresConfig
	LogLevel   string `envconfig:"LOG_LEVEL" default:"debug"`
	ServerAddr string `envconfig:"SERVER_ADDR" default:":8080"`
}

type PostgresConfig struct {
	Host     string `envconfig:"HOST" default:"localhost"`
	Port     string `envconfig:"PORT" default:"5432"`
	DataBase string `envconfig:"DB" default:"stars"`
	User     string `envconfig:"USER" default:"postgres"`
	Password string `envconfig:"PASSWORD" default:"postgres"`
}

func NewFromEnv(prefix string) Config {
	c := Config{}
	envconfig.MustProcess(prefix, &c)
	return c
}
