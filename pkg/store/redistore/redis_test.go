package redistore_test

import (
	"github.com/danikarik/okpock/pkg/store/redistore"
	"github.com/go-redis/redis"
)

func newTestPool(host, pass string) (*redistore.Pool, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pass,
		DB:       9,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redistore.New(client), nil
}
