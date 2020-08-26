package core

import (
	"fmt"
	ontSdk "github.com/ontio/ontology-go-sdk"
	common3 "github.com/ontio/ontology/common"
	"github.com/ontio/oscore/oscoreconfig"
	"testing"
)

func TestWhetherTest(t *testing.T) {
	wallet, err := ontSdk.OpenWallet("wallet.dat")
	if err != nil {
		panic(err)
	}
	oscoreAccount, err := wallet.GetDefaultAccount([]byte("123456"))
	if err != nil {
		panic(err)
	}

	oscoreconfig.DefOscoreConfig.OscoreAccount = oscoreAccount
	oscoreconfig.DefOscoreConfig.DDXFAPIServer = "http://106.75.209.209:8080"

	addr, err := common3.AddressFromBase58("ALBoDa6bkA3WPPSTSaGogRuxNYgHPLveee")
	if err != nil {
		panic(err)
	}

	txHash, _, err := BuildWetherForcastTransaction(oscoreAccount, "30", "31", addr.ToHexString())
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %s\n", txHash)
}
