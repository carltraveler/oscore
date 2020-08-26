package ddxf

import (
	"fmt"
	ontSdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"testing"
)

func TestPublishAPI(t *testing.T) {
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
	apiKey := &tables.APIKey{
		OrderId: "xxx",
		OntId:   "did:ont:ATaZYj1UoHaeJPLTMYHNuuKRhvWeeLdNJ9",
	}

	txHash, err := BuyerRecordDDXF(apiKey)
	if err != nil {
		panic(err)
	}

	fmt.Printf("txHash %s\n", txHash)
}
