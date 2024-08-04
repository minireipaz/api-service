package common

import (
	"crypto/rand"
	"math/big"
	"minireipaz/pkg/domain/models"
	"time"
)

func RandomDuration(min, max time.Duration, i int) time.Duration {
	rangeDuration := max - min

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(rangeDuration)))
	if err != nil {
		return time.Second * 1
	}

	return min + time.Duration(nBig.Int64()) + models.SleepOffset*time.Duration(i)
}
