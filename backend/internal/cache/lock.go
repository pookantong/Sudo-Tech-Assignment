package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrLockNotAcquired = errors.New(
		"seat is already locked by another user",
	)

	ErrLockNotOwned = errors.New(
		"lock is not owned by this user",
	)
)

type SeatLock struct {
	client *redis.Client
	ttl    time.Duration
}

func NewSeatLock(
	client *redis.Client,
	ttl time.Duration,
) *SeatLock {
	return &SeatLock{
		client: client,
		ttl:    ttl,
	}
}

func seatLockKey(
	showtimeID string,
	seatID string,
) string {
	return "lock:seat:" +
		showtimeID +
		":" +
		seatID
}

// ====================
// Acquire
// ====================

func (s *SeatLock) Acquire(
	ctx context.Context,
	showtimeID string,
	seatID string,
	ownerToken string,
) error {
	key := seatLockKey(
		showtimeID,
		seatID,
	)

	ok, err := s.client.SetNX(
		ctx,
		key,
		ownerToken,
		s.ttl,
	).Result()

	if err != nil {
		return err
	}

	if !ok {
		return ErrLockNotAcquired
	}

	return nil
}

// ====================
// Release
// ====================

var releaseScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end
`)

func (s *SeatLock) Release(
	ctx context.Context,
	showtimeID string,
	seatID string,
	ownerToken string,
) error {
	key := seatLockKey(
		showtimeID,
		seatID,
	)

	res, err := releaseScript.Run(
		ctx,
		s.client,
		[]string{key},
		ownerToken,
	).Int()

	if err != nil {
		return err
	}

	if res == 0 {
		return ErrLockNotOwned
	}

	return nil
}

// ====================
// Is Locked
// ====================

func (s *SeatLock) IsLocked(
	ctx context.Context,
	showtimeID string,
	seatID string,
) (bool, string, error) {
	key := seatLockKey(
		showtimeID,
		seatID,
	)

	val, err := s.client.Get(
		ctx,
		key,
	).Result()

	if errors.Is(err, redis.Nil) {
		return false, "", nil
	}

	if err != nil {
		return false, "", err
	}

	return true, val, nil
}