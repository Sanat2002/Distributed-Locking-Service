package redis

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type Scripts struct {
	Acquire *redis.Script
}

func LoadScripts() *Scripts {
	data, err := os.ReadFile("internal/redis/acquire.lua")
	if err != nil {
		panic("failed to load acquire.lua: " + err.Error())
	}

	return &Scripts{
		Acquire: redis.NewScript(string(data)),
	}
}
