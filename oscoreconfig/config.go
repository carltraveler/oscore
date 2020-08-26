package oscoreconfig

import (
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common/log"
)

var Version = ""

var (
	DEFAULT_LOG_LEVEL = log.InfoLog
	DEFAULT_REST_PORT = uint(8080)
)

type Config struct {
	NetWorkId             uint      `json:"network_id"`
	RestPort              uint      `json:"rest_port"`
	Version               string    `json:"version"`
	DbConfig              *DBConfig `json:"db_config"`
	OperatorPublicKey     string    `json:"operator_public_key"`
	QrCodeCallback        string    `json:"qrcode_callback"`
	NASAAPIKey            string    `json:"nasa_api_key"`
	Collect_Money_Address string    `json:"collect_money_address"`
	NetType               string    `json:"net_type"`
	OntId                 string    `json:"ont_id"`
	OntIdPrivate          string    `json:"ontid_private"`
	OscoreHost            string    `json:"oscore_host"`
	AliRetrunHost         string    `json:"aliRetrunHost"`
	DDXFContractAddress   string    `json:"ddxf_contract_address"`
	DDXFAPIServer         string    `json:"ddxfapiServer"`
	AliPayAddress         string    `json:"aliAayAddress"`
	AliQueryAddress       string    `json:"aliQueryAddress"`
	AccessKeyR            string    `json:"accessKeyR"`
	SecretKeyR            string    `json:"secretKeyR"`
	AccessKeyV            string    `json:"accessKeyV"`
	SecretKeyV            string    `json:"secretKeyV"`
	AppID                 string    `json:"appId"`
	WalletName            string    `json:"walletName"`
	Layer2Contract        string    `json:"layer2Contract"`
	Layer2MainNetNode     string    `json:"layer2MainNetNode"`
	Layer2TestNetNode     string    `json:"layer2TestNetNode"`
	Layer2RecordInterval  int       `json:"layer2RecordInterval"`
	OntSdk                *ontology_go_sdk.OntologySdk
	Layer2Sdk             *ontology_go_sdk.OntologySdk
	OscoreAccount         *ontology_go_sdk.Account
	OntIdAccount          *ontology_go_sdk.Account
}

type DBConfig struct {
	ProjectDBUrl      string `json:"projectdb_url"`
	ProjectDBUser     string `json:"projectdb_user"`
	ProjectDBPassword string `json:"projectdb_password"`
	ProjectDBName     string `json:"projectdb_name"`
}

var DefDBConfigMap = map[int]*DBConfig{
	NETWORK_ID_SOLO_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1:3306",
		ProjectDBUser:     "steven",
		ProjectDBPassword: "abcd1234",
		ProjectDBName:     "oscore",
	},
	NETWORK_ID_MAIN_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1:3306",
		ProjectDBUser:     "root",
		ProjectDBPassword: "111111",
		ProjectDBName:     "oscore",
	},
	NETWORK_ID_TRAVIS_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1",
		ProjectDBUser:     "root",
		ProjectDBPassword: "",
		ProjectDBName:     "oscoreunittest",
	},
	NETWORK_ID_UNIT_TEST_NET: &DBConfig{
		ProjectDBUrl:      "127.0.0.1",
		ProjectDBUser:     "steven",
		ProjectDBPassword: "abcd1234",
		ProjectDBName:     "oscoreunittest",
	},
}

var DefOscoreConfig = &Config{
	RestPort:              DEFAULT_REST_PORT,
	Version:               "1.0.0",
	NetWorkId:             NETWORK_ID_POLARIS_NET,
	DbConfig:              DefDBConfigMap[NETWORK_ID_SOLO_NET],
	OperatorPublicKey:     "02b8fcf42deecc7cccb574ba145f2f627339fbd3ba2b63fda99af0a26a8d5a01da",
	QrCodeCallback:        "http://192.168.1.175:8080/api/v1/ali/sendTxAli",
	OscoreHost:            "http://192.168.1.175",
	NASAAPIKey:            NASA_API_KEY,
	OntId:                 OntId,
	OntIdPrivate:          OntIdPrivate,
	Collect_Money_Address: Collect_Money_Address,
	DDXFAPIServer:         "http://106.75.209.209:8080",
}
