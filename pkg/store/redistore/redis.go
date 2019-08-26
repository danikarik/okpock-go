package redistore

import (
	"github.com/gomodule/redigo/redis"
)

// New creates a new instance of `Pool`.
func New(p *redis.Pool) *Pool { return &Pool{p} }

// Pool is a Redis implementation of store interfaces.
type Pool struct{ *redis.Pool }
