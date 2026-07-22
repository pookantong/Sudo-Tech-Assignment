package bootstrap

import (
	"context"
	"log"
	"time"

	"cinema-booking-backend/internal/cache"
	"cinema-booking-backend/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedis(cfg *config.Config) (*redis.Client, *cache.SeatLock) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal(err)
	}

	lock := cache.NewSeatLock(client, cfg.SeatLockTTL)

	return client, lock
}