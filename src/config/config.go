package config

import (
	"fmt"
	"time"

	. "github.com/kwitsch/go-dockerutils/config"
)

type Config struct {
	Redis   RedisConfig           `koanf:"redis"`
	Hosts   map[string]HostConfig `koanf:"hosts"`
	Verbose bool                  `koanf:"verbose" default:"false"`
}

type RedisConfig struct {
	Address   string        `koanf:"address"`
	Username  string        `koanf:"username"`
	Password  string        `koanf:"password"`
	Database  int           `koanf:"database" default:"0"`
	Attempts  int           `koanf:"attempts" default:"3"`
	Cooldown  time.Duration `koanf:"cooldown" default:"1s"`
	Intervall time.Duration `koan:"intervall" default:"5m"`
	Verbose   bool
}

type HostConfig struct {
	Mac map[int]string `koanf:"mac"`
}

const prefix = "TMD_"

func Get() (*Config, error) {
	var res Config
	err := Load(prefix, &res)
	if err == nil {

		res.Redis.Verbose = res.Verbose

		if res.Verbose {
			logHosts(res.Hosts)
		}

		return &res, nil
	}

	return nil, err
}

func logHosts(hosts map[string]HostConfig) {
	fmt.Println("Configured hosts:")
	for k, v := range hosts {
		macs := ""
		for _, m := range v.Mac {
			if len(macs) > 0 {
				macs += ", "
			}
			macs += m
		}
		fmt.Println("-", k, ":", macs)
	}
}
