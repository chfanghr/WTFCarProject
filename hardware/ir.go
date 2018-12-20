package hardware

type IRData []uint8

const (
	IR_SendData CommandMethod = iota
)

const IR_DataMaxLen = 10

type IR interface {
	Send(IRData) error
}
