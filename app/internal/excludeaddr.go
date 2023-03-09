package internal

import "fmt"

var mapExcludeAddrs = map[string]bool{
	// TODO: Hardcode exclude address for disable receive funds
	//"aura19ad4tprcf9ew4577qph3jfzpf9slcrkpmxwvah": true,
}

func MergeExcludeAddrs(m map[string]bool) map[string]bool {
	for k, v := range mapExcludeAddrs {
		m[k] = v
	}
	fmt.Printf("merged exclude addrs: %v", m)
	return m
}
