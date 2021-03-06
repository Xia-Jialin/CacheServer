package cache

import (
	"log"

	"github.com/Xia-Jialin/CacheServer/internal/pkg/cache/badgerCache"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/cache/inmemory"
)

//New 根据typ参数的值选择储存数据的方式
//typ数据储存的方式
//ttl缓存生存时间
func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "inmemory" {
		c = inmemory.NewInMemoryCache(ttl)
	}
	// if typ == "rocksdb" {
	// 	c = newRocksdbCache(ttl)
	// }
	if typ == "badger" {
		x := badgerCache.NewbadgerCache(ttl)
		c = x
	}
	if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
