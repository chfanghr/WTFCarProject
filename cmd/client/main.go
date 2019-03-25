package main

import (
	"fmt"
	"github.com/chfanghr/WTFCarProject/car"
)

const (
	networkType = "tcp"
	networkAddr = "localhost:8888"
	serviceName = "backendService"
)

func main() {
	client := car.NewGeneralClient(networkType, networkAddr, serviceName)
	fmt.Println(client.GetLocation())
}
