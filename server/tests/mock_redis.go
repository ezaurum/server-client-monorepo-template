package templateapp

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type MockRedis struct {
	store map[string]string
}

func NewMockRedis() *MockRedis {
	return &MockRedis{store: make(map[string]string)}
}

func (m *MockRedis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	valStr := value.(string)
	m.store[key] = valStr
	return nil
}

func (m *MockRedis) Get(ctx context.Context, key string) (string, error) {
	val, exists := m.store[key]
	if !exists {
		return val, redis.Nil
	}
	return val, nil
}

func (m *MockRedis) FlushDB(ctx context.Context) *redis.StatusCmd {
	m.store = make(map[string]string)
	return redis.NewStatusResult("OK", nil)
}
