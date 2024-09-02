package common

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

func RandomDuration(max, min time.Duration, attempt int) time.Duration {
	baseWait := min * time.Duration(math.Pow(2, float64(attempt)))

	jitterRange := int64(max - min)
	jitterBig, err := rand.Int(rand.Reader, big.NewInt(jitterRange))
	if err != nil {
		return baseWait
	}

	jitter := time.Duration(jitterBig.Int64()) + min

	waitTime := baseWait + jitter

	return waitTime
}
