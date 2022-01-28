package badgerCache

import (
	"github.com/Xia-Jialin/CacheServer/internal/pkg/cache/scanner"
	"github.com/dgraph-io/badger"
)

type badgerScanner struct {
	pair
	pairCh  chan *pair
	closeCh chan struct{}
}

func (s *badgerScanner) Close() {
	close(s.closeCh)
}

func (s *badgerScanner) Scan() bool {
	p, ok := <-s.pairCh
	if ok {
		s.k, s.v = p.k, p.v
	}
	return ok
}

func (s *badgerScanner) Key() string {
	return s.k
}

func (s *badgerScanner) Value() []byte {
	return s.v
}

func (c *badgerCache) NewScanner() scanner.Scanner {
	pairCh := make(chan *pair)
	closeCh := make(chan struct{})
	go c.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			var v []byte
			v, _ = item.ValueCopy(nil)
			pairCh <- &pair{k: string(k), v: v}
			select {
			case <-closeCh:
				return nil
			case pairCh <- &pair{k: string(k), v: v}:
			}
		}
		return nil
	})
	return &badgerScanner{pair{}, pairCh, closeCh}
}
