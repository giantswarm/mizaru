package main

import (
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

func (s *Service) Activate(mode string, hosts []Container, cancel chan struct{}) error {
	var hs []string
	for _, host := range hosts {
		hs = append(hs, host.IP)
	}
	top := Topology{Hosts: hs}

	var visibilityRules []Rule
	switch mode {
	case "fast":
		fallthrough
	case "cleanup":
		fallthrough
	case "mr-clean":
		// Everything is normal
		visibilityRules = top.Default()

	case "netsplit":
		// The network is split into two groups. One has the majority of nodes
		visibilityRules = top.NetSplit()
	case "bridge":
		// The network is split into two groups. One has the majority, but one node can see both sides
		visibilityRules = top.Bridge()

	case "single-bridge":
		fallthrough
	case "singlebridge":
		// Like bridge, but the bridge node can only see one node of the minority
		visibilityRules = top.SingleBridge()

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

	if cancel != nil {
		<-cancel // block until cancel is closed or sthg is sent
		for _, rule := range dropRules {
			cmd := s.iptables.RevertDrop(rule.From, rule.To)
			if err := s.executor(cmd[0], cmd[1:]); err != nil {
				return errgo.Mask(err)
			}
		}
	}

	return nil
}
