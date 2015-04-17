package main

import "math/rand"

type Rule struct {
	From string
	To   string
}

func fullMesh(hosts []string) []Rule {
	rules := []Rule{}

	for _, from := range hosts {
		for _, to := range hosts {
			if from == to {
				continue
			}

			rules = append(rules, Rule{from, to})
		}
	}
	return rules
}

// Topology is a helper to simulate weird netsplits.
type Topology struct {
	Hosts []string
}

func (t *Topology) randomHosts() []string {
	result := append([]string{}, t.Hosts...)

	for i := range result {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}
	return result
}

func (t *Topology) Default() []Rule {
	return fullMesh(t.Hosts)
}
func (t *Topology) NetSplit() []Rule {
	hosts := t.randomHosts()
	majorityN := len(hosts)/2 + 1

	majority := hosts[0:majorityN]
	minority := hosts[majorityN:]

	result := []Rule{}
	result = append(result, fullMesh(majority)...)
	result = append(result, fullMesh(minority)...)
	return result
}

func (t *Topology) Bridge() []Rule {
	hosts := t.randomHosts()
	majorityN := len(hosts)/2 + 1

	majority := hosts[0:majorityN]
	minority := hosts[majorityN:]

	result := []Rule{}
	result = append(result, fullMesh(majority)...)
	result = append(result, fullMesh(minority)...)

	// Now take one host from majority and make it see all minority hosts (bidirectional)
	bridge := majority[0]
	for _, host := range minority {
		result = append(result, Rule{bridge, host})
		result = append(result, Rule{host, bridge})
	}
	return result
}

// SingleBridge returns rules
func (t *Topology) SingleBridge() []Rule {
	hosts := t.randomHosts()
	majorityN := len(hosts)/2 + 1

	majority := hosts[0:majorityN]
	minority := hosts[majorityN:]

	result := []Rule{}
	result = append(result, fullMesh(majority)...)
	result = append(result, fullMesh(minority)...)

	// Now take one host from majority and make it see one minority host (bidirectional)
	leftNode := majority[0]
	rightNode := minority[0]
	result = append(result, Rule{leftNode, rightNode})
	result = append(result, Rule{rightNode, leftNode})

	return result
}

func (t *Topology) Ring() []Rule {
	// a ring only makes sense with at least 4 hosts.
	if len(t.Hosts) < 4 {
		return t.Default()
	}

	hosts := t.randomHosts()
	result := []Rule{}

	for i, from := range hosts {
		var sibling string
		if i+1 == len(hosts) {
			sibling = hosts[0]
		} else {
			sibling = hosts[i+1]
		}

		result = append(result, Rule{from, sibling})
		result = append(result, Rule{sibling, from})
	}
	return result
}

func (t *Topology) Invert(rules []Rule) []Rule {
	all := t.Default()
	newAll := []Rule{}

allRules:
	for _, allRule := range all {
		for _, rule := range rules {
			// If the allRule also exists in our "blacklist", we ignore it
			if allRule.From == rule.From && allRule.To == rule.To {
				continue allRules
			}
		}
		newAll = append(newAll, allRule)
	}
	return newAll
}
