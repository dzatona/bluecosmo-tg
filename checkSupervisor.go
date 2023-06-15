package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func checkSupervisor() {
	log.Printf("[x] Checking supervisor...")
	cmd := "apt-cache policy supervisor | grep \"Installed: (none)\" | wc -l"
	supCheck, _ := exec.Command("bash", "-c", cmd).Output()
	if strings.Contains(string(supCheck), "1") {
		log.Printf("[x] Installing supervisor...")
		cmd = "apt-get -y install supervisor"
		_, _ = exec.Command("bash", "-c", cmd).Output()
	}
	log.Printf("[x] Checking supervisor config...")
	_, err := os.Stat("/etc/supervisor/conf.d/bluecosmo-tg.conf")
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[x] Writing supervisor config...")
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
			file, err := os.Create("/etc/supervisor/conf.d/bluecosmo-tg.conf")
			if err != nil {
				log.Fatalf("Failed to create file: %v", err)
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					log.Fatalf("Failed to close file: %v", err)
				}
			}(file)
			_, err = file.WriteString(configTemplate)
			if err != nil {
				log.Fatalf("Failed to write to file: %v", err)
			}
			log.Printf("[x] Rereading supervisor config...")
			cmd = "supervisorctl reread"
			_, _ = exec.Command("bash", "-c", cmd).Output()
			log.Printf("[x] Updating supervisor config...")
			cmd = "supervisorctl update"
			_, _ = exec.Command("bash", "-c", cmd).Output()
			log.Printf("[x] Starting bluecosmo-tg via supervisor...")
			cmd = "supervisorctl start bluecosmo-tg"
			_, _ = exec.Command("bash", "-c", cmd).Output()
			os.Exit(0)
		}
	}
}
