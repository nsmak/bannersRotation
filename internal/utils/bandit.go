package utils

import (
	"math"
)

func PlayWithBandit(counts, rewards []int64) (index int) {
	if len(counts) != len(rewards) {
		panic("\"counts\" length must be equal \"rewards\" length")
	}

	sumCounts := sum(counts...)
	maxValue := math.Inf(-1)
	var maxValueIndex int

	for i, count := range counts {
		k := math.Sqrt((2.0 * math.Log(float64(sumCounts))) / float64(count))
		val := (float64(rewards[i]) / float64(count)) + k
		if val > maxValue {
			maxValue = val
			maxValueIndex = i
		}
	}

	return maxValueIndex
}

func sum(values ...int64) int64 {
	var total int64
	for _, v := range values {
		total += v
	}
	return total
}
