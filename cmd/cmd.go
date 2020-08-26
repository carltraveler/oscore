package cmd

import (
	"encoding/json"
	"fmt"
	ontSdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology-go-sdk/utils"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/ontology/common/password"
	common2 "github.com/ontio/oscore/common"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
)

func SetOntologyConfig(ctx *cli.Context) error {
	cf := ctx.String(GetFlagName(ConfigfileFlag))
	if _, err := os.Stat(cf); os.IsNotExist(err) {
		// if there's no config file, use default config
		updateConfigByCmd(ctx)
		return nil
	}

	file, err := os.Open(cf)
	if err != nil {
		return err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	cfg := &oscoreconfig.Config{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return err
	}
	*oscoreconfig.DefOscoreConfig = *cfg

	if oscoreconfig.DefOscoreConfig.WalletName == "" || oscoreconfig.DefOscoreConfig.Layer2MainNetNode == "" || oscoreconfig.DefOscoreConfig.Layer2TestNetNode == "" {
		return fmt.Errorf("oscore walletName/layer2MainNetAddress/layer2TestNetAddress  is nil")
	}

	wallet, err := ontSdk.OpenWallet(oscoreconfig.DefOscoreConfig.WalletName)
	if err != nil {
		return err
	}
	passwd, err := password.GetAccountPassword()
	if err != nil {
		return err
	}
	oscoreAccount, err := wallet.GetDefaultAccount(passwd)
	if err != nil {
		return err
	}

	sdk := ontSdk.NewOntologySdk()
	layer2Sdk := ontSdk.NewOntologySdk()
	switch oscoreconfig.DefOscoreConfig.NetWorkId {
	case oscoreconfig.NETWORK_ID_MAIN_NET:
		log.Infof("currently Main net")
		sdk.NewRpcClient().SetAddress(oscoreconfig.ONT_MAIN_NET)
		layer2Sdk.NewRpcClient().SetAddress(oscoreconfig.DefOscoreConfig.Layer2MainNetNode)
		oscoreconfig.DefOscoreConfig.NetType = oscoreconfig.MainNet
	case oscoreconfig.NETWORK_ID_POLARIS_NET:
		log.Infof("currently test net")
		sdk.NewRpcClient().SetAddress(oscoreconfig.ONT_TEST_NET)
		layer2Sdk.NewRpcClient().SetAddress(oscoreconfig.DefOscoreConfig.Layer2TestNetNode)
		oscoreconfig.DefOscoreConfig.NetType = oscoreconfig.TestNet
	case oscoreconfig.NETWORK_ID_SOLO_NET:
		log.Infof("currently solo net")
		sdk.NewRpcClient().SetAddress(oscoreconfig.ONT_SOLO_NET)
		// solo simulation with test net. but different contract and owner
		layer2Sdk.NewRpcClient().SetAddress(oscoreconfig.DefOscoreConfig.Layer2TestNetNode)
	default:
		return fmt.Errorf("error network id %d", oscoreconfig.DefOscoreConfig.NetWorkId)
	}

	oscoreconfig.DefOscoreConfig.OntSdk = sdk
	oscoreconfig.DefOscoreConfig.Layer2Sdk = layer2Sdk
	oscoreconfig.DefOscoreConfig.OscoreAccount = oscoreAccount
	//CheckLayer2InitAddress()
	return nil
}

func CheckLayer2InitAddress() error {
	layer2Sdk := oscoreconfig.DefOscoreConfig.Layer2Sdk
	if layer2Sdk == nil {
		return fmt.Errorf("layer2 sdk should not be nil")
	}

	if len(oscoreconfig.DefOscoreConfig.Layer2Contract) == 0 {
		return fmt.Errorf("layer2 contract address or contract not init")
	}

	log.Infof("layer2Contract %s", oscoreconfig.DefOscoreConfig.Layer2Contract)

	contractAddr, err := common.AddressFromHexString(oscoreconfig.DefOscoreConfig.Layer2Contract)
	// incase that is a file name.
	if err != nil || len(oscoreconfig.DefOscoreConfig.Layer2Contract) != common.ADDR_LEN*2 {
		code, err := ioutil.ReadFile(oscoreconfig.DefOscoreConfig.Layer2Contract)
		if err != nil {
			return fmt.Errorf("error in ReadFile: %s, %s\n", oscoreconfig.DefOscoreConfig.Layer2Contract, err)
		}

		codeHash := common.ToHexString(code)
		contractAddr, err = utils.GetContractAddress(codeHash)
		if err != nil {
			return fmt.Errorf("error get contract address %s", err)
		}

		payload, err := layer2Sdk.GetSmartContract(contractAddr.ToHexString())
		if payload == nil || err != nil {
			txHash, err := layer2Sdk.NeoVM.DeployNeoVMSmartContract(0, 200000000, oscoreconfig.DefOscoreConfig.OscoreAccount, true, codeHash, "oscore layer2 contract", "1.0", "oscore", "email", "desc")
			if err != nil {
				return fmt.Errorf("deploy contract %s err: %s", oscoreconfig.DefOscoreConfig.Layer2Contract, err)
			}

			_, err = common2.GetLayer2EventByTxHash(txHash.ToHexString())
			if err != nil {
				return fmt.Errorf("deploy contract failed %s", err)
			}
			log.Infof("deploy concontract success")
		}

		log.Infof("the contractAddr hexstring is %s", contractAddr.ToHexString())
		oscoreconfig.DefOscoreConfig.Layer2Contract = contractAddr.ToHexString()
	}

	contractAddr, err = common.AddressFromHexString(oscoreconfig.DefOscoreConfig.Layer2Contract)
	if err != nil {
		return err
	}

	for {
		res, err := layer2Sdk.NeoVM.PreExecInvokeNeoVMContract(contractAddr, []interface{}{"init_status", 2})
		if err != nil {
			return fmt.Errorf("err get init_status %s", err)
		}

		if res.State == 0 {
			return fmt.Errorf("init statuc exec failed state is 0")
		}

		addrB, err := res.Result.ToByteArray()
		if err != nil {
			return fmt.Errorf("error init_status toByteArray %s", err)
		}

		oscoreAddrBase58 := oscoreconfig.DefOscoreConfig.OscoreAccount.Address.ToBase58()
		if len(addrB) != 0 {
			addrO, err := common.AddressParseFromBytes(addrB)
			if err != nil {
				return fmt.Errorf("AddressParseFromBytes err: %s", err)
			}

			log.Infof("layer2 address already init owner to addr %s", addrO.ToBase58())
			if addrO.ToBase58() != oscoreAddrBase58 {
				return fmt.Errorf("contract addr not equal. owner is %s. but oscoreAccount init to %s", addrO.ToBase58(), oscoreAddrBase58)
			}
			break
		} else {
			//log.Infof("start init layer2 addr owner to address %s", oscoreAddrBase58)
			txHash, err := layer2Sdk.NeoVM.InvokeNeoVMContract(0, 200000, nil, oscoreconfig.DefOscoreConfig.OscoreAccount, contractAddr, []interface{}{"init", oscoreconfig.DefOscoreConfig.OscoreAccount.Address})
			if err != nil {
				return fmt.Errorf("init layer2 owner err0 %s", err)
			}
			_, err = common2.GetLayer2EventByTxHash(txHash.ToHexString())
			if err != nil {
				return fmt.Errorf("init layer2 owner err1: %s", err)
			}
			log.Infof("init layer2 addr owner to address %s success.", oscoreAddrBase58)
		}
	}

	txHash, err := layer2Sdk.NeoVM.InvokeNeoVMContract(0, 200000, nil, oscoreconfig.DefOscoreConfig.OscoreAccount, contractAddr, []interface{}{"StoreUsedNum", []interface{}{"userid_init_test", "orderId_init_test", 9873456}})
	if err != nil {
		return fmt.Errorf("StoreUsedNum test failed %s", err)
	}

	_, err = common2.GetLayer2EventByTxHash(txHash.ToHexString())
	if err != nil {
		return fmt.Errorf("init layer2 owner err1: %s", err)
	}
	log.Infof("test StoreUsedNum success ")
	return nil
}

func updateConfigByCmd(ctx *cli.Context) error {
	port := ctx.Uint(GetFlagName(RestPortFlag))
	if port != 0 {
		oscoreconfig.DefOscoreConfig.RestPort = port
	}
	networkId := ctx.Uint(GetFlagName(NetworkIdFlag))
	if networkId > 3 {
		return fmt.Errorf("networkid should be between 1 and 3, curr: %d", networkId)
	}
	oscoreconfig.DefOscoreConfig.NetWorkId = networkId
	//rpc := oscoreconfig.ONT_MAIN_NET
	//if networkId == oscoreconfig.NETWORK_ID_POLARIS_NET {
	//	rpc = oscoreconfig.ONT_TEST_NET
	//} else if networkId == oscoreconfig.NETWORK_ID_SOLO_NET {
	//	rpc = oscoreconfig.ONT_SOLO_NET
	//}
	sdk := ontSdk.NewOntologySdk()
	sdk.NewRpcClient().SetAddress(oscoreconfig.ONT_TEST_NET)
	oscoreconfig.DefOscoreConfig.OntSdk = sdk
	return nil
}

func PrintErrorMsg(format string, a ...interface{}) {
	format = fmt.Sprintf("\033[31m[ERROR] %s\033[0m\n", format) //Print error msg with red color
	fmt.Printf(format, a...)
}
