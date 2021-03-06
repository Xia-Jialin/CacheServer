package tcp

import (
	"net"

	"github.com/Xia-Jialin/CacheServer/internal/pkg/cache"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/cluster"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", s.Addr()+":12346")
	if e != nil {
		panic(e)
	}
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
