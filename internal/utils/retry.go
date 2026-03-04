package utils

import (
	"context"
	"time"
)

func Retty(ctx context.Context, attempts int, baseDelay time.Duration, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = fn()
		if err == nil {
			return nil
		}

		if i == attempts-1 {
			break
		}

		delay := baseDelay * time.Duration(1<<i)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return err
}
