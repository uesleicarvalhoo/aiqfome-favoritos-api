package cache

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Host     string
	Port     string
	User     string
	Password string
	UseSSL   bool
}

func (o *RedisOptions) getUrl() string {
	return net.JoinHostPort(o.Host, o.Port)
}

type Redis struct {
	cli *redis.Client
}

func NewRedis(opts RedisOptions) (*Redis, error) {
	o := &redis.Options{
		Addr: opts.getUrl(),
		DB:   0,
	}

	if opts.Password != "" {
		o.Password = opts.Password
	}

	cli := redis.NewClient(o)

	if err := redisotel.InstrumentTracing(cli, ); err != nil{
		return nil, err
	}
	if err:= redisotel.InstrumentMetrics(cli); err != nil{
		return nil, err
	}

	_, err := cli.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{
		cli: cli,
	}, nil
}

func (r *Redis) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.cli.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.cli.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return []byte(data), nil
}

func (r *Redis) Del(ctx context.Context, key string) error {
	res := r.cli.Del(ctx, key)
	_, err := res.Result()

	return err
}
