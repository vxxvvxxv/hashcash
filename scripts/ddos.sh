#!/usr/bin/env bash
export CLIENT_DDOS_MODE=true
go run ./cmd/client/main.go --race
