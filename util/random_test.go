package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomInt(t *testing.T) {
	randInt := RandomInt(1, 10)

	require.GreaterOrEqual(t, randInt, int64(1))
	require.LessOrEqual(t, randInt, int64(10))
}

func TestRandomString(t *testing.T) {
	randString := RandomString(6)

	require.Len(t, randString, 6)
}

func TestRandomCurrency(t *testing.T) {
	currencies := []string{"USD", "EUR", "MXN"}
	randCurrency := RandomCurrency()

	require.Contains(t, currencies, randCurrency)
}
