package cache

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -lbz2 -llz4 -ldl -lzstd -O3
import "C"
import (
	"unsafe"

	"github.com/Xia-Jialin/CacheServer/server/cache/scanner"
)

type rocksdbScanner struct {
	i           *C.rocksdb_iterator_t
	initialized bool
}

func (s *rocksdbScanner) Close() {
	C.rocksdb_iter_destroy(s.i)
}

func (s *rocksdbScanner) Scan() bool {
	if !s.initialized {
		C.rocksdb_iter_seek_to_first(s.i)
		s.initialized = true
	} else {
		C.rocksdb_iter_next(s.i)
	}
	return C.rocksdb_iter_valid(s.i) != 0
}

func (s *rocksdbScanner) Key() string {
	var length C.size_t
	k := C.rocksdb_iter_key(s.i, &length)
	return C.GoString(k)
}

func (s *rocksdbScanner) Value() []byte {
	var length C.size_t
	v := C.rocksdb_iter_value(s.i, &length)
	return C.GoBytes(unsafe.Pointer(v), C.int(length))
}

func (c *rocksdbCache) NewScanner() scanner.Scanner {
	return &rocksdbScanner{C.rocksdb_create_iterator(c.db, c.ro), false}
}
