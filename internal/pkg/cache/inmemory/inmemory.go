package inmemory

import (
	"sync"
	"time"

	"github.com/Xia-Jialin/CacheServer/internal/pkg/cache/stat"
)

type value struct {
	v       []byte
	created time.Time
}

type inMemoryCache struct {
	c     map[string]value
	mutex sync.RWMutex
	stat.Stat
	ttl time.Duration
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = value{v, time.Now()}
	c.Add(k, v)
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k].v, nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.Stat.Del(k, v.v)
	}
	return nil
}

func (c *inMemoryCache) GetStat() stat.Stat {
	return c.Stat
}

//newInMemoryCache 创建基于内存储存数据的缓存服务
func NewInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{make(map[string]value), sync.RWMutex{}, stat.Stat{}, time.Duration(ttl) * time.Second}
	if ttl > 0 {
		go c.expirer()
	}
	return c
}

//expirer 删除过期数据
func (c *inMemoryCache) expirer() {
	for {
		time.Sleep(c.ttl)
		c.mutex.RLock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			if v.created.Add(c.ttl).Before(time.Now()) {
				c.Del(k)
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}
}
