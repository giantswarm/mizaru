package main

import "testing"

var Top = Topology{Hosts: []string{
	"host1",
	"host2",
	"host3",
	"host4",
	"host5",
}}

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
