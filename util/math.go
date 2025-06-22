package util

import (
	"math"
	"math/rand"
	"time"
)

func Round(v float64) float64 {
	return math.Floor(v + 0.5)
}

var count int64

func RoundByArea(startNum, endNum int) int {
	if startNum >= endNum {
		return startNum
	}
	count++
	if count >= 1<<8 {
		count = 0
	}
	rand.NewSource(time.Now().UnixNano() + count)
	rnd := rand.Intn(endNum - startNum)
	return rnd + startNum
}
