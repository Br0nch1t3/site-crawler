package logger

import (
	"fmt"
	"time"
)

var prevDuration time.Duration

func WithTimer[T any](fn func() (T, error), prefix ...string) (T, error) {
	now := time.Now()
	res, err := fn()
	duration := time.Since(now) - prevDuration
	prevDuration = duration
	text := duration.String()
	if len(prefix) > 0 {
		text = prefix[0] + ": " + text
	}
	fmt.Println(text)
	return res, err
}
