package main

import (
	"flag"
	"fmt"
	"github.com/chfanghr/WTFCarProject/hardware"
	"log"
	"os"
)

func CMarcoConstantValue(name string, value interface{}) string {
	return fmt.Sprint("#define ", name, " ", value)
}

func main() {
	file := flag.String("output", "", "output .h file")
	flag.Parse()
	var strs = []string{}
	strs = append(strs, "//This code is generated automatically by github.com/chfanghr/WTFCarProject/cmd/cmarco,do not edit")
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("Operation_Succeeded", hardware.OperationSucceeded))
	strs = append(strs, CMarcoConstantValue("Operation_Failed", hardware.OperationFailed))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("Command_GPIO", hardware.CommandGpio))
	strs = append(strs, CMarcoConstantValue("Command_IR", hardware.CommandIr))
	strs = append(strs, CMarcoConstantValue("Command_Data", hardware.CommandData))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("GPIO_HIGH", hardware.GpioHigh))
	strs = append(strs, CMarcoConstantValue("GPIO_LOW", hardware.GpioLow))
	strs = append(strs, CMarcoConstantValue("GPIO_INPUT", hardware.GpioInput))
	strs = append(strs, CMarcoConstantValue("GPIO_INPUT_PULLUP", hardware.GpioInputPullup))
	strs = append(strs, CMarcoConstantValue("GPIO_INPUT_PULLDOWN", hardware.GpioInputPulldown))
	strs = append(strs, CMarcoConstantValue("GPIO_OUTPUT", hardware.GpioOutput))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("GPIO_PinMode", hardware.GpioPinmode))
	strs = append(strs, CMarcoConstantValue("GPIO_DigitalWrite", hardware.GpioDigitalwrite))
	strs = append(strs, CMarcoConstantValue("GPIO_DigitalRead", hardware.GpioDigitalread))
	strs = append(strs, CMarcoConstantValue("GPIO_AnalogWrite", hardware.GpioAnalogwrite))
	strs = append(strs, CMarcoConstantValue("GPIO_AnalogRead", hardware.GpioAnalogread))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("IR_SendData", hardware.IrSenddata))
	if len(*file) > 0 {
		f, err := os.OpenFile(*file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatalln(err)
		}
		for _, v := range strs {
			_, err := f.WriteString(v + "\n")
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		for _, v := range strs {
			fmt.Println(v)
		}
	}
}
