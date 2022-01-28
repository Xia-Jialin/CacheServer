package main

import (
	"flag"
	"log"

	"github.com/Xia-Jialin/CacheServer/internal/pkg/cache"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/cluster"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/http"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/tcp"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/tcpnio"
)

func main() {
	typ := flag.String("type", "badger", "cache type")
	ttl := flag.Int("ttl", 30000, "cache time to live")
	node := flag.String("node", "0.0.0.0", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", *typ)
	log.Println("ttl is", *ttl)
	log.Println("node is", *node)
	log.Println("cluster is", *clus)
	c := cache.New(*typ, *ttl)
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	go tcp.New(c, n).Listen()
	go tcpnio.New(c, n).Listen()
	http.New(c, n).Listen()
}
