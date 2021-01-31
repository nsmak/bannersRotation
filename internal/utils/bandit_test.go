package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlayWithBandit(t *testing.T) {
	counts := []int64{6, 7, 5}
	rewards := []int64{1, 2, 1}
	expected := 2

	index, err := PlayWithBandit(counts, rewards)

	require.NoError(t, err)
	require.Equal(t, expected, index)
}

func TestPlayWithBanditInvalidCounts(t *testing.T) {
	counts := []int64{6, 7}
	rewards := []int64{1, 2, 1}
	expected := 0

	index, err := PlayWithBandit(counts, rewards)

	require.Error(t, err)
	require.Equal(t, expected, index)
}

func TestPlayWithBanditInvalidRewards(t *testing.T) {
	counts := []int64{6, 7, 5}
	rewards := []int64{1, 2, 1, 8}
	expected := 0

	index, err := PlayWithBandit(counts, rewards)

	require.Error(t, err)
	require.Equal(t, expected, index)
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

func TestMaxAt(t *testing.T) {
	values := []float64{1, 2, 3}
	var expected = 2

	m := maxAt(values...)

	require.Equal(t, expected, m)
}

func TestMaxAtWhenAllMax(t *testing.T) {
	values := []float64{3, 3, 3}
	var expected = 0

	m := maxAt(values...)

	require.Equal(t, expected, m)
}

func TestMaxAtWhenZero(t *testing.T) {
	var values []float64
	var expected = 0

	m := maxAt(values...)

	require.Equal(t, expected, m)
}
