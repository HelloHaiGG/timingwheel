package L1

import (
	"github.com/allegro/bigcache"
	"log"
	"time"
)

/**
big-cache
一级缓存存储 热点数据 基本不会改变的数据
*/

var lCache *bigcache.BigCache

func init() {
	cfg := bigcache.Config{
		Shards:             2 << 10,          //缓存分片数
		LifeWindow:         time.Second * 10, //超时时间
		CleanWindow:        2,                //2s 扫描一次
		MaxEntriesInWindow: 2 << 10,          //2M
		MaxEntrySize:       2,
		OnRemove:           onDelete, //删除回调
	}
	if c, err := bigcache.NewBigCache(cfg); err != nil {
		log.Fatal("Init big-cache error:", err)
	} else {
		lCache = c
	}
}

func Get(key string) (string, error) {
	b, err := lCache.Get(key)
	return string(b), err
}

func Put(key, value string) error {
	return Reset(key, value)
}

func Reset(key, value string) error {
	_ = lCache.Delete(key)
	return lCache.Set(key, []byte(value))
}

//删除缓存回调
func onDelete(key string, entry []byte) {
	log.Printf("Key:[%s] Value:[%s] Deleted.\n", key, string(entry))
}
