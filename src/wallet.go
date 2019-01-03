package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

type Wallet struct {
	PriKey *ecdsa.PrivateKey
	PubKey []byte
}

func NewWallet() *Wallet {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publicKey := privateKey.PublicKey
	pubKey := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	wallet := Wallet{
		PriKey: privateKey,
		PubKey: pubKey,
	}
	return &wallet
}

func (this *Wallet) NewAddress() string {
	pubKey := this.PubKey
	hash := sha256.Sum256(pubKey)
	rip160hasher := ripemd160.New()
	_, err := rip160hasher.Write(hash[:])
	if err != nil {
		log.Panic(err)
	}
	rip160HashValue := rip160hasher.Sum(nil)
	version := byte(0)
	payload := append([]byte{version}, rip160HashValue...)

	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])

	checkCode := hash2[:4]
	payload = append(payload, checkCode...)

	address := base58.Encode(payload)
	return address
}
