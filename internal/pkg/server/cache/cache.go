package cache

import (
	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/cache/scanner"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/cache/stat"
)

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() stat.Stat
	NewScanner() scanner.Scanner
}
