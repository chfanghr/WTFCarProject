package car

import (
	"github.com/chfanghr/Backend/hardware"
	"github.com/chfanghr/Backend/location"
	"net/rpc/jsonrpc"
)

type GeneralClient struct {
	networkType, netWorkAddr string
	serviceName              string
}

func NewGeneralClient(networkType string, netWorkAddr string, serviceName string) *GeneralClient {
	return &GeneralClient{
		networkType: networkType,
		netWorkAddr: netWorkAddr,
		serviceName: serviceName,
	}
}

func (g GeneralClient) GetLocation() (p location.Point2D, err error) {
	cl, err := jsonrpc.Dial(g.networkType, g.netWorkAddr)
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()
	if err != nil {
		return location.Point2D{}, err
	}
	err = cl.Call(g.serviceName+".GetLocation", 0, &p)
	return
}

func (g GeneralClient) MoveTo(l location.Point2D) (err error) {
	cl, err := jsonrpc.Dial(g.networkType, g.netWorkAddr)
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()
	if err != nil {
		return err
	}
	err = cl.Call(g.serviceName+".MoveTo", l, new(int))
	return
}

func (g *GeneralClient) LastMovementStatus() (r int) {
	cl, err := jsonrpc.Dial(g.networkType, g.netWorkAddr)
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()
	if err != nil {
		return Failed
	}
	err = cl.Call(g.serviceName+".LastMovementStatus", 0, &r)
	if err != nil {
		return Failed
	}
	return
}

func (g *GeneralClient) StopMovement() (err error) {
	cl, err := jsonrpc.Dial(g.networkType, g.netWorkAddr)
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()
	if err != nil {
		return
	}
	err = cl.Call(g.serviceName+".StopMovement", 0, new(int))
	if err != nil {
		return
	}
	return
}

func (g *GeneralClient) IRSend(d hardware.IRData) (err error) {
	cl, err := jsonrpc.Dial(g.networkType, g.netWorkAddr)
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()
	if err != nil {
		return
	}
	err = cl.Call(g.serviceName+".IRSend", d, new(int))
	return
}

func (g *GeneralClient) IsServiceAvailable() (err error) {
	cl, err := jsonrpc.Dial(g.networkType, g.netWorkAddr)
	defer func() {
		if cl != nil {
			cl.Close()
		}
	}()
	return
}
