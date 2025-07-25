package pokecache

import (
	"testing"
	"time"
)



func TestAdd(t *testing.T) {
	cache := NewCache(8 * time.Second)

	cache.Add("test", []byte("test"))

	if _, found := cache.Get("test"); !found {
		t.Errorf("Expected to find key 'test'")
	}

	time.Sleep(5 * time.Second)
}

func TestReap(t *testing.T) {
	cache := NewCache(8 * time.Second)

	cache.Add("test", []byte("test"))

	time.Sleep(9 * time.Second)

	if _, found := cache.Get("test"); found {
		t.Errorf("Expected to not find key 'test'")
	}
}