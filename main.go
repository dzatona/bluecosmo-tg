package main

import (
	"log"
)

func main() {
	v := "2.0.4"
	log.Printf("[*] BLUECOSMO TELEGRAM BOT v.%s started.", v)
	checkDocker()
	checkEnv()
	checkSupervisor()
	initTelegram()
}
