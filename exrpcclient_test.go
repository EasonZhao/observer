package main

import (
	"testing"
)

//go test -v -run TestExGetBlock
func TestExGetBlock(t *testing.T) {
	block, err := exGetBlock("http://obs:obs@localhost:18332", "0000000000008dea9f388858e176f8e180708ffa03b0c0d2fe70fc770d7779c5")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block.Confirmations)
}
