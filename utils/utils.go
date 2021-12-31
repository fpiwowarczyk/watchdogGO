package utils

import (
	"log"
	"os"

	"github.com/creamdog/gonfig"
)

func GetConfig(value string) string {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	config, err := gonfig.FromJson(file)
	if err != nil {
		log.Println(err)
	}
	output, err := config.GetString(value, nil)
	if err != nil {
		log.Println(err)
	}

	return output
}
