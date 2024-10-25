package lru

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_LRUCache(t *testing.T) {
	obj := NewLRUCache(2)
	time.Sleep(150 * time.Millisecond)
	obj.Put(1, 1)
	time.Sleep(150 * time.Millisecond)
	obj.Put(2, 2)
	time.Sleep(150 * time.Millisecond)
	param1 := obj.Get(1)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, true, param1)
	obj.Put(3, 3)
	time.Sleep(150 * time.Millisecond)
	param1 = obj.Get(2)
	assert.Equal(t, false, param1)
	obj.Put(4, 4)
	time.Sleep(150 * time.Millisecond)
	param1 = obj.Get(1)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, false, param1)
	param1 = obj.Get(3)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, true, param1)
	param1 = obj.Get(4)
	time.Sleep(150 * time.Millisecond)
	assert.Equal(t, true, param1)
}
