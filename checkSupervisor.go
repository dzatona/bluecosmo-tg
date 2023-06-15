package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	bash      = "bash"
	cFlag     = "-c"
	supConfig = "/etc/supervisor/conf.d/bluecosmo-tg.conf"
)

func checkSupervisor() {
	log.Println("[x] Checking supervisor...")
	cmd := "apt-cache policy supervisor | grep \"Installed: (none)\" | wc -l"
	supCheck, err := exec.Command(bash, cFlag, cmd).Output()
	if err != nil {
		log.Fatalf("[x] Failed to check supervisor: %v", err)
	}
	if strings.Contains(string(supCheck), "1") {
		log.Println("[x] Installing supervisor...")
		cmd = "apt-get -y install supervisor"
		if _, err := exec.Command(bash, cFlag, cmd).Output(); err != nil {
			log.Fatalf("[x] Failed to install supervisor: %v", err)
		}
	}
	log.Println("[x] Checking supervisor config...")
	if _, err := os.Stat(supConfig); os.IsNotExist(err) {
		writeSupervisorConfig()
		updateSupervisorConfig()
		startBluecosmoTg()
		os.Exit(0)
	} else {
		log.Println("[x] Supervisor config exists, run: supervisorctl start bluecosmo-tg")
	}
}

func writeSupervisorConfig() {
	log.Println("[x] Writing supervisor config...")
	const configTemplate = `[program:bluecosmo-tg]
command=/etc/bluecosmo/bc
directory=/etc/bluecosmo
autostart=true
autorestart=true
user=root
stdout_logfile=/var/log/bluecosmo-tg.log
redirect_stderr=true
startsecs=0
numprocs=1
`
	file, err := os.Create(supConfig)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close file: %v", err)
		}
	}(file)
	if _, err := file.WriteString(configTemplate); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}

func updateSupervisorConfig() {
	log.Println("[x] Rereading supervisor config...")
	if _, err := exec.Command(bash, cFlag, "supervisorctl reread").Output(); err != nil {
		log.Fatalf("[x] Failed to reread supervisor config: %v", err)
	}
	log.Println("[x] Updating supervisor config...")
	if _, err := exec.Command(bash, cFlag, "supervisorctl update").Output(); err != nil {
		log.Fatalf("[x] Failed to update supervisor config: %v", err)
	}
}

func startBluecosmoTg() {
	log.Println("[x] Starting bluecosmo-tg via supervisor...")
	if _, err := exec.Command(bash, cFlag, "supervisorctl start bluecosmo-tg").Output(); err != nil {
		log.Fatalf("[x] Failed to start bluecosmo-tg: %v", err)
	}
}
