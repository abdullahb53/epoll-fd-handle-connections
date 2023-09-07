package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	poolsize := 100
	pool := NewPool(poolsize)
	for i := 0; i < poolsize; i++ {
		pool.Schedule(func() {})
	}

	values := make([]int, 0)
	for i := 0; i < poolsize; i++ {
		pool.Schedule(func() {
			values = append(values, i)
		})
	}
	assert.Equal(t, poolsize-1, len(values))
}
