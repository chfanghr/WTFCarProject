package main

import (
	"fmt"
	"github.com/chfanghr/Backend/hardware"
)

func CMarcoConstantValue(name string, value interface{}) string {
	return fmt.Sprint("#define ", name, " ", value)
}

func main() {
	var strs = []string{}
	strs = append(strs, "//This code is generated automatically,do not edit")
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("Operation_Succeeded", hardware.Operation_Succeeded))
	strs = append(strs, CMarcoConstantValue("Operation_Failed", hardware.Operation_Failed))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("Command_GPIO", hardware.Command_GPIO))
	strs = append(strs, CMarcoConstantValue("Command_IRRemote", hardware.Command_IRRemote))
	strs = append(strs, CMarcoConstantValue("Command_Data", hardware.Command_Data))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("GPIO_HIGH", hardware.GPIO_HIGH))
	strs = append(strs, CMarcoConstantValue("GPIO_HIGH", hardware.GPIO_HIGH))
	strs = append(strs, CMarcoConstantValue("GPIO_LOW", hardware.GPIO_LOW))
	strs = append(strs, CMarcoConstantValue("GPIO_INPUT", hardware.GPIO_INPUT))
	strs = append(strs, CMarcoConstantValue("GPIO_INPUT_PULLUP", hardware.GPIO_INPUT_PULLUP))
	strs = append(strs, CMarcoConstantValue("GPIO_INPUT_PULLDOWN", hardware.GPIO_INPUT_PULLDOWN))
	strs = append(strs, CMarcoConstantValue("GPIO_OUTPUT", hardware.GPIO_OUTPUT))
	strs = append(strs, "//------------------------------------")
	strs = append(strs, CMarcoConstantValue("GPIO_PinMode", hardware.GPIO_PinMode))
	strs = append(strs, CMarcoConstantValue("GPIO_DigitalWrite", hardware.GPIO_DigitalWrite))
	strs = append(strs, CMarcoConstantValue("GPIO_DigitalRead", hardware.GPIO_DigitalRead))
	strs = append(strs, CMarcoConstantValue("GPIO_AnalogWrite", hardware.GPIO_AnalogWrite))
	strs = append(strs, CMarcoConstantValue("GPIO_AnalogRead", hardware.GPIO_AnalogRead))
	for _, v := range strs {
		fmt.Println(v)
	}
}
