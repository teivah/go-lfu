package lfu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teivah/lfu"
)

func TestCache(t *testing.T) {
	var v any
	var b bool

	c := lfu.NewCache[int, string]()

	_, b = c.GetLFU()
	require.Equal(t, false, b)

	v, b = c.Get(1)
	check(t, "", false, v, b)

	c.Set(1, "one")
	v, b = c.GetLFU()
	check(t, 1, true, v, b)

	v, b = c.Get(1)
	check(t, "one", true, v, b)

	c.Set(2, "two")
	_, _ = c.Get(1)
	v, b = c.GetLFU()
	check(t, 2, true, v, b)

	_, _ = c.Get(2)
	_, _ = c.Get(2)
	v, b = c.GetLFU()
	check(t, 1, true, v, b)
}

func check(t *testing.T, expV any, expB bool, gotV any, gotB bool) {
	require.Equal(t, expB, gotB)
	assert.Equal(t, expV, gotV)
}
