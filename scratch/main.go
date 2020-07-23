package main

import (
	gsrpc "github.com/centrifuge/go-substrate-rpc-client"
	gsrpcTypes "github.com/centrifuge/go-substrate-rpc-client/types"
)

type AppID gsrpcTypes.Bytes256

type RawMessage struct {
	data []uint8
}

func main() {
	api, err := gsrpc.NewSubstrateAPI("wss://127.0.0.1:9944")
	if err != nil {
		panic(err)
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}

	rawMessage := RawMessage{data: []uint8{0, 1, 3, 4}}
	message, err := gsrpcTypes.EncodeToHexString(rawMessage)

}
