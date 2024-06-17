package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
)

type config struct {
	threads int
}

var (
	counter uint64 = 0
	wg      sync.WaitGroup
)

func parseConfig() *config {
	var cfg config
	flag.IntVar(&cfg.threads, "threads", runtime.NumCPU(), "Threads count")
	flag.Parse()
	return &cfg
}

func genWallet(mnemonic chan []string, api wallet.TonAPI, ctx context.Context, block *ton.BlockIDExt, thread int) {
	for {
		_mnemonic := <-mnemonic
		_wallet := getWallet(api, _mnemonic)
		_address := _wallet.WalletAddress()

		balance, err := _wallet.GetBalance(ctx, block)
		if err != nil {
			log.Fatal("getWallet err:", err.Error())
		}

		if balance.Nano().Uint64() > 0 {
			log.Printf("Address: %s \nMnemo: %s \nBalance: %s\n", _address.String(), strings.Join(_mnemonic, " "), balance.Nano().String())
		}

		log.Printf("-- Thread: #%d \nAddress: %s \nMnemo: %s \nBalance: %s\n", thread, _address.String(), strings.Join(_mnemonic, " "), balance.Nano().String())

		atomic.AddUint64(&counter, 1)
		_counter := atomic.LoadUint64(&counter)
		if _counter%1000 == 0 {
			log.Printf("Checked %d addresses...", _counter)
		}
	}
}

func main() {
	appCfg := parseConfig()
	chData := make(chan []string)

	client := liteclient.NewConnectionPool()

	err := client.AddConnectionsFromConfigUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		log.Fatal("add connection err:", err.Error())
	}
	api := ton.NewAPIClient(client, ton.ProofCheckPolicyFast).WithRetry()

	ctx := client.StickyContext(context.Background())

	block, err := api.CurrentMasterchainInfo(context.Background())
	if err != nil {
		log.Fatalln("get masterchain info err: ", err.Error())
		return
	}

	for t := 0; t < appCfg.threads; t++ {
		fmt.Printf("Thread %d start...\n", t)
		go genWallet(chData, api, ctx, block, t)
	}

	for {
		chData <- createMnemonic()
	}
}
