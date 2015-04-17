package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

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
	serverAddr = pflag.String("server-addr", "", "Pass host:port to make it listen on http")
	hosts      = pflag.String("hosts", "", "Command separated list of hosts to work on")
)

func main() {
	pflag.Parse()

	srv := Service{executeLocally, &IPTables{}}

	if *serverAddr != "" {
		select {}
	} else {
		mode := pflag.Arg(0)
		if mode == "" {
			mode = "netsplit"
		}

		h := strings.Split(*hosts, ",")

		srv.Activate(mode, h, 10*time.Second)
	}
}
