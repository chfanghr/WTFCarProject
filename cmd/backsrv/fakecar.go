package main

import (
	"fmt"
	"github.com/chfanghr/WTFCarProject/car"
	"github.com/chfanghr/WTFCarProject/hardware"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
	"log"
	"net/http"
)

type FakeCar struct {
	current *location.Point2D
	//lastMovementStatus int//currently not needed
	l        *log.Logger
	upgrader websocket.Upgrader
	wss      wsService
}

func NewFakeCar(l *log.Logger) *FakeCar {
	res := &FakeCar{}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws, err := res.upgrader.Upgrade(w, r, nil)
		if err != nil {
			l.Println("error upgrade ws", err)
			return
		}
		res.wss.AddConnection(newWsConnection(l, ws))
	})
	res.current = new(location.Point2D)
	return res
}
func (f *FakeCar) GetLocation() (location.Point2D, error) {
	return *f.current, nil
}
func (f *FakeCar) MoveTo(l location.Point2D) error {
	f.l.Println("FakeCar.MoveTo() called:", l) //TODO implement this method
	return nil
}
func (f *FakeCar) LastMovementStatus() int {
	//always success
	return car.Succeeded //FIXME
}
func (f *FakeCar) StopMovement() error {
	return fmt.Errorf("not implemented") //TODO implement this method
}
func (f *FakeCar) IRSend(ir hardware.IRData) error {
	//need no implementation because there's no simulation
	return fmt.Errorf("not implemented")
}
