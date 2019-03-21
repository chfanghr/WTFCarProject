package generalir

import . "github.com/chfanghr/WTFCarProject/hardware"

type IRRequester struct {
	data   IRData
	pin    PinNumber
	method CommandMethod
}

func NewIRRequester(p PinNumber, m CommandMethod, d IRData) *IRRequester {
	return &IRRequester{
		data:   d,
		pin:    p,
		method: m,
	}
}
func (i *IRRequester) Commit(c Controller) error {
	res, err := c.Command(*NewCommandRequest(CommandIr, i.method, i.data))
	if err != nil {
		return err
	}
	return res.Check(CommandIr)
}
func (IRRequester) GetRes() interface{} {
	return nil
}
