package main

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
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
