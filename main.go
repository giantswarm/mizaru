package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/ogier/pflag"
)

func executeLocally(command string, args []string) error {
	log.Printf("Executing: %s %v", command, args)
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

var (
	serverAddr     = pflag.String("server-addr", "", "Pass host:port to make it listen on http")
	timeout        = pflag.Int("timeout", 0, "Seconds to wait before reverting the rules (0=infinite wait)")
	dockerEndpoint = pflag.String("docker-endpoint", "unix:///var/run/docker.sock", "endpoint to use for docker communication")
	iptablesChain  = pflag.String("iptables-chain", "FORWARD", "iptables chain to apply rules to")
	block          = pflag.Bool("block", false, "Block until any signal arrives (which trigger rule cleanup)")
)

func main() {
	pflag.Parse()

	srv := Service{executeLocally, &IPTables{Chain: *iptablesChain}}

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

		var cancel chan struct{}
		if *block || *timeout > 0 {
			cancel = make(chan struct{})
			if *timeout > 0 {
				go func() {
					time.Sleep(time.Duration(*timeout) * time.Second)
					close(cancel)
				}()
			}

			if *block {
				signChan := make(chan os.Signal, 1)
				signal.Notify(signChan,
					os.Interrupt, os.Kill,
					syscall.SIGHUP,
					syscall.SIGINT,
					syscall.SIGTERM,
					syscall.SIGQUIT)
				go func() {
					<-signChan
					close(cancel)
				}()
			}
		}

		log.Println("Activating", mode)
		if err := srv.Activate(mode, hs, cancel); err != nil {
			log.Fatalf("Failed to active %s: %v", mode, err)
		}
	}
}
