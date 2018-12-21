package hardware

import (
	"errors"
)

const (
	Operation_Succeeded CommandStatus = iota
	Operation_Failed
)
const (
	Command_GPIO CommandType = iota
	Command_IRRemote
	Command_Data
)

type CommandType int
type CommandStatus int
type CommandMethod int
type CommandParameter interface{}
type CommandRequest struct {
	CommType      CommandType        `json:"type"`
	CommMethod    CommandMethod      `json:"method"`
	CommParameter []CommandParameter `json:"param,omitempty"`
}
type CommandResponse struct {
	CommType      CommandType        `json:"type"`
	CommStatus    CommandStatus      `json:"status"`
	CommParameter []CommandParameter `json:"param,omitempty"`
}
type Requester interface {
	Commit(c Controller) error
	GetRes() interface{}
}

func NewCommandRequest(t CommandType, m CommandMethod, p ...CommandParameter) *CommandRequest {
	return &CommandRequest{
		CommType:      t,
		CommMethod:    m,
		CommParameter: append([]CommandParameter{}, p...),
	}
}
func NewCommandResponse(t CommandType, s CommandStatus, p ...CommandParameter) *CommandResponse {
	return &CommandResponse{
		CommType:      t,
		CommStatus:    s,
		CommParameter: append([]CommandParameter{}, p...),
	}
}
func (res CommandResponse) Check(t CommandType) error {
	if res.CommType != t {
		return errors.New("type error")
	}
	switch res.CommStatus {
	case Operation_Succeeded:
		return nil
	case Operation_Failed:
		return errors.New("operation failed")
	default:
		return errors.New("unknown error")
	}
}
func (res CommandResponse) GetParameter() []CommandParameter {
	return res.CommParameter
}
