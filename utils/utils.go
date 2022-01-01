package utils

import (
	"errors"
	"os"

	"github.com/creamdog/gonfig"
)

func GetConfig(value string) (string, error) {
	file, err := os.Open("/home/filip/go/src/github.com/fpiwowarczyk/watchdogGO/config.json")
	if err != nil {
		return "", err
	}
	defer file.Close()

	config, err := gonfig.FromJson(file)
	if err != nil {
		return "", err
	}
	output, err := config.GetString(value, nil)
	if err != nil {
		return "", errors.New("Missing property: " + value)
	}

	return output, nil
}
