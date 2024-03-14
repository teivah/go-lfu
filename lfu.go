// Package lfu is Least Frequently Used Cache with all the operations in O(1).
package lfu

import "sync"

type freqNode[K comparable, V any] struct {
	value int
	items map[K]struct{}
	prev  *freqNode[K, V]
	next  *freqNode[K, V]
}

func newFrequencyNode[K comparable, V any]() *freqNode[K, V] {
	n := &freqNode[K, V]{
		items: make(map[K]struct{}),
	}
	n.prev = n
	n.next = n
	return n
}

type item[K comparable, V any] struct {
	data   V
	parent *freqNode[K, V]
}

func newItem[K comparable, V any](data V, parent *freqNode[K, V]) *item[K, V] {
	return &item[K, V]{
		data:   data,
		parent: parent,
	}
}

// Cache structure.
type Cache[K comparable, V any] struct {
	mu        sync.RWMutex
	length    int
	maxLength int
	byKey     map[K]*item[K, V]
	freq      *freqNode[K, V]
}

// NewCache creates a new cache.
func NewCache[K comparable, V any](size int) *Cache[K, V] {
	return &Cache[K, V]{
		maxLength: size,
		byKey:     make(map[K]*item[K, V]),
		freq:      newFrequencyNode[K, V](),
	}
}

// Get gets an element.
func (c *Cache[K, V]) Get(k K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	tmp, exists := c.byKey[k]
	if !exists {
		return zero, false
	}
	freq := tmp.parent
	nextFreq := freq.next
	if nextFreq == c.freq || nextFreq.value != freq.value+1 {
		nextFreq = getNewNode(freq.value+1, freq, nextFreq)
	}
	nextFreq.items[k] = struct{}{}
	tmp.parent = nextFreq
	delete(freq.items, k)
	if len(freq.items) == 0 {
		deleteNode(freq)
	}
	return tmp.data, true
}

// Set inserts an element.
func (c *Cache[K, V]) Set(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.byKey[k]; exists {
		return
	}

	c.length++
	if c.length > c.maxLength {
		c.deleteLFU()
		c.length--
	}

	freq := c.freq.next
	if freq.value != 1 {
		freq = getNewNode(1, c.freq, freq)
	}

	freq.items[k] = struct{}{}
	c.byKey[k] = newItem(v, freq)
}

// GetLFU gets the least frequently used key.
func (c *Cache[K, _]) GetLFU() (K, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var zero K
	if len(c.byKey) == 0 {
		return zero, false
	}

	for k := range c.freq.next.items {
		return k, true
	}
	panic("no element")
}

func (c *Cache[K, _]) deleteLFU() {
	var key K
	head := c.freq.next

	for k := range head.items {
		key = k
		break
	}

	if len(head.items) == 1 {
		c.freq.next = c.freq.next.next
	} else {
		delete(head.items, key)
	}

	delete(c.byKey, key)
}

func getNewNode[K comparable, V any](value int, prev, next *freqNode[K, V]) *freqNode[K, V] {
	n := newFrequencyNode[K, V]()
	n.value = value
	n.prev = prev
	n.next = next
	prev.next = n
	next.prev = n
	return n
}

func deleteNode[K comparable, V any](n *freqNode[K, V]) {
	n.prev.next = n.next
	n.next.prev = n.prev
}
