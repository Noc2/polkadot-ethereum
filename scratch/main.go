package main

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/signature"
	gsrpcTypes "github.com/centrifuge/go-substrate-rpc-client/types"
)

type AppID gsrpcTypes.Bytes256

type RawMessage struct {
	data []uint8
}

func main() {

	api, err := gsrpc.NewSubstrateAPI("ws://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	// Create a call, transferring 12345 units to testing account 5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty
	testAccount, err := gsrpcTypes.NewAddressFromHexAccountID("0x8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48")
	if err != nil {
		panic(err)
	}

	// TODO: The type from the example doesn't exist: gsrpcTypes.NewUCompactFromUInt(10))
	// amount := gsrpcTypes.NewI256(*big.NewInt(10)) // TODO: using this type instead
	c, err := gsrpcTypes.NewCall(meta, "Balances.transfer", testAccount, gsrpcTypes.NewUCompactFromUInt(10))
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

	// var TestKeyringPairAlice = KeyringPair{
	// 	URI:       "//Alice",
	// 	PublicKey: []byte{0xd4, 0x35, 0x93, 0xc7, 0x15, 0xfd, 0xd3, 0x1c, 0x61, 0x14, 0x1a, 0xbd, 0x4, 0xa9, 0x9f, 0xd6, 0x82, 0x2c, 0x85, 0x58, 0x85, 0x4c, 0xcd, 0xe3, 0x9a, 0x56, 0x84, 0xe7, 0xa5, 0x6d, 0xa2, 0x7d}, //nolint:lll
	// 	Address:   "5GrwvaEF5zXb26Fz9rcQpDWS57CtERHpNehXCPcNoHGKutQY",
	// }

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

	_ = o

	// Sign the transaction using Alice's default account
	err = ext.Sign(signature.TestKeyringPairAlice, o)
	if err != nil {
		panic(err)
	}

	// // Send the extrinsic
	// hash, err := api.RPC.Author.SubmitExtrinsic(ext)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("Transfer sent with hash %#x\n", hash)
}
