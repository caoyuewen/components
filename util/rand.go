package util

import (
	"math/rand"
	"strings"
	"time"
)

func RandStr(base string, count int) string {
	if len(base) == 0 {
		panic("RandStr base is empty")
	}
	if count <= 0 {
		return ""
	}
	res := make([]rune, 0, count)
	baseRune := []rune(base)
	for i := 0; count > i; i++ {
		idx := rand.Intn(len(baseRune))
		res = append(res, baseRune[idx])
	}

	return string(res)
}

func RandName() string {
	rand.NewSource(time.Now().UnixNano())
	index := rand.Intn(5) + 5
	var res string
	for i := 0; i < index; i++ {
		if i == 0 {
			res += strings.ToUpper(RandomLetter())
		} else {
			res += RandomLetter()
		}
	}
	return res
}

// RandomLetter 随机生成一个小写字母
func RandomLetter() string {
	index := rand.Intn(26)
	return string(byte('a' + index))
}

func WeightedRandom(s []int) int {
	if len(s) == 0 {
		return -1
	}

	maxWeight := 0
	for _, w := range s {
		if w > 0 {
			maxWeight += w
		}
	}

	if maxWeight == 0 {
		return -1
	}
	tmp := rand.Intn(maxWeight)
	for i, w := range s {
		if w > 0 {
			if w > tmp {
				return i
			}
			tmp -= w
		}
	}
	return -1
}

func WeightedRandomFunc[T any](s []T, f func(T) int) int {
	if len(s) == 0 {
		return -1
	}

	weights := make([]int, len(s))
	for i, v := range s {
		weights[i] = f(v)
	}
	return WeightedRandom(weights)
}

func Shuffle[T any](s []T) {
	for i := 0; len(s) > i; i++ {
		swapIdx := rand.Intn(len(s))
		s[i], s[swapIdx] = s[swapIdx], s[i]
	}
}
