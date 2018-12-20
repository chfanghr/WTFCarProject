package hardware

import "io"

type Serial interface {
	io.ReadWriter
	Flush() error
	Close() error
}
