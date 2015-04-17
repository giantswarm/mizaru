package main

const TargetDrop = "DROP"

func iptables(args ...string) []string {
	return append([]string{"iptables", "--wait"}, args...)
}

type IPTables struct {
	Chain string
}

func (i *IPTables) Drop(from, to string) []string {
	return iptables("-A", i.Chain, "-s", from, "-d", to, "-j", TargetDrop)
}
func (i *IPTables) RevertDrop(from, to string) []string {
	return iptables("-D", i.Chain, "-s", from, "-d", to, "-j", TargetDrop)
}

func (i *IPTables) FlushRules() []string {
	return iptables("-F")
}
func (i *IPTables) CleanupChains() []string {
	return iptables("-X")
}
