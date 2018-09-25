package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

//go test -v -run TestGenerateAddress
func TestGenerateAddress(t *testing.T) {
	// public extend key
	const key = "xpub661MyMwAqRbcG9asPR1fmXH4hyg3GpKPhZxMFivAe1E47sEqRhqXnZeh2xEDVykdj3ECfEUAaQ2RjJRyhAZxLAhXaWzUWLkAk1g6Crwi4ue"
	const start = "/0/1"
	addrs, err := GenerateAddress(key, 2000, start, &chaincfg.MainNetParams, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2000 {
		t.Fatal("error count")
	}
	for _, v := range addrs {
		fmt.Println(v)
	}
}
