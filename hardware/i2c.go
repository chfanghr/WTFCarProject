package hardware

import "io"

type I2C interface {
	ReadBytes(buf []byte) (int, error)
	WriteBytes(buf []byte) (int, error)
	io.ReadWriter
	Close() error
}
