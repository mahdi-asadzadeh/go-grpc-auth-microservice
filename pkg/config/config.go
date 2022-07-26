package config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadSettings(develop bool) {
	var err error

	if develop == true {
		err = godotenv.Load("../pkg/config/envs/dev.env")

	} else {
		err = godotenv.Load("../pkg/config/envs/pro.env")
	}
	if err != nil {
		log.Fatal(err)
	}
}
