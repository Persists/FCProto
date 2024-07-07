package main

import (
	"fmt"

	"github.com/Persists/fcproto/pkg/sensors"
)

func main() {
	cpu := sensors.NewCpuSensor()

	data := cpu.GenerateData()

	fmt.Println(data.ToString())

	mem := sensors.NewMemSensor()

	data = mem.GenerateData()

	fmt.Println(data.ToString())

	virt := sensors.NewVirtualSensor()

	data = virt.GenerateData()

	fmt.Println(data.ToString())

}
