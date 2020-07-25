package main

import (
	"fmt"

	gsrpc "github.com/Snowfork/go-substrate-rpc-client"
	"github.com/Snowfork/go-substrate-rpc-client/signature"
	gsrpcTypes "github.com/Snowfork/go-substrate-rpc-client/types"
)

const (
	API_URL               = "ws://127.0.0.1:9944"
	SENDER_ADDRESS        = "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty"
	SENDER_HEX_ACCOUNT_ID = "0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48"
)

type AppID gsrpcTypes.Bytes256

type RawMessage struct {
	data []uint8
}

func main() {

	api, err := gsrpc.NewSubstrateAPI(API_URL)
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to testing account 5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty
	senderAccount, err := gsrpcTypes.NewAddressFromHexAccountID(SENDER_HEX_ACCOUNT_ID)
	if err != nil {
		panic(err)
	}

	amount := gsrpcTypes.NewUCompactFromUInt(10)
	c, err := gsrpcTypes.NewCall(meta, "Balances.transfer", senderAccount, amount)
	if err != nil {
		panic(err)
	}

	// Create the extrinsic
	ext := gsrpcTypes.NewExtrinsic(c)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(err)
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(err)
	}

	key, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", signature.TestKeyringPairAlice.PublicKey, nil)
	if err != nil {
		panic(err)
	}

	var accountInfo gsrpcTypes.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(err)
	}

	nonce := uint32(accountInfo.Nonce)

	o := gsrpcTypes.SignatureOptions{
		BlockHash:   genesisHash,
		Era:         gsrpcTypes.ExtrinsicEra{IsMortalEra: false},
		GenesisHash: genesisHash,
		Nonce:       gsrpcTypes.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion: rv.SpecVersion,
		Tip:         gsrpcTypes.NewUCompactFromUInt(0),
	}

	fmt.Printf("Sending %v from %#x to %#x with nonce %v", amount, signature.TestKeyringPairAlice.PublicKey, senderAccount.AsAccountID, nonce)

	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	sendTx(api, ext)
	// sendAndWatchTx(api, ext)
}

// Send the extrinsic
func sendTx(api *gsrpc.SubstrateAPI, ext gsrpcTypes.Extrinsic) {
	hash, err := api.RPC.Author.SubmitExtrinsic(ext)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Transfer sent with hash %#x\n", hash)
}

// Do the transfer and track the actual status
func sendAndWatchTx(api *gsrpc.SubstrateAPI, ext gsrpcTypes.Extrinsic) {
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		status := <-sub.Chan()
		fmt.Printf("Transaction status: %#v\n", status)

		if status.IsInBlock {
			fmt.Printf("Completed at block hash: %#x\n", status.AsInBlock)
			return
		}
	}
}

func initWatcher() {
	api, err := gsrpc.NewSubstrateAPI(API_URL)
	if err != nil {
		panic(err)
	}

	sub, err := api.RPC.Chain.SubscribeNewHeads()
	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	count := 0

	for {
		head := <-sub.Chan()
		fmt.Printf("Chain is at block: #%v\n", head.Number)
		count++

		if count == 10 {
			sub.Unsubscribe()
			break
		}
	}
}
