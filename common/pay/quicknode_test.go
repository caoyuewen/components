package pay

import (
	"encoding/json"
	"fmt"
	"testing"
)

var quick = QuickNode{
	ApiKey:      "QN_d0b77e2d5daf4ada8030bf658b602f72",
	NotifyEmail: "liufengzhx@gmail.com",
	Callback:    "https://effic.in/payment/callback/quicknode",
	JumpUrl:     "https://effic.in",
	Domain:      "https://snowy-divine-bridge.tron-mainnet.quiknode.pro/38651257af9f5ef32ca03dce1a09994b1003d1fb/jsonrpc",
}

type PaymentBaseConf struct {
	ApiKey      string
	NotifyEmail string
	Callback    string
	Domain      string
}

func TestName(t *testing.T) {

	config, _ := json.Marshal(quick)
	fmt.Println(string(config))
}

func TestQuickNode_CheckWebhooksConfig(t *testing.T) {

	wallets := []string{
		"TMykpgR4V51s3qDMvk8ztdEFZJjCnm5kcZ",
		"TPT2MUXWRvpy17MPcVdLTfAg38rocspzxB",
		"TXYGHnt9ojxLZY4Gtx3UuKd2c7Fi47ojHY",
	}

	err := quick.CheckWebhooksConfig(wallets)
	if err != nil {
		panic(err)
	}

}

func TestQuickNode_GetTxDetailByHash(t *testing.T) {

	data, err := quick.GetTxDetailByHash("0xdb98eb6ca5d373afd1fd8642951ed98370febfbcb45633d5ecbcce4b9b9ef5f7")
	if err != nil {
		panic(err)
	}

	marshal, _ := json.Marshal(data)
	fmt.Println(string(marshal))

}

func TestQuickNode_WebhooksList(t *testing.T) {

	list, err := quick.WebhooksList()
	if err != nil {
		panic(err)
	}

	for _, v := range list.Data {
		fmt.Println("name:", v.Name)
		fmt.Println("id:", v.Id)
	}

}

func TestQuickNode_WebhooksInsert(t *testing.T) {

	wallets := []string{
		"TMykpgR4V51s3qDMvk8ztdEFZJjCnm5kcZ",
		"TPT2MUXWRvpy17MPcVdLTfAg38rocspzxB",
	}

	name := "tron node test"

	bytes, err := quick.CreateWebhook(name, wallets)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

}

func TestQuickNode_Webhooks(t *testing.T) {

	wallets := []string{
		"TMykpgR4V51s3qDMvk8ztdEFZJjCnm5kcZ",
		"TPT2MUXWRvpy17MPcVdLTfAg38rocspzxB",
		"TXYGHnt9ojxLZY4Gtx3UuKd2c7Fi47ojHY",
	}

	name := "tron node test"

	bytes, err := quick.CreateWebhook(name, wallets)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

}

func TestQuickNode_WebhooksUpdate(t *testing.T) {

	bytes, err := quick.WebhookUpdate("9af9baea-d4b3-454e-bbdb-2f04897d31b4", "active")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

}

func TestQuickNode_WebhooksDelete(t *testing.T) {

	list, err := quick.WebhooksList()
	if err != nil {
		panic(err)
	}

	for _, v := range list.Data {
		fmt.Println("name:", v.Name)
		fmt.Println("id:", v.Id)
		err := quick.WebhooksDelete(v.Id)
		if err != nil {
			panic(err)
		}
	}
}
