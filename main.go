package main

import "time"

func main() {
	hosts := []string{
		"172.17.8.101",
		"172.17.8.102",
		"172.17.8.103",
	}

	srv := Service{NewSSH("core"), &IPTables{}}
	srv.Activate("netsplit", hosts, 10*time.Second)

	select {}
}
