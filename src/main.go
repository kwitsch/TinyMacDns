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
			defer redis.Close()
			fmt.Println("Server starting")

			server := server.New(cache, int(cfg.Redis.Intervall.Seconds()), cfg.Verbose)
			defer server.Stop()
			server.Start()

			ticker := time.NewTicker(cfg.Redis.Intervall)
			defer ticker.Stop()

			intChan := make(chan os.Signal, 1)
			signal.Notify(intChan, os.Interrupt)
			defer close(intChan)

			redis.Poll(&cfg.Hosts)

			for {
				select {
				case sErr := <-server.Error:
					fmt.Println(sErr)
					os.Exit(3)
				case <-ticker.C:
					redis.Poll(&cfg.Hosts)
				case <-intChan:
					fmt.Println("Server stopping")
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
