package main

import (
	"log"
)

func main() {
	v := "1.1.1"
	log.Printf("[*] BLUECOSMO TELEGRAM BOT v.%s started.", v)
	checkEnv()
	checkSupervisor()
	initTelegram()
}
