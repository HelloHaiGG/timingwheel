package L1

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

/**
big-cache
*/

var bCache *bigcache.BigCache

func init() {
	cfg := bigcache.Config{
		Shards:             2 << 10, //缓存分片数
		LifeWindow:         time.Second * 60,
		CleanWindow:        3,
		MaxEntriesInWindow: 2 << 10,
		MaxEntrySize:       2,
		OnRemove:           nil, //删除回调
		OnRemoveWithReason: nil,
		Logger:             nil, //logger
	}
	if c, err := bigcache.NewBigCache(cfg); err != nil {
		log.Fatal("Init big-cache error:", err)
	} else {
		bCache = c
	}
}

func GetL1(key string) (string, error) {
	bCache.Get(key)
	return "", nil
}

func PutL1(key, value string) error {
	bCache.Set(key, []byte(value))
	return nil
}

func ResetL1(key, value string) error {
	return nil
}
