package common

import (
	"encoding/json"
	"fmt"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/restful/api/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuildApiBasicInfo(t *testing.T) {
	res := common.ResponseSuccessOnto()
	bs, _ := json.Marshal(res)
	fmt.Println(string(bs))
}

func TestGenerateOrderId(t *testing.T) {
	orderId := make([]string, 0)
	for i := 0; i < 10; i++ {
		orderId = append(orderId, GenerateUUId(UUID_TYPE_RAW))
		fmt.Println(orderId)
	}

	for i := 0; i < 10; i++ {
		for j := i + 1; j < 10; j++ {
			assert.NotEqual(t, orderId[i], orderId[j])
		}
	}
}

func TestBuildQrCodeResult(t *testing.T) {
	res := ApiDetailResponse{
		ApiId: 1,
		RequestParams: []*tables.RequestParam{
			&tables.RequestParam{
				ApiId:     1,
				ParamName: "aa",
			},
		},
		ErrorCodes: []*tables.ErrorCode{
			&tables.ErrorCode{
				ErrorCode: 2,
			},
		},
	}

	bs, _ := json.Marshal(res)
	fmt.Println(string(bs))
}
