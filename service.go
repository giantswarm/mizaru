package main

import (
	"time"

	"github.com/juju/errgo"
)

type Executor func(command string, args []string) error

var (
	UnknownMode = errgo.Newf("Unknown mode.")
)

type Service struct {
	executor Executor
	iptables *IPTables
}

func (s *Service) Activate(mode string, hosts []string, timeout time.Duration) error {
	top := Topology{Hosts: hosts}

	var visibilityRules []Rule
	switch mode {
	case "fast":
		// Everything is normal
		visibilityRules = top.Default()

	case "netsplit":
		// The network is split into two groups. One has the majority of nodes
		visibilityRules = top.NetSplit()
	case "bridge":
		// The network is split into two groups. One has the majority, but one node can see both sides
		visibilityRules = top.Bridge()
	case "ring":
		// Each node can see two others nodes. They form a ring
		visibilityRules = top.Ring()
	default:
		return errgo.Notef(UnknownMode, "%s", mode)
	}

	dropRules := top.Invert(visibilityRules)

	for _, rule := range dropRules {
		cmd := s.iptables.Drop(rule.From, rule.To)

		if err := s.executor(cmd[0], cmd[1:]); err != nil {
			return errgo.Mask(err)
		}
	}

	if timeout > 0 {
		time.Sleep(timeout)
		for _, rule := range dropRules {
			cmd := s.iptables.RevertDrop(rule.From, rule.To)
			if err := s.executor(cmd[0], cmd[1:]); err != nil {
				return errgo.Mask(err)
			}
		}
	}

	return nil
}
