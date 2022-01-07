package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/kwitsch/TinyMacDns/config"
)

type Client struct {
	cfg    *config.RedisConfig
	client *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new redis client
func New(cfg *config.RedisConfig) (*Client, error) {
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
