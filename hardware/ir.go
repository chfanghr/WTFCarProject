package hardware

type IRData []uint8

const (
	IrSendData CommandMethod = iota
)

const IrDatamaxlen = 10

type IR interface {
	Send(IRData) error
}
