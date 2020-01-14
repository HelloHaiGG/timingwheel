package L1

import (
	"errors"
	"github.com/allegro/bigcache"
	"log"
	"sync"
	"time"
)

/**
big-cache
一级缓存存储 热点数据 基本不会改变的数据
*/

var bCache *bigcache.BigCache
var mFlag sync.Map

func init() {
	cfg := bigcache.Config{
		Shards:             2 << 10, //缓存分片数
		LifeWindow:         time.Second * 10,
		CleanWindow:        2,
		MaxEntriesInWindow: 2 << 10, //2M
		MaxEntrySize:       2,
		OnRemove:           onDelete, //删除回调
	}
	if c, err := bigcache.NewBigCache(cfg); err != nil {
		log.Fatal("Init big-cache error:", err)
	} else {
		bCache = c
	}
}

func Get(key string) (string, error) {
	var retry int
RETRY:
	if retry < 3 {
		if v, ok := mFlag.Load(key); ok && v.(bool) {
			retry++
			time.Sleep(time.Millisecond * 50 * 1 << retry) //100,200,400
			goto RETRY
		}
	} else {
		return "", errors.New("Key has be used. ")
	}
	b, err := bCache.Get(key)
	return string(b), err
}

func Put(key, value string) error {
	return Reset(key, value)
}

func Reset(key, value string) error {
	var retry int
RETRY:
	if retry < 3 {
		if v, ok := mFlag.Load(key); ok && v.(bool) {
			retry++
			time.Sleep(time.Millisecond * 50 * 1 << retry) //100,200,400
			goto RETRY
		}
	} else {
		return errors.New("Key has be used. ")
	}
	mFlag.Store(key, true)
	_ = bCache.Delete(key)
	if err := bCache.Set(key, []byte(value)); err != nil {
		return err
	}
	mFlag.Delete(key)
	return nil
}

//删除缓存回调
func onDelete(key string, entry []byte) {

}

//
func Demo(month int) int {
	if month <= 2 {
		return 1
	}
	monthN := 0
	monthN2 := 1
	monthN1 := 1
	for i := 0; i < month-2; i++ {
		monthN = monthN1 + monthN2
		monthN2 = monthN1
		monthN1 = monthN
	}
	return monthN
}
