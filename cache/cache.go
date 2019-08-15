package cache

import (
	"fmt"

	"github.com/mediocregopher/radix"
)

// Cache represents a cache for counter values
type Cache struct {
	client radix.Client
}

func NewCache(conUrl string) (*Cache, error) {
	pool, err := radix.NewPool("tcp", conUrl, 10)
	if err != nil {
		fmt.Printf("error while connecting to redis: %s", err)
		return &Cache{}, err
	}
	return &Cache{
		client: pool,
	}, nil
}

func (cache *Cache) Close() {
	cache.client.Close()
}

func (cache *Cache) CreateIfNotExist(key string, value []byte) {
	mvalue := map[string]int{"a": 1, "b": 2, "c": 3}
	err := cache.client.Do(radix.FlatCmd(nil, "MSETNX", key, mvalue))
	if err != nil {
		fmt.Printf("error: %+v", err)
	}
}

func (cache *Cache) IncreaseBy(key string, value uint64) {
	err := cache.client.Do(radix.Cmd(nil, "INCRBY", key, fmt.Sprintf("%d", value)))
	if err != nil {
		fmt.Printf("error while increasing value: %s\n", err)
	}
}

func (cache *Cache) Delete(key string) {
	err := cache.client.Do(radix.Cmd(nil, "DEL", key))
	if err != nil {
		fmt.Printf("error while deleting value: %s\n", err)
	}
}

func (cache *Cache) FindKeys(match string) []string {
	scanner := radix.NewScanner(cache.client, radix.ScanOpts{Command: "SCAN", Pattern: match})
	defer scanner.Close()
	var key string
	keys := make([]string, 0)
	for scanner.Next(&key) {
		keys = append(keys, key)
	}
	return keys
}

func (cache *Cache) Get(key string) []byte {
	var val []byte
	err := cache.client.Do(radix.Cmd(&val, "GET", key))
	if err != nil {
		fmt.Printf("error while getting cache value %s: %s", key, err)
	}
	return val
}
