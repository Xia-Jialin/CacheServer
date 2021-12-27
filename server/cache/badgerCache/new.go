package badgerCache

import (
	"log"
	"time"

	"github.com/Xia-Jialin/CacheServer/server/cache/stat"
	"github.com/dgraph-io/badger"
)

type badgerCache struct {
	db  *badger.DB
	ttl int
	ch  chan *pair
}

type pair struct {
	k   string
	v   []byte
	ttl int
}

const BATCH_SIZE = 100

func NewbadgerCache(ttl int) *badgerCache {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan *pair, 5000)
	go write_func(db, c)

	return &badgerCache{db: db, ttl: ttl, ch: c}
}

func (c *badgerCache) Set(key string, value []byte) error {
	c.ch <- &pair{key, value, c.ttl}
	return nil
}

func (c *badgerCache) Get(k string) ([]byte, error) {
	var b []byte
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(k))
		if err != nil {
			return err
		}
		b, err = item.ValueCopy(nil)
		return err
	})
	return b, err
}

func (c *badgerCache) Del(k string) error {
	err := c.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(k))
		return err
	})
	return err
}
func (c *badgerCache) GetStat() stat.Stat {
	var count, KeySize, ValueSize int64
	err := c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			count++
			item := it.Item()
			k := item.Key()
			KeySize += int64(len(k))
			var b []byte
			b, _ = item.ValueCopy(nil)
			ValueSize += int64(len(b))
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return stat.Stat{
		Count:     count,
		KeySize:   KeySize,
		ValueSize: ValueSize,
	}
}

func write_func(db *badger.DB, c chan *pair) {
	t := time.NewTimer(time.Second)
	wb := db.NewWriteBatch()
	defer wb.Cancel()
	wb.SetMaxPendingTxns(BATCH_SIZE)
	for {
		select {
		case p := <-c:
			e := badger.NewEntry([]byte(p.k), p.v).WithTTL(time.Duration(p.ttl) * time.Second)
			wb.SetEntry(e)
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second)
		case <-t.C:
			wb.Flush()
			wb = db.NewWriteBatch()
			wb.SetMaxPendingTxns(BATCH_SIZE)
			t.Reset(time.Second)
		}
	}
}
