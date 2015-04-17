package main

import (
	"reflect"
	"testing"
)

var Top = Topology{Hosts: []string{
	"host1",
	"host2",
	"host3",
	"host4",
	"host5",
}}

func hasRoute(routes []Rule, rule Rule) bool {
	for _, r := range routes {
		if r.From == rule.From && r.To == rule.To {
			return true
		}
	}
	return false
}

func TestRing(t *testing.T) {
	ring := Top.Ring()

	// 5 nodes that can talk to 2 other nodes each
	expectedRingSize := 5 * 2
	if len(ring) != expectedRingSize {
		t.Fatalf("Expected ring-size to be %d, but is %d", expectedRingSize, len(ring))
	}

	sees := map[string]int{}
	seen := map[string]int{}
	for _, rule := range ring {
		sees[rule.From]++
		seen[rule.To]++
	}

	for host, seenHosts := range sees {
		if seenHosts > 2 {
			t.Fatalf("Host %s can see %d hosts: %v", host, seenHosts, sees[host])
		}
	}
	for host, seenByHosts := range seen {
		if seenByHosts > 2 {
			t.Fatalf("Host %s can see %d hosts: %v", host, seenByHosts, seen[host])
		}
	}
}

func TestInvert(t *testing.T) {
	top := Topology{Hosts: []string{"host1", "host2"}}

	full := top.Default()
	inverted := top.Invert(full)

	if len(inverted) != 0 {
		t.Fatalf("Expected the invert of the Full() network to be empty, but is %v", inverted)
	}

	full2 := top.Invert(inverted)

	if !reflect.DeepEqual(full, full2) {
		t.Fatalf("Expected the Invert() of an empty list to be the Full() network, but is %v", full2)
	}
}

func TestInvert__SubNet(t *testing.T) {
	top := Topology{Hosts: []string{"host1", "host2", "host3"}}

	rules := []Rule{
		Rule{"host1", "host2"},
		Rule{"host2", "host1"},
	}

	inverted := top.Invert(rules)

	if hasRoute(inverted, rules[0]) {
		t.Fatalf("Expected rule host1->host2 to no longer exist")
	}
	if hasRoute(inverted, rules[1]) {
		t.Fatalf("Expected rule host1->host2 to no longer exist")
	}
	if !hasRoute(inverted, Rule{"host3", "host1"}) {
		t.Fatalf("Expected rule host3->host1 to exist")
	}
}
