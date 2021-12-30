package utils

import (
	"fmt"
	"os"

	"github.com/creamdog/gonfig"
)

func GetConfig(value string) string {
	file, err := os.Open("./config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	config, err := gonfig.FromJson(file)
	if err != nil {
		// TODO: error handle
	}
	watchdogTable, err := config.GetString(value, nil)
	if err != nil {
		fmt.Println(err)
	}

	return watchdogTable
}
