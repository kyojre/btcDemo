package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

func test() {
	curve := elliptic.P256()
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := privateKey.PublicKey

	data := "hello world"
	hash := sha256.Sum256([]byte(data))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Panic(err)
	}
	signature := append(r.Bytes(), s.Bytes()...)

	fmt.Printf("%x\n", r.Bytes)
	fmt.Printf("%x\n", s.Bytes)
	fmt.Printf("%x\n", signature)

	res := ecdsa.Verify(&pubKey, hash[:], r, s)
	fmt.Printf("%v\n", res)

	sigLen := len(signature)
	r1 := big.Int{}
	s1 := big.Int{}
	r1.SetBytes(signature[:sigLen/2])
	s1.SetBytes(signature[sigLen/2:])

	res1 := ecdsa.Verify(&pubKey, hash[:], &r1, &s1)
	fmt.Printf("%v\n", res1)
}
