package main

import (
	"fmt"
	"github.com/chfanghr/WTFCarProject/car"
	"github.com/chfanghr/WTFCarProject/hardware"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/chfanghr/WTFCarProject/rpcprotocal"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type FakeCar struct {
	mu      *sync.Mutex
	current *location.Point2D
	//lastMovementStatus int//currently not needed
	l        *log.Logger
	upgrader websocket.Upgrader
	wss      wsService
}
type message struct {
	CurrentLocation  *rpcprotocal.Point2D `json:"cur,omitempty"`
	DestineLocation  *rpcprotocal.Point2D `json:"dest,omitempty"`
	StopMovementFlag bool                 `json:"stop,omitempty"`
}

func NewFakeCar(l *log.Logger, listenAddr string) *FakeCar {
	res := &FakeCar{mu: new(sync.Mutex)}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := res.upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.Println("error upgrade ws", err)
			return
		}
		conn := newWsConnection(l, ws)
		res.mu.Lock()
		cur := *res.current
		res.mu.Unlock()
		_ = conn.WriteJSON(message{CurrentLocation: rpcprotocal.Point2DFromLocationPoint2D(cur)})
		res.wss.AddConnection(conn)
	})
	go l.Fatalln(http.ListenAndServe(listenAddr, http.DefaultServeMux))
	res.current = new(location.Point2D)
	return res
}
func (f *FakeCar) GetLocation() (location.Point2D, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return *f.current, nil
}
func (f *FakeCar) MoveTo(l location.Point2D) error {
	f.mu.Lock()
	f.current = &l
	f.mu.Unlock()
	return f.wss.Update(message{DestineLocation: rpcprotocal.Point2DFromLocationPoint2D(l)})
}
func (f *FakeCar) LastMovementStatus() int {
	//always success
	return car.Succeeded //FIXME
}
func (f *FakeCar) StopMovement() error {
	return f.wss.Update(message{StopMovementFlag: true})
}
func (f *FakeCar) IRSend(ir hardware.IRData) error {
	//need no implementation because there's no simulation
	return fmt.Errorf("not implemented")
}
