package config

import "github.com/ilyakaznacheev/cleanenv"

type Service struct {
	Name string `env:"NAME" env-default:"subscriptions"`
	Host string `env:"HOST" env-required:"true"`
	Port int    `env:"PORT" env-required:"true"`
}

type Postgres struct {
	Host     string `env:"HOST"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
	SSLMode  string `env:"SSL_MODE"`
	Port     int    `env:"PORT"`
	MaxConns int    `env:"MAX_CONNS" env-default:"20"`
	MinConns int    `env:"MIN_CONNS" env-default:"2"`
}

type Config struct {
	Postgres Postgres `env-prefix:"DB_"`
	Service  Service  `env-prefix:"SERVICE_"`
}

func Load() (Config, error) {
	cfg := Config{}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
