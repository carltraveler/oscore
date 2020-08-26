package core

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/ontio/ontology-crypto/signature"
	ontSdk "github.com/ontio/ontology-go-sdk"
	common3 "github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/aksk"
	common2 "github.com/ontio/oscore/common"
	"github.com/ontio/oscore/core/http"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/ddxf"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/ontio/oscore/utils"
	"net/url"
	"strings"
	"time"
)

const (
	TRADE_SUCCESS  string = "TRADE_SUCCESS"
	TRADE_FINISHED string = "TRADE_FINISHED"
	TRADE_CLOSED   string = "TRADE_CLOSED"
	WAIT_BUYER_PAY string = "WAIT_BUYER_PAY"
)

type AliPayCallBackArg struct {
	OutTradeNo string `json:"out_trade_no"`
	Status     string `json:"status"`
}

type AliPayRefundArg struct {
	Amount     string `json:"amount"`
	App        int    `json:"app"`
	OutTradeNo string `json:"out_trade_no"`
	Reason     string `json:"reason"`
}

func RefundFlow(order *tables.Order) {
	for {
		log.Infof("RundFlow. Y. 0 %v", *order)
		req := AliPayRefundArg{
			Amount:     order.Amount,
			App:        0,
			OutTradeNo: order.OrderId,
			Reason:     "oscore api payment",
		}
		b, err := json.Marshal(req)
		if err != nil {
			log.Errorf("RundFlow.N.0 %s. %v", err, *order)
			continue
		}

		target := oscoreconfig.DefOscoreConfig.AliPayAddress + "/aksk/alipay/refund"
		ru, err := url.Parse(target)
		if err != nil {
			log.Errorf("RundFlow.N.1 %s. %v", err, *order)
			continue
		}
		_, err = aksk.Post(ru, b)
		if err != nil {
			log.Errorf("RundFlow.N.2 %s %v", err, *order)
			continue
		}

		err = fmt.Errorf("refund db start")
		for err != nil {
			if order.State == tables.ORDER_STATE_CANCEL {
				err = dao.DefOscoreApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, order.Result, tables.ORDER_STATE_CANCEL_REFUND_DONE, time.Now().Unix())
			} else {
				err = dao.DefOscoreApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, order.Result, tables.ORDER_STATE_DEL_REFUND_DONE, time.Now().Unix())
			}
			log.Errorf("RundFlow.N.3 refund done but update order err %s %v", err, *order)
		}

		log.Infof("RundFlow.N.3 refund done %v", *order)

		break
	}
}

const (
	AlreadyNotifyHandle string = "AlreadyNotifyHandle"
)

