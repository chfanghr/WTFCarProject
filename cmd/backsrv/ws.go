package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"sync"
	"time"
)

const pingPeriod = time.Second
const cleanupPeriod = time.Second
const writeWait = time.Second

type wsConnection struct {
	*websocket.Conn
	closed chan struct{}
}

func newWsConnection(logger *log.Logger, conn *websocket.Conn) *wsConnection {
	res := &wsConnection{}
	res.Conn = conn
	res.closed = make(chan struct{})
	go res.worker(logger)
	return res
}
func (c *wsConnection) worker(logger *log.Logger) {
	if _, ok := <-c.closed; !ok {
		return
	}

	defer func() { _ = c.Close(); close(c.closed) }()

	for {
		<-time.AfterFunc(pingPeriod, func() {
			if err := c.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				logger.Println("ws ping error", err)
				return
			}
		}).C
	}
}

type wsService struct {
	*sync.Mutex
	conn []*wsConnection
}

func newWsService() *wsService {
	return &wsService{Mutex: new(sync.Mutex)}
}
func (w *wsService) AddConnection(c *wsConnection) {
	if c == nil {
		return
	}
	w.conn = append(w.conn, c)
	idx := len(w.conn) - 1
	go func() {
		<-c.closed
		if idx == 0 {
			w.conn = w.conn[1:]
			return
		}
		if idx == len(w.conn)-1 {
			w.conn = w.conn[:idx-1]
		}
		w.conn = append(w.conn[:idx-1], w.conn[idx:]...)
	}()
}
func (w *wsService) Update(data interface{}) error {
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var writers []io.Writer
	for _, c := range w.conn {
		w, err := c.NextWriter(websocket.TextMessage)
		if err != nil {
			continue
		}
		writers = append(writers, w)
	}
	writer := io.MultiWriter(writers...)
	_, _ = writer.Write([]byte(msg))
	for _, w := range writers {
		_ = w.(io.WriteCloser).Close()
	}
	return nil
}
