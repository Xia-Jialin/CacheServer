package tcpnio

import (
	"fmt"
	"log"

	"github.com/Xia-Jialin/CacheServer/server/cache"
	"github.com/Xia-Jialin/CacheServer/server/cluster"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
)

type Server struct {
	Es *echoServer
}

type echoServer struct {
	cluster.Node
	cache.Cache
	*gnet.EventServer
	pool *goroutine.Pool
}

func (es *echoServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	data := append([]byte{}, frame...)

	com := bytes.Split(data, []byte("\r\n"))
	op := data[0]
	if op == '*' {
		vlen := fmt.Sprintf("%d ", len(""))
		data = []byte(vlen)
		c.AsyncWrite(data)
	}
	// Use ants pool to unblock the event-loop.
	// _ = es.pool.Submit(func() {
	// })

	return
}
func (s *Server) Listen() {
	log.Fatal(gnet.Serve(s.Es, "tcp://"+s.Es.Addr()+":12346", gnet.WithMulticore(true)))
}

func New(c cache.Cache, n cluster.Node) *Server {
	p := goroutine.Default()
	defer p.Release()
	es := &echoServer{Cache: c, Node: n, pool: p}
	return &Server{es}
}
