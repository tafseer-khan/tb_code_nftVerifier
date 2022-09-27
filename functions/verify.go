package lib

import (
	"fmt"
	"math/big"
	"bitbucket.org/taubyte/go-sdk/event"
	ethereum "bitbucket.org/taubyte/go-sdk/ethereum/client"
	ethBytes "bitbucket.org/taubyte/go-sdk/ethereum/client/bytes"
	http "bitbucket.org/taubyte/go-sdk/http/client"
)

//export verify
func verify(e event.Event) uint32 {
	h, err := e.HTTP()
	if err != nil {
		panic(err)
		return 1
	}

	errReturn := func(msg string) {
		h.Write([]byte(msg))
		h.Return(404)
	}

	client, err := ethereum.New("https://goerli.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161")
	if err != nil {
		errReturn("Unable to connect to rpc client")
		return 1
	}

	contractAddress := "0xf4910C763eD4e47A585E2D34baA9A4b611aE448C"

	httpClient, err := http.New()
	if err != nil {
		errReturn("Unable to get http client")
		return 1
	}

	req, err := httpClient.Request(fmt.Sprintf("https://api-goerli.etherscan.io/api?module=contract&action=getabi&format=raw&address=%s&apikey=GJ9AZ69URIFNGXCUGXN2GSSRSPB16HI3M7", contractAddress))
	if err != nil {
		errReturn("Unable to get http request")
		return 1
	}

	res, err := req.Do()
	if err != nil {
		errReturn("Unable to get http client")
		return 1
	}

	contract, err := client.NewBoundContract(res.Body(), contractAddress)
	if err != nil {
		errReturn("Unable to create bound contract")
		return 1
	}

	balanceOf, err := contract.Method("balanceOf")
	if err != nil {
		errReturn("Unable to create method balanceOf")
		return 1
	}

	addressString := h.Query().Get("address")
	address := ethBytes.AddressFromHex(addressString)
	tokenId, ok := new(big.Int).SetString("80867650201096745079196794753906950580251458356280840071563152651088098754660", 10)
	if ok == false {
		errReturn("Unable to create tokenId")
		return 1
	}

	outputs, err := balanceOf.Call(address, tokenId)
	if err != nil {
		errReturn("Cannot call balance of with: " + err.Error())
		return 1
	}

	if outputs[0].(*big.Int).Cmp(big.NewInt(0)) == 0 {
		errReturn("Holder does not own NFT")
		return 1
	}

	h.Return(200)

	return 0
}
