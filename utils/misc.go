package utils

import (
	"time"
)

// RetriableExecute повторяет функцию до attempt раз с задержкой delay при ошибке.
func RetriableExecute(fn func() error, attempts int, delay time.Duration) error {
	var err error
	for attempts > 0 {
		err = fn()
		if err == nil {
			return nil
		}
		attempts--
		if attempts > 0 {
			time.Sleep(delay)
		}
	}
	return err
}
