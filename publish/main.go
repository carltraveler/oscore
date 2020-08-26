package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	ontSdk "github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/cmd/utils"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/cmd"
	"github.com/ontio/oscore/common"
	"github.com/ontio/oscore/core"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/urfave/cli"
)

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		cmd.PrintErrorMsg(err.Error())
		os.Exit(1)
	}
}

var (
	OperationTypeFlag = cli.UintFlag{
		Name:  "t",
		Usage: "Set the operation to `<type>` (0~1). 0:publish 1:takeOrder",
		Value: uint(0),
	}
)

const (
	OPTYPE_PUBLISH   uint = 0
	OPTYPE_TAKEORDER uint = 1
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "publish CLI"
	app.Action = PublishAPIJSON
	app.Version = oscoreconfig.Version
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		cmd.LogLevelFlag,
		cmd.ConfigfileFlag,
		OperationTypeFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

type TakeOderApiIDSpec struct {
	ApiId  uint32 `json:"apiId"`
	SpecId uint32 `json:"specId"`
}

type PublishConfig struct {
	DBConfig   *oscoreconfig.DBConfig `json:"db_config"`
	AccessMode int32                  `json:"accessMode"`
	OntId      string                 `json:"ontId"`
	Author     string                 `json:"author"`
	Specs      []TakeOderApiIDSpec    `json:"specs"`
	RestPort   uint                   `json:"restPort"`
	OscoreHost string                 `json:"oscoreHost"`
	JsonFiles  []string               `json:"jsonFiles"`
}

var DefPublishConfig = PublishConfig{
	DBConfig: &oscoreconfig.DBConfig{
		ProjectDBUrl:      "127.0.0.1:3306",
		ProjectDBUser:     "steven",
		ProjectDBPassword: "abcd1234",
		ProjectDBName:     "oscoreunittest",
	},
	OntId:      "did:ont:Ad4pjz2bqep4RhQrUAzMuZJkBC3qJ1tZuT",
	Author:     "admin",
	AccessMode: core.ACCESS_BUT_IGNORE_ERR,
	JsonFiles:  []string{"./publish2.json", "./airvisual2.json", "./barchart2.json", "./bitcoinaverage.json", "./coingecko.json", "./currency.json"},
}

func InitConfig(ctx *cli.Context) (*PublishConfig, error) {
	cf := ctx.String(cmd.GetFlagName(cmd.ConfigfileFlag))
	if _, err := os.Stat(cf); os.IsNotExist(err) {
		// if there's no config file, use default config
		return &DefPublishConfig, nil
	}
	file, err := os.Open(cf)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	cfg := &PublishConfig{}
	err = json.Unmarshal(bs, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func SendPublushRequest(list []*core.PublishAPI, reset bool, pubConfig *PublishConfig) error {
	log.InitLog(log.DebugLog, log.Stdout)
	var err error
	for _, l := range list {
		if reset {
			api, err := dao.DefOscoreApiDB.QueryApiBasicInfoByTitle(nil, l.Name)
			if err != nil && !dao.IsErrNoRows(err) {
				return err
			}

			if err == nil {
				tx, errl := dao.DefOscoreApiDB.DB.Beginx()
				if errl != nil {
					return errl
				}
				defer func() {
					if errl != nil {
						tx.Rollback()
					}
				}()

				errl = dao.DefOscoreApiDB.DeleteApiBasicInfoByApiId(nil, api.ApiId)
				if errl != nil {
					return errl
				}

				errl = tx.Commit()
				if errl != nil {
					return errl
				}
			}
		} else {
			_, err = dao.DefOscoreApiDB.QueryApiBasicInfoByTitle(nil, l.Name)
		}
		if reset || dao.IsErrNoRows(err) {
			err = core.PublishAPIHandleCore(l, pubConfig.OntId, pubConfig.Author, pubConfig.AccessMode, oscoreconfig.AdminUserId)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	return err
}

func PublishAPIJSON(ctx *cli.Context) error {
	OpType := ctx.Uint(utils.GetFlagName(OperationTypeFlag))
	log.InitLog(log.DebugLog, log.Stdout)

	pconfig, err := InitConfig(ctx)
	if err != nil {
		return err
	}

	oscoreconfig.DefOscoreConfig.OscoreHost = pconfig.OscoreHost
	oscoreconfig.DefOscoreConfig.RestPort = pconfig.RestPort
	log.Infof("OscoreHost %s. RestPort %d", oscoreconfig.DefOscoreConfig.OscoreHost, oscoreconfig.DefOscoreConfig.RestPort)

	PubDB, err := dao.NewOscoreApiDB(pconfig.DBConfig)
	if err != nil {
		log.Errorf("PublishAPIJSON OpType %d", OpType)
		return err
	}
	dao.DefOscoreApiDB = PubDB
	core.DefOscoreApi = core.NewOscoreApi()
	log.Infof("PublishAPIJSON OpType %d", OpType)

	if OpType == OPTYPE_PUBLISH {
		log.Infof("PublishAPIJSON OpType publish")
		for _, f := range pconfig.JsonFiles {
			list, err := core.GetPulishFunctionList(f)
			if err != nil {
				fmt.Printf("PublishAPIJSON N.0 %s", err)
				return err
			}
			data, err := json.Marshal(list)
			if err != nil {
				return err
			}
			fmt.Printf("File: %s\n", string(data))
			err = SendPublushRequest(list, false, pconfig)
			if err != nil {
				log.Errorf("PublishAPIJSON %s", err)
				return err
			}
		}
	} else if OpType == OPTYPE_TAKEORDER {
		log.Infof("PublishAPIJSON OpType admin take order")
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

		for _, spec := range pconfig.Specs {
			param := &common.TakeOrderParam{
				OntId:            pconfig.OntId,
				ApiId:            spec.ApiId,
				SpecificationsId: spec.SpecId,
			}

			user, err := dao.DefOscoreApiDB.QueryUserNameByOntId(nil, param.OntId)
			if err != nil {
				return err
			}

			orderId, err := core.DefOscoreApi.OscoreOrder.TakeOrder(param, user.Id, tables.TAKE_ORDER_ADMIN)
			if err != nil {
				log.Errorf("Order Api %v failed. %s", *param, err)
				return err
			}

			apiKey, err := dao.DefOscoreApiDB.QueryApiKeyByOrderId(nil, *orderId)
			if err != nil {
				log.Errorf("Order Api %v failed. %s", *param, err)
				return err
			}

			log.Infof("order %v . apiKey: %v success, orderId: %s", *param, *apiKey, *orderId)
		}
	} else {
		log.Infof("nothing to do. with OpType %d", OpType)
	}

	return nil
}
