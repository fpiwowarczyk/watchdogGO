package utils

import (
	"errors"

	"github.com/creamdog/gonfig"
)

func GetConfig(value string, fs FileSystem) (string, error) {
	file, err := fs.Open("/home/filip/go/src/github.com/fpiwowarczyk/watchdogGO/config.json") // Can change to env
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

func Equals(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, el := range a {
		if el != b[i] {
			return false
		}
	}
	return true
}
