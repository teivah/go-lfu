package lfu_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teivah/lfu"
)

func TestCache(t *testing.T) {
	var v any
	var b bool

	c := lfu.NewCache[int, string](2)

	_, b = c.GetLFU()
	require.Equal(t, false, b)

	v, b = c.Get(1)
	check(t, "", false, v, b)

	c.Set(1, "one")
	v, b = c.GetLFU()
	check(t, 1, true, v, b)
	v, b = c.Get(1)
	check(t, "one", true, v, b)

	// Update the same value
	c.Set(1, "onex")
	v, b = c.Get(1)
	check(t, "onex", true, v, b)

	c.Set(2, "two")
	_, _ = c.Get(1)
	v, b = c.GetLFU()
	check(t, 2, true, v, b)

	_, _ = c.Get(2)
	_, _ = c.Get(2)
	_, _ = c.Get(2)
	v, b = c.GetLFU()
	check(t, 1, true, v, b)

	c.Set(3, "three")
	v, b = c.GetLFU()
	check(t, 3, true, v, b)
	v, b = c.Get(1)
	check(t, "", false, v, b)

	c.Set(4, "four")
	v, b = c.GetLFU()
	check(t, 4, true, v, b)
	v, b = c.Get(3)
	check(t, "", false, v, b)

	_, _ = c.Get(4)
	_, _ = c.Get(4)
	_, _ = c.Get(4)
	c.Set(5, "five")
}

func check(t *testing.T, expV any, expB bool, gotV any, gotB bool) {
	require.Equal(t, expB, gotB)
	assert.Equal(t, expV, gotV)
}

func TestRace(t *testing.T) {
	n := 1_000_000
	rand.Seed(time.Now().UnixNano())

	c := lfu.NewCache[int, int](100)
	go func() {
		for i := 0; i < n; i++ {
			c.Get(rand.Intn(200))
		}
	}()
	go func() {
		for i := 0; i < n; i++ {
			c.GetLFU()
		}
	}()
	go func() {
		for i := 0; i < n; i++ {
			v := rand.Intn(200)
			c.Set(v, v)
		}
	}()
}
