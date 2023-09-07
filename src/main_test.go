package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const N = 100

func TestPool(t *testing.T) {
	poolsize := 20
	pool := NewPool(poolsize)
	for i := 0; i < poolsize; i++ {
		pool.Schedule(func() {})
	}

	values := make([]int, 0)
	for i := 0; i < N; i++ {
		val := i
		pool.Schedule(func() {
			values = append(values, val+1)
		})
	}

	// Imitate some other stuffs.
	pool.Schedule(func() {})
	pool.Schedule(func() {})
	pool.Schedule(func() {})

	sum := 0
	for i := 0; i < len(values); i++ {
		sum += values[i]
	}
	t.Log("values:", values, "len:", len(values))
	assert.Equal(t, sum, (N)*(N+1)/2)

}
