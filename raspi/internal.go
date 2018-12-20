package raspi

import (
	"io/ioutil"
	"os"
	"strconv"
)

func generateGPIODirectoryFilePath(pin uint8) string {
	return _GPIOClassPath + "/gpio" + pinUint8ToString(pin)
}

func generateGpioValueFilePath(pin uint8) string {
	return generateGPIODirectoryFilePath(pin) + "/value"
}

func generateGpioDirectionFilePath(pin uint8) string {
	return generateGPIODirectoryFilePath(pin) + "/direction"
}

const _GPIOExportFilePath = _GPIOClassPath + "/export"
const _GPIOUnexportFilePath = _GPIOClassPath + "/unexport"

func pinUint8ToString(pin uint8) string {
	return strconv.Itoa(int(pin))
}

func modeUint8ToString(mode uint8) (string, error) {
	switch mode {
	case INPUT:
		return "in", nil
	case OUTPUT:
		return "out", nil
	default:
		return emptyString, ErrInvalidPinMode
	}
}

func modeStringToUint8(mode string) (uint8, error) {
	switch mode {
	case "in\n":
		return INPUT, nil
	case "out\n":
		return OUTPUT, nil
	default:
		return emptyValue, ErrInvalidPinMode
	}
}

func valueStringToUint8(value string) (uint8, error) {
	tmp, err := strconv.Atoi(value)
	if err != nil {
		return emptyValue, err
	}
	switch uint8(tmp) {
	case LOW:
		return LOW, nil
	case HIGH:
		return HIGH, nil
	default:
		return emptyValue, ErrInvalidPinValue
	}
}

func valueUint8ToString(value uint8) (string, error) {
	switch uint8(value) {
	case LOW:
		return strconv.Itoa(LOW), nil
	case HIGH:
		return strconv.Itoa(HIGH), nil
	default:
		return emptyString, ErrInvalidPinValue
	}
}

// Only use these functions in a critical scope!!!
// https://www.kernel.org/doc/Documentation/gpio/sysfs.txt
func exportPin(pin uint8) error {
	exportFile, err := os.OpenFile(_GPIOExportFilePath, os.O_WRONLY, os.ModeType)
	defer func() {
		if exportFile != nil {
			exportFile.Close()
		}
	}()
	if err != nil {
		return err
	}
	_, err = exportFile.WriteString(pinUint8ToString(pin))
	if err != nil {
		return err
	}
	return nil
}

func unexportPin(pin uint8) error {
	unexportFile, err := os.OpenFile(_GPIOUnexportFilePath, os.O_WRONLY, os.ModeType)
	defer func() {
		if unexportFile != nil {
			unexportFile.Close()
		}
	}()
	if err != nil {
		return err
	}
	_, err = unexportFile.WriteString(pinUint8ToString(pin))
	if err != nil {
		return err
	}
	return nil
}

func digitalWrite(pin, value uint8) error {
	valueString, err := valueUint8ToString(value)
	if err != nil {
		return err
	}
	if v, _ := getPinMode(pin); v != OUTPUT {
		err = setPinMode(pin, OUTPUT)
		if v, _ = getPinMode(pin); v != OUTPUT {
			return err
		}
	}
	err = ioutil.WriteFile(generateGpioValueFilePath(pin), []byte(valueString), os.ModeType)
	if err != nil {
		return err
	}
	return nil
}

func digitalRead(pin uint8) (uint8, error) {
	data, err := ioutil.ReadFile(generateGpioValueFilePath(pin))
	if err != nil {
		return emptyValue, err
	}
	dataUint8, err := valueStringToUint8(string(data[:1]))
	if err != nil {
		return emptyValue, err
	}
	return dataUint8, nil
}

func setPinMode(pin uint8, mode uint8) error {
	modeString, err := modeUint8ToString(mode)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(generateGpioDirectionFilePath(pin), []byte(modeString), os.ModeType)
	if err != nil {
		return err
	}
	return nil
}

func getPinMode(pin uint8) (uint8, error) {
	data, err := ioutil.ReadFile(generateGpioDirectionFilePath(pin))
	if err != nil {
		return emptyMode, err
	}
	return modeStringToUint8(string(data))
}

func isPinExported(pin uint8) bool {
	_, err := os.Stat(generateGPIODirectoryFilePath(pin))
	return !os.IsNotExist(err)
}

func hasRightPermission() bool {
	if os.Getuid() != 0 {
		return false
	}
	return true
}
