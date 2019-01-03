package main

import (
	"bytes"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
)

func IsValidAddress(address string) bool {
	addressByte := base58.Decode(address)
	addressLen := len(addressByte)
	if addressLen < 4 {
		return false
	}
	payload := addressByte[:addressLen-4]
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return bytes.Equal(hash2[:4], addressByte[addressLen-4:])
}

func GetPubKeyHashByAddress(address string) []byte {
	//check todo
	addressByte := base58.Decode(address)
	addressLen := len(addressByte)
	return addressByte[1 : addressLen-4]
}
