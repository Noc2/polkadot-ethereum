package main

import (
	"github.com/Snowfork/polkadot-ethereum/scratch/crypto/sr25519"
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

	// --------------------------------------- OUR ACCOUNT INFORMATION ---------------------------------------
	// Secret phrase `coral piece session toilet february finger furnace shine metal advance summer health` is account:
	//   Network ID/version: substrate
	//   Secret seed:        0x2d268c3c96800f649d9b91f20e4fdf3f25b08fb0051f126942bd195d94ba844e
	//   Public key (hex):   0x46e7708880bbe783e6e22ee29fd6c07d76ea44ff36b6401a88bccadd30135c74
	//   Account ID:         0x46e7708880bbe783e6e22ee29fd6c07d76ea44ff36b6401a88bccadd30135c74
	//   SS58 Address:       5Dffxh5ZwZzHiFNsWE1XDkvVAnEfcW1hc8qGcBJSZcfWi8JN
	// -------------------------------------------------------------------------------------------------------

	seedPhrase := "coral piece session toilet february finger furnace shine metal advance summer health"
	// seed := "0x2d268c3c96800f649d9b91f20e4fdf3f25b08fb0051f126942bd195d94ba844e"

	network := "Development"
	keyringPair, err := sr25519.GenerateKeypair(network)
	if err != nil {
		panic(err)
	}

	// keyringPair, err := signature.KeyringPairFromSecret()
	// if err != nil {
	// 	panic(err)
	// }

	key, err := gsrpcTypes.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey, nil)
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
