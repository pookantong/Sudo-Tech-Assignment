package cache

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupTestRedis(
	t *testing.T,
) *redis.Client {
	t.Helper()

	mr := miniredis.RunT(t)

	return redis.NewClient(
		&redis.Options{
			Addr: mr.Addr(),
		},
	)
}

func TestConcurrentAcquire_OnlyOneWinner(
	t *testing.T,
) {
	client := setupTestRedis(t)

	lock := NewSeatLock(
		client,
		5*time.Minute,
	)

	const attackers = 100

	showtimeID := "showtime-1"
	seatID := "seat-A1"

	var successCount int64

	var wg sync.WaitGroup

	start := make(chan struct{})

	for i := 0; i < attackers; i++ {
		wg.Add(1)

		go func(userNum int) {
			defer wg.Done()

			<-start

			ownerToken := fmt.Sprintf(
				"user-%d-token",
				userNum,
			)

			err := lock.Acquire(
				context.Background(),
				showtimeID,
				seatID,
				ownerToken,
			)

			if err == nil {
				atomic.AddInt64(
					&successCount,
					1,
				)
			}
		}(i)
	}

	close(start)

	wg.Wait()

	if successCount != 1 {
		t.Fatalf(
			"expected exactly 1 successful lock acquisition, got %d — DOUBLE BOOKING BUG",
			successCount,
		)
	}
}

func TestConcurrentAcquire_DifferentSeatsAllSucceed(
	t *testing.T,
) {
	client := setupTestRedis(t)

	lock := NewSeatLock(
		client,
		5*time.Minute,
	)

	const seats = 50

	showtimeID := "showtime-1"

	var successCount int64

	var wg sync.WaitGroup

	for i := 0; i < seats; i++ {
		wg.Add(1)

		go func(seatNum int) {
			defer wg.Done()

			seatID := fmt.Sprintf(
				"seat-%d",
				seatNum,
			)

			err := lock.Acquire(
				context.Background(),
				showtimeID,
				seatID,
				fmt.Sprintf(
					"owner-%d",
					seatNum,
				),
			)

			if err == nil {
				atomic.AddInt64(
					&successCount,
					1,
				)
			}
		}(i)
	}

	wg.Wait()

	if successCount != seats {
		t.Fatalf(
			"expected all %d distinct seats to lock independently, got %d successes",
			seats,
			successCount,
		)
	}
}

func TestRelease_OnlyOwnerCanRelease(
	t *testing.T,
) {
	client := setupTestRedis(t)

	lock := NewSeatLock(
		client,
		5*time.Minute,
	)

	ctx := context.Background()

	showtimeID := "showtime-1"
	seatID := "seat-A1"

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-A-token",
	); err != nil {
		t.Fatalf(
			"initial acquire should succeed: %v",
			err,
		)
	}

	if err := lock.Release(
		ctx,
		showtimeID,
		seatID,
		"user-A-token",
	); err != nil {
		t.Fatalf(
			"owner release should succeed: %v",
			err,
		)
	}

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-B-token",
	); err != nil {
		t.Fatalf(
			"user B acquire should succeed: %v",
			err,
		)
	}

	if err := lock.Release(
		ctx,
		showtimeID,
		seatID,
		"user-A-token",
	); err != ErrLockNotOwned {
		t.Fatalf(
			"expected ErrLockNotOwned when non-owner releases, got: %v",
			err,
		)
	}

	locked, owner, err := lock.IsLocked(
		ctx,
		showtimeID,
		seatID,
	)

	if err != nil {
		t.Fatalf(
			"IsLocked error: %v",
			err,
		)
	}

	if !locked || owner != "user-B-token" {
		t.Fatalf(
			"expected user-B-token to still hold the lock, got locked=%v owner=%s",
			locked,
			owner,
		)
	}
}

func TestAcquire_TTLExpiryFreesSeat(
	t *testing.T,
) {
	client := setupTestRedis(t)

	lock := NewSeatLock(
		client,
		50*time.Millisecond,
	)

	ctx := context.Background()

	showtimeID := "showtime-1"
	seatID := "seat-A1"

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-A-token",
	); err != nil {
		t.Fatalf(
			"initial acquire should succeed: %v",
			err,
		)
	}

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-B-token",
	); err != ErrLockNotAcquired {
		t.Fatalf(
			"expected user B to be blocked before TTL expiry, got: %v",
			err,
		)
	}

	time.Sleep(
		100 * time.Millisecond,
	)

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-B-token",
	); err != nil {
		t.Fatalf(
			"expected user B to acquire after TTL expiry, got: %v",
			err,
		)
	}
}

func TestExtend_OnlyOwnerCanExtend(
	t *testing.T,
) {
	client := setupTestRedis(t)

	lock := NewSeatLock(
		client,
		50*time.Millisecond,
	)

	ctx := context.Background()

	showtimeID := "showtime-1"
	seatID := "seat-A1"

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-A-token",
	); err != nil {
		t.Fatalf(
			"initial acquire should succeed: %v",
			err,
		)
	}

	if err := lock.Extend(
		ctx,
		showtimeID,
		seatID,
		"user-B-token",
		5*time.Minute,
	); err != ErrLockNotOwned {
		t.Fatalf(
			"expected ErrLockNotOwned for non-owner extend, got: %v",
			err,
		)
	}

	if err := lock.Extend(
		ctx,
		showtimeID,
		seatID,
		"user-A-token",
		5*time.Minute,
	); err != nil {
		t.Fatalf(
			"owner should be able to extend lock: %v",
			err,
		)
	}

	time.Sleep(
		100 * time.Millisecond,
	)

	if err := lock.Acquire(
		ctx,
		showtimeID,
		seatID,
		"user-B-token",
	); err != ErrLockNotAcquired {
		t.Fatalf(
			"expected lock to remain active after extension, got: %v",
			err,
		)
	}
}

func TestIsLocked_ReturnsFalseWhenMissing(
	t *testing.T,
) {
	client := setupTestRedis(t)

	lock := NewSeatLock(
		client,
		5*time.Minute,
	)

	locked, owner, err := lock.IsLocked(
		context.Background(),
		"showtime-1",
		"seat-A1",
	)

	if err != nil {
		t.Fatalf(
			"IsLocked should not return error for missing lock: %v",
			err,
		)
	}

	if locked {
		t.Fatal(
			"expected seat to be unlocked",
		)
	}

	if owner != "" {
		t.Fatalf(
			"expected empty owner, got %q",
			owner,
		)
	}
}