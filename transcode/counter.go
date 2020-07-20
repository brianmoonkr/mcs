package transcode

import "sync/atomic"

type Counter uint64

// NewCounter ...
func NewCounter() *Counter {
	return new(Counter)
}

// Add ...
func (c *Counter) Add(val uint64) uint64 {
	return atomic.AddUint64((*uint64)(c), val)
}

// Up ...
func (c *Counter) Up() uint64 {
	return atomic.AddUint64((*uint64)(c), 1)
}

func (c *Counter) Down() uint64 {
	for v := c.Get(); v > 0; v = c.Get() {
		if atomic.CompareAndSwapUint64((*uint64)(c), v, v-1) {
			return v - 1
		}
	}
	return 0
}

func (c *Counter) Subtract(val uint64) uint64 {
	for v := c.Get(); (v - val) >= 0; v = c.Get() {
		if atomic.CompareAndSwapUint64((*uint64)(c), v, v-val) {
			return v - val
		}
	}
	return 0
}

func (c *Counter) Set(v uint64) {
	atomic.StoreUint64((*uint64)(c), v)
}

func (c *Counter) Get() uint64 {
	return atomic.LoadUint64((*uint64)(c))
}
