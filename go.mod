module github.com/aura-nw/aura

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.45.1
	github.com/cosmos/ibc-go/v2 v2.2.0
	github.com/gogo/protobuf v1.3.3
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	// github.com/tendermint/spm v0.1.8
	github.com/tendermint/starport v0.19.2
	github.com/tendermint/tendermint v0.34.14
	github.com/tendermint/tm-db v0.6.6
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	google.golang.org/genproto v0.0.0-20220302033224-9aa15565e42a
	google.golang.org/grpc v1.44.0
)

replace (
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.7-0.20210622111912-ef00f8ac3d76
	github.com/cosmos/ibc-go => github.com/cosmos/ibc-go/v2 v2.0.3
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
