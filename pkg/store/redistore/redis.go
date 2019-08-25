package redistore

import (
	"github.com/go-redis/redis"
)

// Pool is a Redis implementation of store interfaces.
type Pool struct{}

// New creates a new instance of `Pool`.
func New(client *redis.Client) *Pool {
	return &Pool{}
}
