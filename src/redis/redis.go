package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kwitsch/TinyMacDns/cache"
	"github.com/kwitsch/TinyMacDns/config"
)

type Client struct {
	cfg    *config.RedisConfig
	client *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
	cache  *cache.Cache
}

// New creates a new redis client
func New(cfg *config.RedisConfig, cache *cache.Cache) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Address,
		Username:        cfg.Username,
		Password:        cfg.Password,
		DB:              cfg.Database,
		MaxRetries:      cfg.Attempts,
		MaxRetryBackoff: time.Duration(cfg.Cooldown),
	})
	ctx, cancel := context.WithCancel(context.Background())

	_, err := rdb.Ping(ctx).Result()
	if err == nil {
		res := &Client{
			cfg:    cfg,
			client: rdb,
			ctx:    ctx,
			cancel: cancel,
			cache:  cache,
		}
		return res, nil
	}
	cancel()
	return nil, err
}

// Close discards the redis client
func (c *Client) Close() {
	c.cancel()
}

func (c *Client) Poll(hosts *map[string]config.HostConfig) {
	if c.cfg.Verbose {
		fmt.Println("redis.Client.Poll")
	}

	for hostname, host := range *hosts {
		c.pollHost(hostname, host)
	}
}

func (c *Client) pollHost(hostname string, host config.HostConfig) {
	found := false
	for _, mac := range host.Mac {
		ip, err := c.client.Get(c.ctx, mac).Result()
		if err == nil {
			if c.cfg.Verbose {
				fmt.Println("redis.Client.pollhost", hostname, "=", ip)
			}
			c.cache.Update(hostname, ip)
			found = true
			break
		} else {
			fmt.Println(err)
		}
	}
	if !found {
		if c.cfg.Verbose {
			fmt.Println("redis.Client.pollhost", hostname, "not found")
		}
		c.cache.Delete(hostname)
	}
}
