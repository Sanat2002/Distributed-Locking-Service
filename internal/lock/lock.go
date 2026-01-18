package lock

import (
	"context"

	rds "github.com/username/distributed-lock-service/internal/redis"
)

// Manager owns all lock-related state
type Manager struct {
	redis *rds.Client
}

func NewManager(redisClient *rds.Client) *Manager {
	return &Manager{
		redis: redisClient,
	}
}

// HealthCheck verifies server -> Redis connectivity
func (m *Manager) HealthCheck(ctx context.Context) error {
	return m.redis.RDB.Ping(ctx).Err()
}