// before enter must be ORDER_STATE_PAYING or ORDER_STATE_CANCEL, and ali trade status is TRADE_SUCCESS
func SendTxAliCore(param *AliPayCallBackArg) error {
	log.Debugf("SendTxAliCore %v", *param)
	DefOscoreApi.Cache.CallBackLock.Lock()
	order, err := dao.DefOscoreApiDB.QueryOrderByOrderId(nil, param.OutTradeNo)
	if err != nil {
		DefOscoreApi.Cache.CallBackLock.Unlock()
		return err
	}

	if order.Result == AlreadyNotifyHandle || order.State == tables.ORDER_STATE_COMPLETE || order.State >= tables.ORDER_STATE_CANCEL_REFUNDING {
		log.Infof("SendTxAliCore already handled orderId %s", order.OrderId)
		DefOscoreApi.Cache.CallBackLock.Unlock()
		return nil
	}

	err = dao.DefOscoreApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, AlreadyNotifyHandle, tables.ORDER_STATE_CALLBACK_HANDLING, time.Now().Unix())
	if err != nil {
		log.Debugf("SendTxAliCore .N.first %s", err)
		DefOscoreApi.Cache.CallBackLock.Unlock()
		return err
	}
	DefOscoreApi.Cache.CallBackLock.Unlock()

	switch param.Status {
	case TRADE_CLOSED:
		log.Warnf("SendTxAli Y.0 %s", "received TRADE_CLOSED statues")
		return nil
	case WAIT_BUYER_PAY:
		log.Warnf("SendTxAli Y.0 %s", "received WAIT_BUYER_PAY statues")
		return nil
	case TRADE_FINISHED:
		log.Warnf("SendTxAli Y.0 %s", "received TRADE_FINISHED statues")
		return nil
	case TRADE_SUCCESS:
		// handle canceled or api removed logic.
		if order.State == tables.ORDER_STATE_CANCEL {
			err = dao.DefOscoreApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, order.Result, tables.ORDER_STATE_CANCEL_REFUNDING, time.Now().Unix())
			log.Infof("SendTxAli N.refund canceled %s", err)
			go RefundFlow(order)
			return nil
		}
		if order.OrderType == oscoreconfig.Api {
			_, err = dao.DefOscoreApiDB.QueryApiBasicInfoByApiId(nil, order.ApiId, tables.API_STATE_BUILTIN)
			if err != nil {
				err = dao.DefOscoreApiDB.UpdateTxInfoByOrderId(nil, order.OrderId, order.Result, tables.ORDER_STATE_DEL_REFUNDING, time.Now().Unix())
				log.Infof("SendTxAli N.refund del %s", err)
				go RefundFlow(order)
				return nil
			}
		}
	default:
		log.Errorf("SendTxAli Y.0 %s", "received Other statues")
		return nil
	}

	txdb, errl := dao.DefOscoreApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("SendTx.N.1: %s", errl)
		return errl
	}

	defer func() {
		if errl != nil {
			txdb.Rollback()
		}
	}()

	errl = CoreOrderCallBackHandle(order, txdb)
	if errl != nil {
		log.Errorf("SendTxAli CoreOrderCallBackHandle End %s", errl)
		return errl
	}

	errl = txdb.Commit()
	if errl != nil {
		log.Debugf("SendTx.N.13 %v", errl)
		return errl
	}
	return nil
}

