package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/juju/errgo"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func NewSSH(username string) *SSH {
	var a agent.Agent
	if con, e := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); e == nil {
		a = agent.NewClient(con)
	} else {
		panic("ssh-agent connect failed: " + e.Error())
	}

	signers, err := a.Signers()
	if err != nil {
		log.Fatal(err)
	}

	return &SSH{
		config: &ssh.ClientConfig{
			User: "username",
			Auth: []ssh.AuthMethod{ssh.PublicKeys(signers...)},
		},
	}
}

type SSH struct {
	config *ssh.ClientConfig
}

func (s *SSH) Do(host string, command []string) error {
	client, err := ssh.Dial("tcp", host+":22", s.config)
	if err != nil {
		return errgo.Mask(err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return errgo.Mask(err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		return errgo.Mask(err)
	}
	fmt.Println(b.String())
	return nil
}
