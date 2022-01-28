package tcpnio

import (
	"bytes"
	"log"
	"strconv"

	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/cache"
	"github.com/Xia-Jialin/CacheServer/internal/pkg/server/cluster"
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

	index := bytes.IndexByte(data, byte('*'))
	if index == -1 {
		return
	}
	var argv [][]byte
	argc, _ := strconv.ParseUint(string(data[index+1:bytes.Index(data, []byte("\r\n"))]), 0, 64)
	data = data[bytes.Index(data, []byte("\r\n"))+2:]
	for i := 0; i < int(argc); i++ {
		index = bytes.IndexByte(data, byte('$'))
		length, _ := strconv.ParseUint(string(data[index+1:bytes.Index(data, []byte("\r\n"))]), 0, 64)
		data = data[bytes.Index(data, []byte("\r\n"))+2:]
		//argv = append(argv, data[bytes.Index(data, []byte("\r\n"))+2:bytes.Index(data, []byte("\r\n"))+2+int(length)])
		argv = append(argv, data[:length])
		data = data[int(length)+2:]
	}
	if len(argv) <= 0 {
		return
	}
	if bytes.Equal(bytes.ToLower(argv[0]), []byte("set")) && len(argv) >= 3 {
		es.Set(string(argv[1]), argv[2])
		c.AsyncWrite([]byte("+OK\r\n"))
	}
	if bytes.Equal(bytes.ToLower(argv[0]), []byte("get")) && len(argv) >= 2 {
		value, _ := es.Get(string(argv[1]))
		acc := append([]byte{}, ([]byte("$" + strconv.Itoa(len(value)) + "\r\n" + string(value) + "\r\n"))...)
		c.AsyncWrite(acc)
	}
	if bytes.Equal(bytes.ToLower(argv[0]), []byte("del")) && len(argv) >= 2 {
		es.Del(string(argv[1]))
		c.AsyncWrite([]byte("+OK\r\n"))
	}
	return
}
func (s *Server) Listen() {
	log.Fatal(gnet.Serve(s.Es, "tcp://"+s.Es.Addr()+":6379", gnet.WithMulticore(true)))
}

func New(c cache.Cache, n cluster.Node) *Server {
	p := goroutine.Default()
	defer p.Release()
	es := &echoServer{Cache: c, Node: n, pool: p}
	return &Server{es}
}
