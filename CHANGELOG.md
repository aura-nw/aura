<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking CLI commands and REST routes used by end-users.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.

Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## Unreleased

## [v0.6.1] - 2023-07-07

### Improvements
- Add config swagger
- Change before check tx of smart account from query to execution

## [v0.6.0] - 2023-06-30

### Features
- Support module SmartAccount

## [v0.5.2] - 2023-06-19

### Improvements
- Fix bug create vesting account

## [v0.5.1] - 2023-04-19

### Improvements
- Applying the patch of wasmvm, upgrade from v1.2.1 to v1.2.3
- Add makefile support ledger

## [v0.5.0] - 2023-04-18

### Features
- Implement ibc-middleware support GMP protocol of Axelar network
- Upgrade wasmd from 0.29.1 to 0.31.0

## [v0.4.4] - 2023-03-13

### Improvements
- Prevent excluded addresses from receiving fund
- Fix minor bugs

## [v0.4.3] - 2023-02-02

### Improvements
- Upgrade Comos SDK from v0.45.9 -> v0.45.11
- Change tendermint repo to InformalSystem
- Refactor some code unused

## [v0.4.2] - 2022-12-05

### Improvements
- Allow IBC clients to be recoverd after expired

## [v0.4.1] - 2022-11-11

### Improvements
- Register `uaura` as default denom

## [v0.4.0] - 2022-11-02

### Features
- Add tx type and CLI command to create periodic vesting account 


## [v0.3.3] - 2022-10-15

### Improvements
- Bump cosmos-sdk to version `v0.45.9`
- Bump wasmd to version `v0.29.1`
- Bump ibc-go to version `v3.3.0`

### Bug Fixes
- Apply dragonberry patch

## [v0.3.2] - 2022-09-29

### Bug Fixes

- Fix state-sync mode does not sync wasm data.
- Fix spend limit feature when using feegrant to allow user execute specific contract. 

## [v0.3.1] - 2022-09-07

### Improvements

- Bump cosmos-sdk to version `v0.45.6`
- Exclude whitelist balances when querying current supply by specific denom

## [v0.3.0] - 2022-08-09

### Features

- Bump wasmd version to version `v0.18.0`
- Bump cosmos-sdk to version `v0.45.5`
- Feegrant allows to use contract address 