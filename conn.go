package dap

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"sync/atomic"
)

type Conn struct {
	io       io.ReadWriteCloser
	handler  Handler
	seq      int64
	awaitMap map[int]chan []byte
}

func NewConn(io io.ReadWriteCloser, handler Handler) *Conn {
	conn := &Conn{
		io:       io,
		handler:  handler,
		seq:      0,
		awaitMap: make(map[int]chan []byte),
	}
	return conn
}

type Handler interface {
	Handle(*Conn, Message)
}

// 阻塞等待结果
func (conn *Conn) SendRequest(request RequestMessage, response ResponseMessage) error {
	oldSeq := request.GetSeq()
	request.GetRequest().Seq = int(atomic.AddInt64(&conn.seq, 1))
	await := make(chan []byte, 1)
	conn.awaitMap[request.GetSeq()] = await
	defer func() {
		request.GetRequest().Seq = oldSeq
		close(await)
		delete(conn.awaitMap, request.GetSeq())
	}()
	if err := conn.Send(request); err != nil {
		return err
	}
	message := <-await
	if message == nil {
		return errors.New("conn close")
	}
	err := json.Unmarshal(message, response)
	if err == nil {
		response.GetResponse().RequestSeq = oldSeq
	}
	return err
}

// 非阻塞，不会拿到结果
func (conn *Conn) Send(message Message) error {
	return WriteProtocolMessage(conn.io, message)
}

func (conn *Conn) Run() {
	reader := bufio.NewReader(conn.io)
	for {
		content, err := ReadBaseMessage(reader)
		if err != nil {
			return
		}
		message, err := DecodeProtocolMessage(content)
		if err != nil {
			// TODO 协议解码错误
			continue
		}
		switch message := message.(type) {
		case ResponseMessage:
			seq := message.GetResponse().RequestSeq
			if await, ok := conn.awaitMap[seq]; ok {
				await <- content
				delete(conn.awaitMap, seq)
				continue
			}
		}
		if conn.handler != nil {
			conn.handler.Handle(conn, message)
		}
	}
	for k, v := range conn.awaitMap {
		v <- nil
		delete(conn.awaitMap, k)
	}
}

func (conn *Conn) Close() error {
	if conn.io != nil {
		return conn.io.Close()
	}
	return nil
}
