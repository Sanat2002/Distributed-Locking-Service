package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// Client is a thin wrapper around go-redis client
type Client struct {
	RDB *goredis.Client
	Scripts *Scripts
}

// NewClient connects to Redis and fails fast if unavailable
func NewClient(addr string) *Client {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return &Client{
		RDB: rdb,
		Scripts: LoadScripts(),
	}
}
