package hardware

type Controller interface {
	Command(CommandRequest) (CommandResponse, error)
	IsValidPin(PinNumber) error
}
