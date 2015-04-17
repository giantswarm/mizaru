package main

const ChainInput = "INPUT"
const TargetDrop = "DROP"

func iptables(args ...string) []string {
	return append([]string{"iptables"}, args...)
}

type IPTables struct {
}

func (i *IPTables) Drop(from, to string) []string {
	return iptables("-A", ChainInput, "-s", from, "-d", to, "-j", TargetDrop)
}
func (i *IPTables) RevertDrop(from, to string) []string {
	return iptables("-D", ChainInput, "-s", from, "-d", to, "-j", TargetDrop)
}

func (i *IPTables) FlushRules() []string {
	return iptables("-F")
}
func (i *IPTables) CleanupChains() []string {
	return iptables("-X")
}
