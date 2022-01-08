package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kwitsch/TinyMacDns/cache"
	"github.com/kwitsch/TinyMacDns/config"
	"github.com/kwitsch/TinyMacDns/redis"
	"github.com/kwitsch/TinyMacDns/server"

	_ "github.com/kwitsch/go-dockerutils"
)

func main() {
	cfg, cErr := config.Get()
	if cErr == nil {
		cache := cache.New()

		redis, rErr := redis.New(&cfg.Redis, cache)
		if rErr == nil {
			server := server.New(cache, int(cfg.Redis.Intervall.Seconds()))
			server.Start()

			ticker := time.NewTicker(cfg.Redis.Intervall)

			intChan := make(chan os.Signal, 1)
			signal.Notify(intChan, os.Interrupt)

			for {
				select {
				case <-ticker.C:
					redis.Poll(&cfg.Hosts)
				case <-intChan:
					fmt.Println("Collector stopping")
					server.Stop()
					ticker.Stop()
					redis.Close()
					close(intChan)
					os.Exit(0)
				}
			}
		} else {
			fmt.Println(rErr)
			os.Exit(2)
		}
	} else {
		fmt.Println(cErr)
		os.Exit(1)
	}
}
