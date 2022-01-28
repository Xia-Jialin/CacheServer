package tcpnio

import (
	"bufio"
	"errors"
	"io"
)

//readKey 读取键
func (s *echoServer) readKey(r *bufio.Reader) (string, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", e
	}
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", e
	}
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", errors.New("redirect " + addr)
	}
	return key, nil
}

//readKeyAndValue 读取键和值
func (s *echoServer) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	klen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	vlen, e := readLen(r)
	if e != nil {
		return "", nil, e
	}
	k := make([]byte, klen)
	_, e = io.ReadFull(r, k)
	if e != nil {
		return "", nil, e
	}
	key := string(k)
	addr, ok := s.ShouldProcess(key)
	if !ok {
		return "", nil, errors.New("redirect " + addr)
	}
	v := make([]byte, vlen)
	_, e = io.ReadFull(r, v)
	if e != nil {
		return "", nil, e
	}
	return key, v, nil
}
