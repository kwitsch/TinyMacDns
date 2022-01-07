package config

import (
	"time"

	. "github.com/kwitsch/go-dockerutils/config"
)

type Config struct {
	Redis   RedisConfig `koanf:"redis"`
	Verbose bool        `koanf:"verbose" default:"false"`
}

type RedisConfig struct {
	Address  string        `koanf:"address"`
	Username string        `koanf:"username"`
	Password string        `koanf:"password"`
	Database int           `koanf:"database" default:"0"`
	Attempts int           `koanf:"attempts" default:"3"`
	Cooldown time.Duration `koanf:"cooldown" default:"1s"`
	Verbose  bool
}

const prefix = "TMD_"

func Get() (*Config, error) {
	var res Config
	err := Load(prefix, &res)
	if err == nil {

		res.Redis.Verbose = res.Verbose

		return &res, nil
	}
	return nil, err
}
