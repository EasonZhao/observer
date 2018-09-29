package main

import (
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	util "github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"strconv"
	"strings"
)

//GenerateAddress generate bitcoin address
func GenerateAddress(key string, count int, start string, net *chaincfg.Params, isSigWit bool) ([]string, error) {
	master, err := hdkeychain.NewKeyFromString(key)
	if err != nil {
		return nil, err
	}
	strs := strings.Split(start, "/")
	for _, v := range strs {
		if v == "" {
			continue
		}
		i, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		child, err := master.Child(uint32(i))
		if err != nil {
			return nil, err
		}
		master = child
	}
	result := make([]string, count)
	for i := 0; i < count; i++ {
		child, err := master.Child(uint32(i))
		if err != nil {
			return nil, err
		}
		if isSigWit {
			pubKey, err := child.ECPubKey()
			address, err := GetWitnessAddress(pubKey, net)
			if err != nil {
				return nil, err
			}
			result[i] = address
		} else {
			pubHash, err := child.Address(net)
			if err != nil {
				return nil, err
			}
			result[i] = pubHash.EncodeAddress()
		}
	}
	return result, nil
}

// GetWitnessAddress get witness address
func GetWitnessAddress(pubKey *btcec.PublicKey, net *chaincfg.Params) (string, error) {
	pubKeyHash := util.Hash160(pubKey.SerializeCompressed())
	witAddr, err := util.NewAddressWitnessPubKeyHash(pubKeyHash, net)
	witnessProgram, err := txscript.PayToAddrScript(witAddr)
	if err != nil {
		return "", err
	}
	address, err := util.NewAddressScriptHash(witnessProgram, net)
	if err != nil {
		return "", err
	}
	return address.EncodeAddress(), nil
}
