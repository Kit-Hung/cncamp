package util

import (
	"math/rand"
	"time"
)

func RandInt(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + rand.Intn(max-min+1)
}
