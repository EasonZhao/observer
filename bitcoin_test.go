package main

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"testing"
)

//go test -v -run TestGenerateAddress
func TestGenerateAddress(t *testing.T) {
	// public extend key
	const key = "xpub661MyMwAqRbcFXeLasdACMrCN7iP2aJPAyNtFDy5Gc3QxUGu7iGvpRY6LPyRzy6T3meQHP54UzQQqs6gsMUryedaQxz7rMPbsoqUFBCkMFD"
	const start = "/0/77"
	count := 10
	addrs, err := GenerateAddress(key, count, start, &chaincfg.TestNet3Params, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != count {
		t.Fatal("error count")
	}
	if addrs[0] != "2Mx8enZb39sqkVX6qoAeP9NzBYYMtmYgMGq" {
		t.Fatalf("address is %s not 2Mx8enZb39sqkVX6qoAeP9NzBYYMtmYgMGq", addrs[0])
	} else if addrs[count-1] != "2Mw6g4Zk3bMJE5fkhttMwYNwwiqiDkn8HVr" {
		t.Fatalf("address is %s not 2Mw6g4Zk3bMJE5fkhttMwYNwwiqiDkn8HVr", addrs[0])
	}
}

//go test -v -run TestGetWitnessAddress
func TestGetWitnessAddress(t *testing.T) {
	const str = "cUeH6WQemR1BAxgbjXnzr4tNQuA4VpUQ4MHxBKKSSS4JY5Hzc3vi"
	wif, err := btcutil.DecodeWIF(str)
	if err != nil {
		t.Fatal(err)
	}
	pubKey := wif.PrivKey.PubKey()
	address, err := GetWitnessAddress(pubKey, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	if address != "2MyVZ51ZTKQve4PA29Bz5gs4yA2QcSPxmHb" {
		t.Fatalf("address is %s not 2MyVZ51ZTKQve4PA29Bz5gs4yA2QcSPxmHb", address)
	}
}
