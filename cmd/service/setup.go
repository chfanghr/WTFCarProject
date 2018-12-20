package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chfanghr/backend/car"
	"github.com/chfanghr/backend/hardware"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

const (
	DefaultI2CAddr uint8 = 0x30
	DefaultI2CBus  int   = 0

	DefaultMotorAIN1 hardware.PinNumber = 3
	DefaultMotorAIN2 hardware.PinNumber = 5
	DefaultMotorBIN1 hardware.PinNumber = 6
	DefaultMotorBIN2 hardware.PinNumber = 9
	DefaultIRPin     hardware.PinNumber = 11

	DefaultIBeaconDeviceName string = "Car"
)

type Config struct {
	I2CAddr uint8 `json:"i2c_addr"`
	I2CBus  int   `json:"i2c_bus"`

	MotorAIN1 hardware.PinNumber `json:"motor_A1"`
	MotorAIN2 hardware.PinNumber `json:"motor_A2"`
	MotorBIN1 hardware.PinNumber `json:"motor_B1"`
	MotorBIN2 hardware.PinNumber `json:"motor_B2"`
	IRPin     hardware.PinNumber `json:"ir_pin"`

	SerialIBeacon string `json:"serial_ibeacon"`
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
		CleanUpFuncs.Add(func(i *log.Logger) {
			i.Println("remove pid file")
			os.Remove(pidFilePath)
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
	logger = log.New(w, serviceName+":", log.LstdFlags)
	logger.Println(serviceName, "starting")
	return logger, nil
}
func SetupListener(listenNetwork, listenAddress string) (net.Listener, error) {
	return net.Listen(listenNetwork, listenAddress)
}
func LoadCarService(configFile string) (car.Service, error) {
	if len(configFile) > 0 {
		buf, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}

		var c Config
		err = json.Unmarshal(buf, &c)
		if err != nil {
			return nil, err
		}
		ca, err := NewCar(c.I2CAddr, c.I2CBus, c.MotorAIN1, c.MotorAIN2, c.MotorBIN1, c.MotorBIN2, c.IRPin, c.SerialIBeacon)
		if err != nil {
			return nil, err
		}
		return car.NewGeneralServiceHandler(ca), nil
	} else {
		return LoadDefaultCarService()
	}
	return nil, errors.New("unknown error")
}
func LoadDefaultCarService() (car.Service, error) {
	c, err := NewCar(DefaultI2CAddr, DefaultI2CBus, DefaultMotorAIN1, DefaultMotorAIN2, DefaultMotorBIN1, DefaultMotorBIN2, DefaultIRPin, DefaultIBeaconDeviceName)
	if err != nil {
		return nil, err
	}
	return car.NewGeneralServiceHandler(c), nil
}