func CoreOrderCallBackHandle(order *tables.Order, txdb *sqlx.Tx) error {
	if txdb == nil {
		return fmt.Errorf("CoreOrderCallBackHandle txdb can not nil")
	}
	log.Debugf("CoreOrderCallBackHandle start Y.0 %v", *order)
	var err error
	arr := strings.Split(order.OntId, ":")
	if len(arr) < 3 {
		log.Errorf("SendTxAli strart N.0.0 error order Ontid")
		return err
	}

	result := ""
	apiKey := ""
	txHash := ""
	if order.OrderType == oscoreconfig.Api {
		switch order.OrderKind {
		case tables.ORDER_KIND_API:
			apiKey, err = generateApiKey(txdb, order.OrderId, order.OntId)
			if err != nil {
				log.Debugf("SendTx.N.2: %s", err)
				return err
			}
		case tables.ORDER_KIND_API_RENEW:
			apiKey, err = RenewApiKey(txdb, order, order.OntId)
			if err != nil {
				log.Debugf("SendTx.N.2.0: %s", err)
				return err
			}
			// here add to old apk key statistics.
		default:
			return fmt.Errorf("SendTx.N.2.1: error order kind %d", order.OrderKind)
		}

		recordApiKey, err := dao.DefOscoreApiDB.QueryApiKeyByApiKey(txdb, apiKey)
		if err != nil {
			log.Errorf("SendTx.N.2.2: %s", err)
		}

		var count uint32
		for {
			txHash, err = ddxf.BuyerRecordDDXF(recordApiKey)
			if err != nil && count < 4 {
				log.Errorf("CoreOrderCallBackHandle .N.2.3: %s", err)
				count += 1
				continue
			}

			log.Infof("CoreOrderCallBackHandle BuyerRecordDDXF txHash %s", txHash)
			break
		}
	} else if order.OrderType == oscoreconfig.ApiProcess {
		// send request to server. check the result. and update status.
		// may check other request like coin. now default is WetherForcastServiceRequest
		api, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiId(txdb, order.ApiId, tables.API_STATE_BUILTIN)
		if err != nil {
			log.Debugf("SendTx.N.3: %s", err)
			return err
		}

		if order.OrderKind == tables.ORDER_KIND_DATA_PROCESS_WETHER {
			paramWether := &common2.WetherForcastRequest{}
			err = json.Unmarshal([]byte(order.Request), paramWether)
			if err != nil {
				log.Debugf("SendTx.N.4: %s", err)
				return err
			}
			env, err := dao.DefOscoreApiDB.QueryEnvById(txdb, paramWether.EnvId)
			if err != nil {
				log.Debugf("SendTx.N.5: %s", err)
				return err
			}
			alg, err := dao.DefOscoreApiDB.QueryAlgorithmById(nil, paramWether.AlgorithmId)
			if err != nil {
				log.Debugf("SendTx.N.6: %s", err)
				return err
			}
			tm := time.Unix(paramWether.TargetDate, 0)
			mm, err := time.ParseDuration("-240h")
			if err != nil {
				log.Debugf("SendTx.N.7 %s", err)
				return err
			}
			targetTime := tm.Add(mm)

			user, err := dao.DefOscoreApiDB.QueryUserNameByOntId(nil, order.OntId)
			if err != nil {
				log.Debugf("SendTx.N.7.0 %s", err)
				return err
			}

			pri, err := hex.DecodeString(user.PrivateKey.String)
			if err != nil {
				log.Debugf("SendTx.N.7.0.1 %s", err)
				return err
			}

			userAccount, err := ontSdk.NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)
			if err != nil {
				log.Debugf("SendTx.N.7.1 %s", err)
				return err
			}

			addr, err := common3.AddressFromBase58(env.OwnerAddress)
			if err != nil {
				log.Debugf("SendTx.N.7.2 %s", err)
				return err
			}

			log.Infof("BuildWetherForcastTransaction agentAddr: %s OwnerAddr: %s", env.OwnerAddress, userAccount.Address.ToBase58())
			var apiTokenId string
			txHash, apiTokenId, err = BuildWetherForcastTransaction(userAccount, api.TokenHash, env.TokenHash, addr.ToHexString())
			if err != nil {
				log.Debugf("SendTx.N.7.2 %s", err)
				return err
			}

			r := common2.WetherForcastServiceRequest{
				DataUrl:       api.ApiProvider,
				Header:        make(map[string]interface{}),
				Param:         make(map[string]interface{}),
				RequestMethod: api.RequestType,
				AlgorithmName: alg.AlgName,
				TokenId:       apiTokenId,
				OwnerAddress:  userAccount.Address.ToHexString(),
			}

			r.Header["Authorization"] = "4cdb5582-90d8-11ea-af71-0242ac130002-4cdb5640-90d8-11ea-af71-0242ac130002"
			r.Param["params"] = "airTemperature"
			r.Param["lat"] = paramWether.Location.Lat
			r.Param["lng"] = paramWether.Location.Lng
			r.Param["start"] = targetTime.Format("2006-01-02") // should be string.
			log.Debugf("SendTx.Y.0 %v", r)

			data, err := json.Marshal(r)
			if err != nil {
				log.Debugf("SendTx.N.8 %v", err)
				return err
			}
			res, err := http.NewClient().Post(env.ServiceUrl, data)
			if err != nil {
				log.Debugf("SendTx.N.9 %v", err)
				return err
			}

			type DataProcessServiceRespone struct {
				Action  string      `json:"action"`
				Error   uint32      `json:"error"`
				Desc    string      `json:"desc"`
				Result  interface{} `json:"result"`
				Version string      `json:"version"`
			}
			var t DataProcessServiceRespone
			err = json.Unmarshal(res, &t)
			if err != nil {
				log.Debugf("SendTx.N.9.0 %s", err)
				return err
			}

			if t.Error != 0 {
				err = errors.New(t.Desc)
				log.Debugf("SendTx.N.9.1 %s", err)
				return err
			}

			log.Debugf("%s: %s", alg.AlgName, string(res))
			result = string(res)
		} else {
			err = errors.New("wrong data process type")
			return err
		}
	}

	err = dao.DefOscoreApiDB.UpdateTxInfoByOrderId(txdb, order.OrderId, result, tables.ORDER_STATE_COMPLETE, time.Now().Unix())
	if err != nil {
		log.Debugf("SendTx.N.11 %v", err)
		return err
	}

	err = dao.DefOscoreApiDB.UpdateOrderApiKey(txdb, txHash, order.OrderId, apiKey)
	if err != nil {
		log.Debugf("SendTx.N.12 %v", err)
		return err
	}

	return nil
}

