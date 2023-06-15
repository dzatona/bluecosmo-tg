package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const envFileName = ".env"

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		log.Printf("[x] An error occurred while closing the file %s: %s\n", envFileName, err)
	}
}

func checkAndWriteDefaultEnv(file *os.File) {
	strs := []string{
		"BLUECOSMO_USERNAME=",
		"BLUECOSMO_PASSWORD=",
		"TG_AUTHORIZED_USERS=",
		"TG_BOT_TOKEN=",
	}
	for _, str := range strs {
		if _, err := file.WriteString(str + "\n"); err != nil {
			log.Fatalf("[x] An error occurred while writing to the file %s: %s", envFileName, err)
		}
	}
	if err := file.Sync(); err != nil {
		log.Fatalf("[x] An error occurred while syncing the file %s: %s", envFileName, err)
	}
	log.Printf("[x] Config written to file successfully. Please restart the program.")
	os.Exit(1)
}

func checkEnv() {
	log.Printf("[x] Checking .env file...")
	if _, err := os.Stat(envFileName); os.IsNotExist(err) {
		file, err := os.Create(envFileName)
		if err != nil {
			log.Fatalf("[x] An error occurred while creating the file %s: %s", envFileName, err)
		}
		defer closeFile(file)
		checkAndWriteDefaultEnv(file)
	} else {
		err := godotenv.Load(envFileName)
		if err != nil {
			log.Fatalf("[x] An error occurred while loading the environment variables from file %s: %s\n", envFileName, err)
		}
		strs := []string{
			"BLUECOSMO_USERNAME",
			"BLUECOSMO_PASSWORD",
			"TG_AUTHORIZED_USERS",
			"TG_BOT_TOKEN",
		}
		for _, str := range strs {
			value := os.Getenv(str)
			if value == "" {
				log.Fatalf("[x] Environment variable %s is not set. Please set it and try again.\n", str)
			}
		}
	}
}
