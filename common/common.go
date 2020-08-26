package common

import (
	"fmt"
	sdkcom "github.com/ontio/ontology-go-sdk/common"
	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

const (
	WETHER_DATA_PROCESS string = "Weather Forecast"
)

const (
	TEST_APIKEY_PREFIX   string = "test_"
	NORMAL_APIKEY_PREFIX string = "apikey_"
	OSCORE_URL_PREFIX    string = "oscoreurl_"
	ORDER_ID_PREFIX      string = "orderId_"
	QRCODE_ID_PREFIX     string = "qrcodeId_"
)

const (
	UUID_TYPE_RAW          int32 = 1
	UUID_TYPE_TEST_API_KEY int32 = 2
	UUID_TYPE_API_KEY      int32 = 3
	UUID_TYPE_OSCORE_URL   int32 = 4
	UUID_TYPE_ORDER_ID     int32 = 5
	UUID_TYPE_QRCODE_ID    int32 = 6
)

func GenerateUUId(uuidType int32) string {
	u1 := uuid.NewV4()
	switch uuidType {
	case UUID_TYPE_RAW:
		return u1.String()
	case UUID_TYPE_TEST_API_KEY:
		return TEST_APIKEY_PREFIX + u1.String()
	case UUID_TYPE_OSCORE_URL:
		return OSCORE_URL_PREFIX + u1.String()
	case UUID_TYPE_API_KEY:
		return NORMAL_APIKEY_PREFIX + u1.String()
	case UUID_TYPE_ORDER_ID:
		return ORDER_ID_PREFIX + u1.String()
	case UUID_TYPE_QRCODE_ID:
		return QRCODE_ID_PREFIX + u1.String()
	}

	return u1.String()
}

func IsTestKey(apiKey string) bool {
	return strings.HasPrefix(apiKey, TEST_APIKEY_PREFIX)
}

func IsApiKey(apiKey string) bool {
	return strings.HasPrefix(apiKey, NORMAL_APIKEY_PREFIX)
}

func IsOscoreUrlKey(oscoreUrl string) bool {
	return strings.HasPrefix(oscoreUrl, OSCORE_URL_PREFIX)
}

func GetLayer2EventByTxHash(txHash string) (*sdkcom.SmartContactEvent, error) {
	var events *sdkcom.SmartContactEvent
	var err error
	var count uint32
	for {
		events, err = oscoreconfig.DefOscoreConfig.Layer2Sdk.GetSmartContractEvent(txHash)
		if err != nil {
			log.Errorf("GetLayer2EventByTxHash N.0 :%s\n", err)
			return nil, err
		}

		if events == nil && count < 30 {
			time.Sleep(time.Second)
			count++
			continue
		}

		break
	}

	if events != nil {
		if events.State == 0 {
			return nil, fmt.Errorf("error in events.State is 0 failed.")
		}

		return events, nil
	} else {
		return nil, fmt.Errorf("GetLayer2EventByTxHash failed counter over 30 times")
	}
}
