#!/usr/bin/env bash

mockgen_cmd="mockgen"

$mockgen_cmd -source=x/mint/types/expected_keepers.go -package testutil -destination x/mint/testutil/expected_keepers_mocks.go
