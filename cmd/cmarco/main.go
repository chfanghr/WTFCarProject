package main

import (
	"flag"
	"fmt"
	"github.com/chfanghr/Backend/hardware"
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
	strs = append(strs, "//This code is generated automatically by github.com/chfanghr/Backend/cmd/cmarco,do not edit")
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("OperationSucceeded", hardware.OperationSucceeded))
	strs = append(strs, CMarcoConstantValue("OperationFailed", hardware.OperationFailed))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("CommandGPIO", hardware.CommandGpio))
	strs = append(strs, CMarcoConstantValue("CommandIR", hardware.CommandIr))
	strs = append(strs, CMarcoConstantValue("CommandData", hardware.CommandData))
	strs = append(strs, CMarcoConstantValue("CommandSerial", hardware.CommandSerial))
	strs = append(strs, CMarcoConstantValue("CommandSerialHost", hardware.CommandSerialHost))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("GpioHigh", hardware.GpioHigh))
	strs = append(strs, CMarcoConstantValue("GpioLow", hardware.GpioLow))
	strs = append(strs, CMarcoConstantValue("GpioInput", hardware.GpioInput))
	strs = append(strs, CMarcoConstantValue("GpioInputPullUp", hardware.GpioInputPullUp))
	strs = append(strs, CMarcoConstantValue("GpioInputPullDown", hardware.GpioInputPullDown))
	strs = append(strs, CMarcoConstantValue("GpioOutput", hardware.GpioOutput))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("GpioPinmode", hardware.GpioPinmode))
	strs = append(strs, CMarcoConstantValue("GpioDigitalWrite", hardware.GpioDigitalWrite))
	strs = append(strs, CMarcoConstantValue("GpioDigitalRead", hardware.GpioDigitalRead))
	strs = append(strs, CMarcoConstantValue("GpioAnalogWrite", hardware.GpioAnalogWrite))
	strs = append(strs, CMarcoConstantValue("GpioAnalogRead", hardware.GpioAnalogRead))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("IRSendData", hardware.IrSendData))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("IrDatamaxlen", hardware.IrDatamaxlen))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("SerialHostWrite", hardware.SerialHostWrite))
	strs = append(strs, CMarcoConstantValue("SerialHostAvailable", hardware.SerialHostAvailable))
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
