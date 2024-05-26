package infra

import (
	"errors"
	"time"
)

// WaitFor retries a check function up to maxAttempts times with an exponential backoff, starting at
// the given timeout. It returns an error if the check does not return true after all attempts.
func WaitFor(check func() bool, maxAttempts int, timeout time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		if check() {
			return nil
		}
		if i == maxAttempts-1 {
			return errors.New("max attempts reached")
		}
		time.Sleep(timeout)
		timeout *= 2
	}
	return nil
}
