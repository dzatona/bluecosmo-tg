package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func checkDocker() {
	log.Printf("[x] Checking docker...")
	cmd := exec.Command("docker", "--version")
	err := cmd.Run()
	if err != nil {
		log.Printf("[x] Installing docker...")
		cmds := []string{
			"sudo apt-get update -y",
			"sudo apt-get install -y ca-certificates curl gnupg",
			"sudo install -m 0755 -d /etc/apt/keyrings",
			"curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg",
			"sudo chmod a+r /etc/apt/keyrings/docker.gpg",
		}
		for _, cmd := range cmds {
			_, err = exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				log.Printf("[*] Error: %s", err)
			}
		}
		out, err := exec.Command("bash", "-c", "dpkg --print-architecture").Output()
		if err != nil {
			log.Printf("[*] Error: %s", err)
			return
		}
		arch := strings.TrimSpace(string(out))
		out, err = exec.Command("bash", "-c", ". /etc/os-release && echo $VERSION_CODENAME").Output()
		if err != nil {
			log.Printf("[*] Error: %s", err)
			return
		}
		codename := strings.TrimSpace(string(out))
		dockerList := fmt.Sprintf("deb [arch=%s signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu %s stable", arch, codename)
		err = os.WriteFile("/etc/apt/sources.list.d/docker.list", []byte(dockerList), 0644)
		if err != nil {
			log.Printf("[*] Error: %s", err)
			return
		}
		cmds = []string{
			"sudo apt-get update -y",
			"sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y",
			"sudo docker run hello-world",
		}
		for _, cmd := range cmds {
			out, err = exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				log.Printf("[*] Error: %s", err)
			} else {
				log.Printf("[x] Docker: %s", out)
			}
		}
	}
	command := "docker pull alpeware/chrome-headless-trunk"
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		log.Printf("[*] Error: %s", err)
	} else {
		log.Printf("[x] Docker: %s", out)
	}
}

func isDockerContainerRunning(containerName string) (bool, error) {
	cmd := exec.Command("docker", "ps", "-q", "-f", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("error executing command 'docker ps': %w", err)
	}
	if strings.TrimSpace(string(output)) == "" {
		return false, nil
	}
	return true, nil
}

func checkContainer() {
	isRunning, err := isDockerContainerRunning("headless-shell")
	if err != nil {
		log.Fatalf("[*] Error checking if Docker container is running: %v", err)
	}
	if !isRunning {
		log.Println("[x] Docker: container is not running. Starting...")
		command := "docker run -d -p 9222:9222 --rm --name headless-shell alpeware/chrome-headless-trunk"
		_, _ = exec.Command("bash", "-c", command).Output()
	} else {
		cmds := []string{
			"docker stop headless-shell",
			"docker run -d -p 9222:9222 --rm --name headless-shell alpeware/chrome-headless-trunk",
		}
		for _, cmd := range cmds {
			out, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				log.Printf("[*] Error: %s", err)
				os.Exit(1)
			} else {
				log.Printf("[x] Docker: %s", out)
			}
		}
	}
	time.Sleep(5 * time.Second)
}
