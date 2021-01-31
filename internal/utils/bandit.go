package utils

import (
	"errors"
	"math"
)

func PlayWithBandit(counts, rewards []int64) (index int, err error) {
	if len(counts) != len(rewards) {
		return 0, errors.New("invalid length of counts/rewards")
	}

	sumCounts := sum(counts...)
	values := make([]float64, len(counts))

	for i, count := range counts {
		k := math.Sqrt((2.0 * math.Log(float64(sumCounts))) / float64(count))
		val := (float64(rewards[i]) / float64(count)) + k
		values[i] = val
	}

	return maxAt(values...), nil
}

func sum(values ...int64) int64 {
	var total int64
	for _, v := range values {
		total += v
	}
	return total
}

func maxAt(values ...float64) (index int) {
	value := math.Inf(-1)
	for i, v := range values {
		if v > value {
			value = v
			index = i
		}
	}
	return
}
