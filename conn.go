package dap

import (
	"bufio"
	"encoding/json"
	"io"
	"sync/atomic"
)

type Conn struct {
	io       io.ReadWriter
	handler  Handler
	seq      int64
	awaitMap map[int]chan []byte
}

func NewConn(io io.ReadWriter, handler Handler) *Conn {
	conn := &Conn{
		io:       io,
		handler:  handler,
		seq:      0,
		awaitMap: make(map[int]chan []byte),
	}
	return conn
}

// 处理请求和事件
type Handler interface {
	Handle(*Conn, Message)
}

// 阻塞等待结果
func (conn *Conn) SendRequest(request RequestMessage, response ResponseMessage) error {
	request.GetRequest().Seq = int(atomic.AddInt64(&conn.seq, 1))
	await := make(chan []byte, 1)
	conn.awaitMap[request.GetSeq()] = await
	if err := conn.Send(request); err != nil {
		close(await)
		delete(conn.awaitMap, request.GetSeq())
		return err
	}
	message := <-await
	close(await)
	return json.Unmarshal(message, response)
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
			if await, ok := conn.awaitMap[message.GetSeq()]; ok {
				await <- content
				delete(conn.awaitMap, message.GetSeq())
			}
		default:
			if conn.handler != nil {
				conn.handler.Handle(conn, message)
			}
		}
	}
}
