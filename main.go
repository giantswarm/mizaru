package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
)

func executeLocally(command string, args []string) error {
	log.Printf("Executing: %s %v", command, args)
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	defer log.Printf("done: %v", cmd.ProcessState)
	return cmd.Run()
}

var (
	serverAddr     = pflag.String("server-addr", "", "Pass host:port to make it listen on http")
	timeout        = pflag.Int("timeout", 10, "Seconds to wait before reverting the rules (0=no wait)")
	dockerEndpoint = pflag.String("docker-endpoint", "unix:///var/run/docker.sock", "endpoint to use for docker communication")
)

func main() {
	pflag.Parse()

	srv := Service{executeLocally, &IPTables{}}

	dockerClient, err := docker.NewClient(*dockerEndpoint)
	if err != nil {
		log.Fatalf("Could not conncet to docker host: %v", err)
	}
	dck := DockerInspector{client: dockerClient}

	if *serverAddr != "" {
		select {}
	} else {
		mode := pflag.Arg(0)
		if mode == "" {
			mode = "netsplit"
		}

		hs, err := dck.getContainers()
		if err != nil {
			log.Fatalf("Could not fetch docker containers: %v", err)
		}

		log.Println("Activating", mode)
		if err := srv.Activate(mode, hs, time.Duration(*timeout)*time.Second); err != nil {
			log.Fatalf("Failed to active %s: %v", mode, err)
		}
	}
}
