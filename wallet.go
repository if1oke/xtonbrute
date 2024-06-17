package main

import (
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
)

func createMnemonic() []string {
	return wallet.NewSeed()
}

func getWallet(api wallet.TonAPI, mnemonic []string) *wallet.Wallet {
	_wallet, err := wallet.FromSeed(api, mnemonic, wallet.V4R2)
	if err != nil {
		log.Fatal("getWallet err:", err.Error())
	}
	return _wallet
}
