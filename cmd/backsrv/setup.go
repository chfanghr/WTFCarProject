package main

import (
	"encoding/json"
	"fmt"
	"github.com/chfanghr/WTFCarProject/car"
	"github.com/chfanghr/WTFCarProject/hardware"
	"github.com/chfanghr/cleanuphandler"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

//const (
//	DefaultI2CAddr uint8 = 0x30
//	DefaultI2CBus  int   = 0
//
//	DefaultMotorAIN1 hardware.PinNumber = 3
//	DefaultMotorAIN2 hardware.PinNumber = 5
//	DefaultMotorBIN1 hardware.PinNumber = 6
//	DefaultMotorBIN2 hardware.PinNumber = 9
//	DefaultIRPin     hardware.PinNumber = 11
//
//	DefaultBluetoothHost string = "hci0"
//	DefaultIBeaconName   string = "Car"
//	DefaultIbeaconUUID   string = "c3468b29-38a3-4d80-a921-908450fcdd0e"
//)

type Config struct {
	I2CAddr uint8 `json:"i2c_addr"`
	I2CBus  int   `json:"i2c_bus"`

	MotorAIN1 hardware.PinNumber `json:"motor_A1"`
	MotorAIN2 hardware.PinNumber `json:"motor_A2"`
	MotorBIN1 hardware.PinNumber `json:"motor_B1"`
	MotorBIN2 hardware.PinNumber `json:"motor_B2"`
	IRPin     hardware.PinNumber `json:"ir_pin"`

	BluetoothHost string `json:"bluetooth_host"`
	IBeaconName   string `json:"iBeacon_name"`
	IBeaconUUID   string `json:"iBeacon_uuid"`
}

func SetupPidFile(pidFilePath string) error {
	if len(pidFilePath) > 0 {
		if _, err := os.Create(pidFilePath); err != nil {
			return err
		}
		err := ioutil.WriteFile(pidFilePath, []byte(string(fmt.Sprint(os.Getpid()))), 0666)
		if err != nil {
			return err
		}
		cleanuphandler.AddCleanupHandlers(func(i *log.Logger) {
			i.Println("remove pid file")
			_ = os.Remove(pidFilePath)
		})
		return nil
	}
	return nil
}

func SetupLogger(logFilePath string, useStdio bool) (logger *log.Logger, err error) {
	var logFile, stdio *os.File
	if len(logFilePath) > 0 {
		logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Fatalln("open log file :", err)
			return nil, err
		}
	}
	if useStdio {
		stdio = os.Stdout
	}
	var w io.Writer = nil
	if len(logFilePath) > 0 {
		if useStdio {
			w = io.MultiWriter(logFile, stdio)
		} else {
			w = logFile
		}
	} else {
		if useStdio {
			w = stdio
		} else {
			f, _ := os.Open("/dev/null")
			w = f
		}
	}
	logger = log.New(w, ServiceName+":", log.LstdFlags)
	logger.Println(ServiceName, "starting")
	return logger, nil
}

func SetupListener(listenNetwork, listenAddress string) (net.Listener, error) {
	return net.Listen(listenNetwork, listenAddress)
}

func LoadCarService(configFile string) (car.Service, error) {
	if len(configFile) > 0 {
		Logger.Println("load from config file")
		buf, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		var c Config
		err = json.Unmarshal(buf, &c)
		if err != nil {
			return nil, err
		}
		ca, err := NewCar(c.I2CAddr, c.I2CBus, c.MotorAIN1, c.MotorAIN2, c.MotorBIN1, c.MotorBIN2, c.IRPin, c.BluetoothHost, c.IBeaconName, c.IBeaconUUID)
		if err != nil {
			return nil, err
		}
		return car.NewGeneralServiceHandler(ca), nil
	} else {
		//return LoadDefaultCarService()
		Logger.Println("load fake car")
		return LoadFakeCarService()
	}
}

//func LoadDefaultCarService() (car.Service, error) {
//	c, err := NewCar(DefaultI2CAddr, DefaultI2CBus, DefaultMotorAIN1, DefaultMotorAIN2, DefaultMotorBIN1, DefaultMotorBIN2, DefaultIRPin, DefaultBluetoothHost, DefaultIBeaconName, DefaultIbeaconUUID)
//	if err != nil {
//		return nil, err
//	}
//	return car.NewGeneralServiceHandler(c), nil
//}

func LoadFakeCarService() (car.Service, error) {
	return car.NewGeneralServiceHandler(NewFakeCar(Logger, fca)), nil
}
