package redistore_test

import (
	"time"

	"github.com/danikarik/okpock/pkg/store/redistore"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

var requiredVars = []string{"REDIS_HOST", "REDIS_PASS"}

func fakeUsername() string {
	return uuid.NewV4().String()
}

func fakeEmail() string {
	return uuid.NewV4().String() + "@example.com"
}

func newTestPool(host, pass string) (*redistore.Pool, error) {
	dialOptions := []redis.DialOption{
		redis.DialDatabase(9),
		redis.DialPassword(pass),
		redis.DialConnectTimeout(15 * time.Second),
		redis.DialReadTimeout(15 * time.Second),
		redis.DialWriteTimeout(15 * time.Second),
	}
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", host, dialOptions...)
		},
		MaxIdle:   5,
		MaxActive: 100,
		Wait:      true,
	}

	c := pool.Get()
	defer c.Close()

	_, err := c.Do("FLUSHDB")
	if err != nil {
		return nil, err
	}

	return redistore.New(pool), nil
}
