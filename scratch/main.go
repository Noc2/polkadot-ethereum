package main

import (
	"fmt"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	gsrpcTypes "github.com/centrifuge/go-substrate-rpc-client/types"
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

	// publicKeyOne := []byte{47, 140, 97, 41, 216, 22, 207, 81, 195, 116, 188, 127, 8, 195, 230, 62, 209, 86, 207, 120, 174, 251, 74, 101, 80, 217, 123, 135, 153, 121, 119, 238}
	charliePubKey := []byte{0x90, 0xb5, 0xab, 0x20, 0x5c, 0x69, 0x74, 0xc9, 0xea, 0x84, 0x1b, 0xe6, 0x88, 0x86, 0x46, 0x33, 0xdc, 0x9c, 0xa8, 0xa3, 0x57, 0x84, 0x3e, 0xea, 0xcf, 0x23, 0x14, 0x64, 0x99, 0x65, 0xfe, 0x22}

	var CharlieSr25519 = signature.KeyringPair{
		URI:       "//Charlie",
		Address:   "5FLSigC9HGRKVhB9FiEo4Y3koPsNmBmLJbpXg2mp1hXcS59Y",
		PublicKey: publicKeyOne,
	}

	key, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", publicKeyOne, nil)
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

	fmt.Printf("Sending %v from %#x to %#x with nonce %v", amount, CharlieSr25519.PublicKey, senderAccount.AsAccountID, nonce)

	// Sign the transaction using Charlie's default account
	err = ext.Sign(CharlieSr25519, o)
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
