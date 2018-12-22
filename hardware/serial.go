package hardware

import (
	"io"
)

const (
	SerialHostWrite CommandMethod = iota
	SerialHostAvailable
	SerialHostAvailableToWrite
	SerialHostBegin
	SerialHostEnd
	SerialHostRead
	SerialHostReadBytes
	SerialHostReadLine
	SerialHostReadString
	SerialHostReadUntil
	SerialHostSetTimeout
)

type LocalSerial interface {
	io.ReadWriter
	Flush() error
	Close() error
}

type SerialData []byte
type SerialTimeout uint16
type SerialBusIdentifier uint8
type SerialBaudRate int
type SerialHost interface {
	Begin(int) error
	End() error
	Available() (int, error)
	AvailableToWrite() (int, error)
	SetTimeout(int) error
	io.ReadWriter
	WriteAndRead(w, r []byte) (int, int, error)
}
