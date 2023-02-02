#!/usr/bin/env bash

mockgen_cmd="mockgen"

$mockgen_cmd -source=x/mint/types/expected_keepers.go -package testutil -destination x/mint/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/bank/types/expected_keepers.go -package testutil -destination x/bank/testutil/expected_keepers_mocks.go
$mockgen_cmd -source=x/auth/vesting/types/expected_keepers.go -package testutil -destination x/auth/testutil/expected_keepers_mocks.go
