package internal


var mapExcludeAddrs = map[string]bool{
	// TODO: Hardcode exclude address for disable receive funds
	//"aura19ad4tprcf9ew4577qph3jfzpf9slcrkpmxwvah": true,
	"aura15y8jhg5f5frp7g4uznxsj98fzyruuxaepalrzw": true, // aura-eco-growth
	"aura1kjgu375v4yqjffl0cfq4n6qw0r60w32k9nthh6": true, // aura-reserve
}

func MergeExcludeAddrs(m map[string]bool) map[string]bool {
	for k, v := range mapExcludeAddrs {
		m[k] = v
	}
	return m
}
