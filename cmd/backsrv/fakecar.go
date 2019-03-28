package main

import (
	"flag"
	"github.com/chfanghr/WTFCarProject/car"
	"github.com/chfanghr/WTFCarProject/hardware"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type FakeCar struct {
	current            *location.Point2D
	lastmovementstatus int
	l                  *log.Logger
}

var InitialMessage = "Hello World!!"
var current string
var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func Worker(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if message != nil {
			err = c.WriteMessage(mt, []byte(InitialMessage))
			if err != nil {
				log.Println("write:", err)
				break
			}
		}

		//TODO get the current
		err = c.WriteMessage(mt, []byte(current))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func NewFakeCar(l *log.Logger) *FakeCar {
	var p *location.Point2D
	p.SetX(0)
	p.SetY(0)

	go func() {
		flag.Parse()
		http.HandleFunc("/Worker", Worker)
		log.Fatal(http.ListenAndServe(*addr, nil))

	}() //TODO Create a goroutine inside this constructor to provide websocket service and other staff
	return &FakeCar{
		current:            p,
		lastmovementstatus: car.Succeeded,
		l:                  l,
	}
}
func (f *FakeCar) GetLocation() (location.Point2D, error) {
	f.l.Println("FakeCar.GetLocation() called")
	return *f.current, nil
}
func (f *FakeCar) MoveTo(l location.Point2D) error {
	*f.current = l
	f.l.Println("FakeCar.MoveTo() called:", l)
	return nil
}
func (f *FakeCar) LastMovementStatus() int {
	f.l.Println("LastMovementStatus() called")
	return car.Succeeded
}
func (f *FakeCar) StopMovement() error {
	f.l.Println(" StopMovement() called")
	return nil
}
func (f *FakeCar) IRSend(ir hardware.IRData) error {
	f.l.Println("IRSend() called")
	return nil
}
