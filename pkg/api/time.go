package api

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

// Time custom wrapper around time.
type Time struct{ time.Time }

// NewTime creates wrapper around std time.
func NewTime(t time.Time) Time { return Time{t} }

// Now creates time now.
func Now() Time { return Time{time.Now()} }

// RedisArg implements redis interface.
func (t Time) RedisArg() interface{} {
	return t.Unix()
}

// RedisScan implements redis interface.
func (t *Time) RedisScan(src interface{}) error {
	data, err := redis.Int64(src, nil)
	if err != nil {
		return err
	}
	*t = Time{time.Unix(data, 0)}
	return nil
}
