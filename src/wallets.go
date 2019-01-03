package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

const walletFileName = "../db/wallets.dat"

type Wallets struct {
	WalletsMap map[string]*Wallet
}

func NewWallets() *Wallets {
	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.LoadFile()
	return &wallets
}

func (this *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()

	this.WalletsMap[address] = wallet
	this.SaveToFile()
	return address
}

func (this *Wallets) SaveToFile() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	//encode 接口时候要注册
	gob.Register(elliptic.P256())
	err := encoder.Encode(*this)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(walletFileName, buffer.Bytes(), 0644)
	if err != nil {
		log.Panic(err)
	}
}

func (this *Wallets) LoadFile() {
	_, err := os.Stat(walletFileName)
	if os.IsNotExist(err) {
		return
	}
	buffer, err := ioutil.ReadFile(walletFileName)
	if err != nil {
		log.Panic(err)
	}
	decoder := gob.NewDecoder(bytes.NewReader(buffer))
	//decode 接口时候要注册
	gob.Register(elliptic.P256())
	err = decoder.Decode(this)
	if err != nil {
		log.Panic(err)
	}
}
