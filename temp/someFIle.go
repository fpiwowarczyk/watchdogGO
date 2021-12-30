package main

import (
	"fmt"
	"os"
)

const (
	initDPAth          = "/etc/init.d/"
	serviceCmd         = "service"
	startCmd           = "start"
	statusCmd          = "status"
	serviceDown        = iota
	serviceStart       = iota
	serviceCannotStart = iota
)

func main() {
	stat, err := os.Stat("/etc/init.d/craa")
	if err != nil {
		fmt.Print("Error")
	}
	fmt.Printf(stat.Name())
}
