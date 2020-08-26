package main

import (
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ontio/ontology-crypto/signature"
	"github.com/ontio/ontology-go-sdk"
	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/cmd"
	"github.com/ontio/oscore/core"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/restful"
	"github.com/ontio/oscore/oscoreconfig"
)

func setupAPP() *cli.App {
	app := cli.NewApp()
	app.Usage = "oscoreapi CLI"
	app.Action = startOscore
	app.Version = oscoreconfig.Version
	app.Copyright = "Copyright in 2018 The Ontology Authors"
	app.Flags = []cli.Flag{
		cmd.LogLevelFlag,
		cmd.RestPortFlag,
		cmd.NetworkIdFlag,
		cmd.ConfigfileFlag,
	}
	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}
	return app
}

func main() {
	if err := setupAPP().Run(os.Args); err != nil {
		cmd.PrintErrorMsg(err.Error())
		os.Exit(1)
	}
}

func startOscore(ctx *cli.Context) {
	initLog(ctx)
	if err := initConfig(ctx); err != nil {
		log.Errorf("[initConfig] error: %s", err)
		return
	}
	if err := initDB(ctx); err != nil {
		log.Errorf("[initDB] error: %s", err)
		return
	}
	if err := initAccount(); err != nil {
		log.Errorf("[initAccount] error: %s", err)
		return
	}
	core.DefOscoreApi = core.NewOscoreApi()
	log.Infof("config: %v\n", oscoreconfig.DefOscoreConfig)
	log.Info("QrCodeCallback:", oscoreconfig.DefOscoreConfig.QrCodeCallback)
	startServer()
	waitToExit()
}

func initAccount() error {
	pri, _ := common.HexToBytes(oscoreconfig.DefOscoreConfig.OntIdPrivate)
	acc, err := ontology_go_sdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
	if err != nil {
		return err
	}
	oscoreconfig.DefOscoreConfig.OntIdAccount = acc
	return nil
}

func initLog(ctx *cli.Context) {
	//init log module
	logLevel := ctx.GlobalInt(cmd.GetFlagName(cmd.LogLevelFlag))
	//logName := fmt.Sprintf("%s%s", oscoreconfig.LogPath, string(os.PathSeparator))
	log.InitLog(logLevel, log.Stdout)
}

func initDB(ctx *cli.Context) error {
	var dbConfig = *oscoreconfig.DefOscoreConfig.DbConfig
	db, err := dao.NewOscoreApiDB(&dbConfig)
	if err != nil {
		return err
	}
	dao.DefOscoreApiDB = db
	return nil
}

func getDBUserName() (string, error) {
	fmt.Printf("DB UserName:")
	var userName string
	n, err := fmt.Scanln(&userName)
	if n == 0 {
		return "", fmt.Errorf("db username is wrong")
	}
	if err != nil {
		return "", err
	}
	return userName, nil
}

// GetPassword gets password from user input
func getDBPassword() ([]byte, error) {
	fmt.Printf("DB Password:")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return passwd, nil
}

func initConfig(ctx *cli.Context) error {
	//init config
	return cmd.SetOntologyConfig(ctx)
}

func startServer() {
	router := restful.NewRouter()
	port := fmt.Sprintf("%d", oscoreconfig.DefOscoreConfig.RestPort)
	log.Infof("list to port: %s", port)
	go router.Run(":" + port)
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			err := dao.DefOscoreApiDB.Close()
			if err != nil {
				log.Errorf("close db error: %s", err)
			}
			log.Infof("oscore server received exit signal: %s.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
