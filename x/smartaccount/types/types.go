package types

var (
	SmartAccountI = "smart"
)

type SmartAccountValue struct {
	Type   string `json:"type"`
	Active bool   `json:"active"`
}
