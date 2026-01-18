package lock

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	rds "github.com/username/distributed-lock-service/internal/redis"
)

type Manager struct {
	redis  *rds.Client
	nodeID string
}

func NewManager(redis *rds.Client) *Manager {
	return &Manager{
		redis:  redis,
		nodeID: uuid.NewString(),
	}
}

func (m *Manager) Acquire(
	ctx context.Context,
	resourceID string,
	ttl time.Duration,
) (int64, error) {

	lockKey := fmt.Sprintf("lock:%s", resourceID)
	fencingKey := fmt.Sprintf("fencing:%s", resourceID)

	now := time.Now().UnixMilli()

	log.Println("[LOCK] trying to acquire lock:", lockKey)

	result, err := m.redis.Scripts.Acquire.Run(
		ctx,
		m.redis.RDB,
		[]string{lockKey, fencingKey},
		m.nodeID,
		ttl.Milliseconds(),
		now,
	).Result()

	if err != nil {
		log.Println("[LOCK] redis error:", err)
		return 0, err
	}

	values := result.([]interface{})
	granted := values[0].(int64)

	if granted == 0 {
		log.Println("[LOCK] lock busy")
		return 0, fmt.Errorf("lock busy")
	}

	token := values[1].(int64)
	log.Println("[LOCK] lock granted, token:", token)

	return token, nil
}
