package stm

import (
	"runtime/pprof"
)

var retries = pprof.NewProfile("stmRetries")

// retry is a sentinel value. When thrown via panic, it indicates that a
// transaction should be retried.
var retry = &struct{}{}

// catchRetry returns true if fn calls tx.Retry.
func catchRetry[R any](fn Operation[R], tx *Tx) (result R, gotRetry bool) {
	defer func() {
		if r := recover(); r == retry {
			gotRetry = true
		} else if r != nil {
			panic(r)
		}
	}()
	result = fn(tx)
	return
}
