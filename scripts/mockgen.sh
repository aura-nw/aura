#!/usr/bin/env bash

mockgen_cmd="mockgen"

$mockgen_cmd -source=x/mint/types/expected_keepers.go -package testutil -destination tests/mocks/mint/expected_keepers_mocks.go
$mockgen_cmd -source=x/bank/types/expected_keepers.go -package testutil -destination tests/mocks/bank/expected_keepers_mocks.go
$mockgen_cmd -source=x/auth/vesting/types/expected_keepers.go -package testutil -destination tests/mocks/auth/expected_keepers_mocks.go
