package generalir

import . "github.com/chfanghr/Backend/hardware"

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
	res, err := c.Command(*NewCommandRequest(Command_IRRemote, i.method, i.data))
	if err != nil {
		return err
	}
	return res.Check(Command_IRRemote)
}
func (IRRequester) GetRes() interface{} {
	return nil
}