func generateApiKey(tx *sqlx.Tx, orderId, ontId string) (string, error) {
	order, err := dao.DefOscoreApiDB.QueryOrderByOrderId(tx, orderId)
	if err != nil {
		return "", err
	}

	spec, err := dao.DefOscoreApiDB.QuerySpecificationsById(tx, order.SpecificationsId)
	if err != nil {
		return "", err
	}

	var id string
	var outDate int64
	var apiKeyType int32

	switch spec.SpecType {
	case tables.SPEC_TYPE_COUNT:
		apiKeyType = tables.API_KEY_TYPE_COUNT
	case tables.SPEC_TYPE_DURATION:
		apiKeyType = tables.API_KEY_TYPE_DURATION
	default:
		return "", errors.New("error spec type")
	}

	id = common2.GenerateUUId(common2.UUID_TYPE_API_KEY)
	apiKey := &tables.APIKey{
		OrderId:      orderId,
		ApiKey:       id,
		ApiId:        order.ApiId,
		RequestLimit: spec.Amount,
		UsedNum:      0,
		OntId:        ontId,
		UserId:       order.UserId,
		OutDate:      outDate,
		ApiKeyType:   apiKeyType,
		CreateTime:   time.Now().Unix(),
	}

	outDate, err = utils.CalOutDateByMonth(apiKey, spec.EffectiveDuration)
	if err != nil {
		return "", err
	}

	apiKey.OutDate = outDate

	return apiKey.ApiKey, dao.DefOscoreApiDB.InsertApiKey(tx, apiKey)
}

func RenewApiKey(tx *sqlx.Tx, order *tables.Order, ontId string) (string, error) {
	// use order.Request to specify the new renewOrder related apikey. here use Request not ApiKey. because ApiKey means already payed.
	key, err := dao.DefOscoreApiDB.QueryApiKeyByApiKey(tx, order.Request)
	if err != nil {
		log.Debugf("RenewApiKey : N.0 %s", err)
		return "", err
	}
	// note key's other info like key.used_num etc is not right. may be old.

	spec, err := dao.DefOscoreApiDB.QuerySpecificationsById(nil, order.SpecificationsId)
	if err != nil {
		log.Debugf("RenewApiKey : N.1 %s", err)
		return "", err
	}

	switch spec.SpecType {
	case tables.SPEC_TYPE_COUNT:
		// here must load apikey from cache to get new UsedNum OutDate etc new info.
		return key.ApiKey, DefOscoreApi.Cache.AtomicRenewApiKeyCacheCount(key.ApiKey, spec.Amount)
	case tables.SPEC_TYPE_DURATION:
		return key.ApiKey, DefOscoreApi.Cache.AtomicRenewApiKeyCacheDate(key.ApiKey, spec.EffectiveDuration)
	default:
		return "", errors.New("error spec type.")
	}
}
