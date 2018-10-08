package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
)

func main() {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		panic(err)
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		panic(err)
	}
	if masterKey.IsPrivate() {
		fmt.Println(masterKey.String())
		masterKey, err = masterKey.Neuter()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(masterKey.String())
}

func mytest() {
	const extendKey = "tprv8ZgxMBicQKsPeyGNZnjScfnJ6gCExqXjoQRrojefJF6L35LdG743Bh6haXrg89cgBaBE28UStzweiPaG5QTqD6qPsra2wuCY88v1eQWTXGg"
	master, err := hdkeychain.NewKeyFromString(extendKey)
	if err != nil {
		panic(err)
	}
	child, err := master.Child(0x80000000 + 0)
	if err != nil {
		panic(err)
	}
	child, err = child.Child(0x80000000 + 0)
	if err != nil {
		panic(err)
	}
	child, err = child.Child(0x80000000 + 0)
	if err != nil {
		panic(err)
	}
	privKey, err := child.ECPrivKey()
	if err != nil {
		panic(err)
	}
	wif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(wif.String())
}
