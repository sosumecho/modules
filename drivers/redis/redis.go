package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// Redis redis实例
var Redis *Client

type Client struct {
	*redis.Client
}

type Conf struct {
	Host       []string `mapstructure:"host"`
	Password   string   `mapstructure:"password"`
	Master     string   `mapstructure:"master"`
	DB         int      `mapstructure:"db"`
	PoolSize   int      `mapstructure:"pool_size"`
	IsSentinel bool     `mapstructure:"is_sentinel"`
}

func New(conf *Conf) *Client {
	if Redis == nil {
		var c *redis.Client
		if conf.IsSentinel {
			c = redis.NewFailoverClient(&redis.FailoverOptions{
				MasterName:    conf.Master,
				SentinelAddrs: conf.Host,
				Password:      conf.Password,
				DB:            conf.DB,
				PoolSize:      conf.PoolSize,
			})
		} else {
			c = redis.NewClient(&redis.Options{
				Addr:     conf.Host[0],
				Password: conf.Password,
				DB:       conf.DB,
				PoolSize: conf.PoolSize,
			})
		}
		Redis = &Client{
			Client: c,
		}
	}
	return Redis
}

// Lock 加锁
func (c *Client) Lock(key string, expiration time.Duration) error {
	if c.SetNX(context.Background(), key, 1, expiration).Val() {
		return nil
	}
	return errors.New("too frequent")
}

// UnLock 解锁
func (c *Client) UnLock(key string) {
	if c.Exists(context.Background(), key).Val() > 0 {
		c.Del(context.Background(), key)
	}
}

// LockN 加锁
func (c *Client) LockN(key string, expiration time.Duration, n int64) error {
	ctx := context.Background()
	if c.Incr(ctx, key).Val() <= n {
		c.Expire(ctx, key, expiration)
		return nil
	}
	return errors.New("too frequent")
}

// UnLockN 解锁
func (c *Client) UnLockN(key string) {
	if c.Exists(context.Background(), key).Val() > 0 {
		c.Decr(context.Background(), key)
	}
}

// Cache 缓存
func (c *Client) Cache(key string, expiration time.Duration, f func() string) string {
	cmd := c.Get(context.Background(), key)
	var val string
	result, _ := cmd.Result()
	if len(result) == 0 {
		val = f()
		c.Set(context.Background(), key, val, expiration)
		return val
	}
	return result
}
