package main

import (
	"github.com/chfanghr/Backend/arduino"
	"github.com/chfanghr/Backend/generalserialhost"
	"github.com/chfanghr/Backend/hardware"
	"github.com/chfanghr/Backend/raspi"
	"log"
	"os"
)

const I2CAddr = 0x30
const I2CBus = 1

var logger *log.Logger

type RSSI int
type DevicesList map[string] /*address*/ RSSI
type Receiver struct {
	h hardware.SerialHost
	l DevicesList
}

func NewReceiver(h hardware.SerialHost) *Receiver {
	return &Receiver{
		h: h,
		l: make(DevicesList),
	}
}

func main() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
	i2cdev, err := raspi.NewI2C(I2CAddr, I2CBus)
	if err != nil {
		logger.Fatalln(err)
	}
	ardCtrl := arduino.NewArduinoControllerViaI2C(arduino.ArduinoProMini{}, i2cdev)
	s0 := generalserialhost.NewGeneralSerialHost(ardCtrl, 0)
	s1 := generalserialhost.NewGeneralSerialHost(ardCtrl, 1)
	s2 := generalserialhost.NewGeneralSerialHost(ardCtrl, 2)

	recvs := []*Receiver{NewReceiver(s0), NewReceiver(s1), NewReceiver(s2)}

	for _, v := range recvs {
		err := v.h.Begin(115200)
		if err != nil {
			logger.Fatalln(err)
		}
		err = v.h.SetTimeout(5000)
		if err != nil {
			logger.Fatalln(err)
		}
	}

	for _, v := range recvs {
		go func(receiver *Receiver) {

		}(v)
	}
}
