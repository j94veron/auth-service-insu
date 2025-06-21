package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewClient(addr, password string, db int) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *Client) SaveToken(ctx context.Context, uuid string, userID uint, expiration time.Duration) error {
	return c.client.Set(ctx, uuid, userID, expiration).Err()
}

func (c *Client) DeleteToken(ctx context.Context, uuid string) error {
	return c.client.Del(ctx, uuid).Err()
}

func (c *Client) GetUserID(ctx context.Context, uuid string) (uint, error) {
	val, err := c.client.Get(ctx, uuid).Uint64()
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
