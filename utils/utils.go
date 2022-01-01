package utils

import (
	"os"

	"github.com/creamdog/gonfig"
)

func GetConfig(value string) (string, error) {
	file, err := os.Open("./config.json")
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
		return "", err
	}

	return output, nil
}
