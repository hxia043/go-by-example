package lru

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func Test_CLRUCache_SyncMap(t *testing.T) {
	obj := NewCLRUCache2(2)
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

func Test_CLRUCache_BucketMap(t *testing.T) {
	obj := NewCLRUCache(2)
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

func BenchmarkGetAndPutWithSyncMap(b *testing.B) {
	b.ResetTimer()
	obj := NewCLRUCache2(128)
	wg := sync.WaitGroup{}
	wg.Add(b.N * 2)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			obj.Get(rand.Intn(200))
		}()
		go func() {
			defer wg.Done()
			obj.Put(rand.Intn(200), rand.Intn(200))
		}()
	}
	wg.Wait()
}

func BenchmarkGetAndPutWithLock(b *testing.B) {
	b.ResetTimer()
	obj := NewCLRUCache1(128)
	wg := sync.WaitGroup{}
	wg.Add(b.N * 2)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			obj.Get(rand.Intn(200))
		}()
		go func() {
			defer wg.Done()
			obj.Put(rand.Intn(200), rand.Intn(200))
		}()
	}
	wg.Wait()
}

func BenchmarkGetAndPutWithBucketMap(b *testing.B) {
	b.ResetTimer()
	obj := NewCLRUCache(128)
	wg := sync.WaitGroup{}
	wg.Add(b.N * 2)
	for i := 0; i < b.N; i++ {
		go func() {
			defer wg.Done()
			obj.Get(rand.Intn(200))
		}()
		go func() {
			defer wg.Done()
			obj.Put(rand.Intn(200), rand.Intn(200))
		}()
	}
	wg.Wait()
}
