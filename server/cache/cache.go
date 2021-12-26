package cache

import (
	"github.com/Xia-Jialin/CacheServer/server/cache/scanner"
	"github.com/Xia-Jialin/CacheServer/server/cache/stat"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() stat.Stat
	NewScanner() scanner.Scanner
}
