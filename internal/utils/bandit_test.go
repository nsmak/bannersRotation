package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayWithBandit(t *testing.T) {
	counts := []int64{6, 7, 5}
	rewards := []int64{1, 2, 1}
	expected := 2

	var index int

	require.NotPanics(t, func() {
		index = PlayWithBandit(counts, rewards)
	})
	require.Equal(t, expected, index)
}

func TestPlayWithBanditInvalidCounts(t *testing.T) {
	counts := []int64{6, 7}
	rewards := []int64{1, 2, 1}

	require.Panics(t, func() {
		PlayWithBandit(counts, rewards)
	})
}

func TestPlayWithBanditInvalidRewards(t *testing.T) {
	counts := []int64{6, 7, 5}
	rewards := []int64{1, 2, 1, 8}

	require.Panics(t, func() {
		PlayWithBandit(counts, rewards)
	})
}

func TestSum(t *testing.T) {
	values := []int64{1, 2, 3}
	var expected int64 = 6

	s := sum(values...)

	require.Equal(t, expected, s)
}

func TestSumZero(t *testing.T) {
	var values []int64
	var expected int64
	s := sum(values...)

	require.Equal(t, expected, s)
}
