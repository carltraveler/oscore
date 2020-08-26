package dao

import (
	"github.com/ontio/oscore/common"
	"github.com/ontio/oscore/models/tables"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestApiDB_InsertApiBasicInfo(t *testing.T) {
	TestDB.ClearAll()
	l := 11
	infos := make([]*tables.ApiBasicInfo, l)
	for i := 0; i < len(infos); i++ {
		info := &tables.ApiBasicInfo{
			Icon:            "",
			Title:           "mytestasdx" + strconv.Itoa(i),
			ApiProvider:     common.GenerateUUId(common.UUID_TYPE_RAW),
			ApiOscoreUrlKey: common.GenerateUUId(common.UUID_TYPE_RAW),
			ApiUrl:          "",
			Price:           "",
			ApiState:        int32(1),
			ApiDesc:         "",
			Specifications:  1,
			Popularity:      0,
			Delay:           0,
			SuccessRate:     0,
			InvokeFrequency: 0,
		}

		infos[i] = info
	}
	assert.Nil(t, TestDB.InsertApiBasicInfo(nil, infos))

	ApiType := "ApiType"
	Icon := "Icon"
	Title := "Title"
	ApiProvider := "ApiProvider"
	ApiOscoreUrlKey := "ApiOscoreUrlKey"
	ApiUrl := "ApiUrl"
	Price := "Price"
	ApiDesc := "ApiDesc"
	ErrorDesc := "ErrorDesc"
	Specifications := uint32(9)
	Popularity := uint32(10)
	Delay := uint32(11)
	SuccessRate := uint32(100)
	InvokeFrequency := uint64(12345)
	ApiState := int32(1)
	RequestType := "RequestType"
	Mark := "Mark"
	ResponseParam := "ResponseParam"
	ResponseExample := "ResponseExample"
	DataDesc := "DataDesc"
	DataSource := "DataSource"
	ApplicationScenario := "ApplicationScenario"

	info2 := &tables.ApiBasicInfo{
		ApiType:             ApiType,
		Icon:                Icon,
		Title:               Title,
		ApiProvider:         ApiProvider,
		ApiOscoreUrlKey:     ApiOscoreUrlKey,
		ApiUrl:              ApiUrl,
		Price:               Price,
		ApiDesc:             ApiDesc,
		ErrorDesc:           ErrorDesc,
		Specifications:      Specifications,
		Popularity:          Popularity,
		Delay:               Delay,
		SuccessRate:         SuccessRate,
		InvokeFrequency:     InvokeFrequency,
		ApiState:            ApiState,
		RequestType:         RequestType,
		Mark:                Mark,
		ResponseParam:       ResponseParam,
		ResponseExample:     ResponseExample,
		DataDesc:            DataDesc,
		DataSource:          DataSource,
		ApplicationScenario: ApplicationScenario,
	}

	assert.Nil(t, TestDB.InsertApiBasicInfo(nil, []*tables.ApiBasicInfo{info2}))

	basic, err := TestDB.QueryApiBasicInfoByPage(1, 1, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	info, err := TestDB.QueryApiBasicInfoByApiId(nil, basic[0].ApiId, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	assert.Equal(t, info.ApiId, basic[0].ApiId)

	assert.Nil(t, err)
	params := &tables.RequestParam{
		ApiId:      info.ApiId,
		ParamName:  "qqq",
		Required:   true,
		ParamType:  "",
		ParamWhere: tables.URL_PARAM_RESTFUL,
		Note:       "xxx",
		ValueDesc:  "yyy",
	}

	assert.Nil(t, TestDB.InsertRequestParam(nil, []*tables.RequestParam{params}))

	requestParams, err := TestDB.QueryRequestParamByApiId(nil, info.ApiId, tables.REQUEST_PARAM_TAG_PARAM)
	assert.Nil(t, err)
	assert.Equal(t, len(requestParams), 1)

	infos, err = TestDB.QueryApiBasicInfoByCategoryId(nil, 1, 0, 1)
	assert.Nil(t, err)

	err = TestDB.ClearRequestParamDB()
	assert.Nil(t, err)
	err = TestDB.ClearApiBasicDB()
	assert.Nil(t, err)
}

func InsertTestApi() ([]*tables.ApiBasicInfo, error) {
	ApiType := "ApiType"
	Icon := "Icon"
	Title := "Title"
	ApiProvider := "ApiProvider"
	ApiOscoreUrlKey := "ApiOscoreUrlKey"
	ApiUrl := "ApiUrl"
	Price := "Price"
	ApiDesc := "ApiDesc"
	ErrorDesc := "ErrorDesc"
	Specifications := uint32(9)
	Popularity := uint32(10)
	Delay := uint32(11)
	SuccessRate := uint32(100)
	InvokeFrequency := uint64(12345)
	ApiState := tables.API_STATE_BUILTIN
	RequestType := "RequestType"
	Mark := "Mark"
	ResponseParam := "ResponseParam"
	ResponseExample := "ResponseExample"
	DataDesc := "DataDesc"
	DataSource := "DataSource"
	ApplicationScenario := "ApplicationScenario"

	info2 := &tables.ApiBasicInfo{
		ApiType:             ApiType,
		Icon:                Icon,
		Title:               Title,
		ApiProvider:         ApiProvider,
		ApiOscoreUrlKey:     ApiOscoreUrlKey,
		ApiUrl:              ApiUrl,
		Price:               Price,
		ApiDesc:             ApiDesc,
		ErrorDesc:           ErrorDesc,
		Specifications:      Specifications,
		Popularity:          Popularity,
		Delay:               Delay,
		SuccessRate:         SuccessRate,
		InvokeFrequency:     InvokeFrequency,
		ApiState:            ApiState,
		RequestType:         RequestType,
		Mark:                Mark,
		ResponseParam:       ResponseParam,
		ResponseExample:     ResponseExample,
		DataDesc:            DataDesc,
		DataSource:          DataSource,
		ApplicationScenario: ApplicationScenario,
	}

	err := TestDB.InsertApiBasicInfo(nil, []*tables.ApiBasicInfo{info2})
	if err != nil {
		return nil, err
	}

	basic, err := TestDB.QueryApiBasicInfoByPage(0, 10, tables.API_STATE_BUILTIN)
	if err != nil {
		return nil, err
	}
	return basic, nil
}

func TestApiDB_InsertApiKey(t *testing.T) {
	TestDB.ClearApiKeyDB()
	TestDB.ClearOrderDB()
	TestDB.ClearApiBasicDB()
	basic, err := InsertTestApi()
	assert.Nil(t, err)
	orderId := "orderId"
	tt := time.Now().Unix()
	order := &tables.Order{
		ApiId:     basic[0].ApiId,
		OrderId:   orderId,
		OntId:     "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
		OrderTime: tt,
	}
	err = TestDB.InsertOrder(nil, order)
	assert.Nil(t, err)
	apikey := "apikey"
	key := &tables.APIKey{
		ApiId:        basic[0].ApiId,
		OrderId:      orderId,
		ApiKey:       apikey,
		RequestLimit: 2,
		UsedNum:      1,
		OntId:        "did:ont:APe4yT5B6KnvR7LenkZD6eQGhG52Qrdjuo",
	}

	err = TestDB.InsertApiKey(nil, key)
	assert.Nil(t, err)

	key, err = TestDB.QueryApiKeyByApiKey(nil, apikey)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), key.UsedNum)

	TestDB.ClearApiKeyDB()
	TestDB.ClearOrderDB()
	TestDB.ClearApiBasicDB()
}

func TestApiDB_QuerySpecificationsByApiDetailId(t *testing.T) {
	assert.Nil(t, TestDB.ClearSpecificationsDB())
	basic, err := InsertTestApi()
	assert.Nil(t, err)
	params := []*tables.Specifications{
		&tables.Specifications{
			ApiId:  basic[0].ApiId,
			Price:  "0",
			Amount: 500,
		},
		&tables.Specifications{
			ApiId:  basic[0].ApiId,
			Price:  "0.01",
			Amount: 1000,
		},
	}
	err = TestDB.InsertSpecifications(nil, params)
	assert.Nil(t, err)

	specs, err := TestDB.QuerySpecificationsByApiId(nil, basic[0].ApiId)
	assert.Nil(t, err)
	assert.Equal(t, specs[0].ApiId, basic[0].ApiId)

	spec, err := TestDB.QuerySpecificationsById(nil, specs[0].Id)
	assert.Nil(t, err)

	assert.Equal(t, spec.Id, specs[0].Id)

	assert.Nil(t, TestDB.ClearSpecificationsDB())
	assert.Nil(t, TestDB.ClearApiBasicDB())
}

func TestApiDB_QueryApiBasicInfoByOscoreUrlKey(t *testing.T) {
	TestDB.ClearAll()
	info := &tables.ApiBasicInfo{
		Icon:            "",
		Title:           "mytestasd",
		ApiOscoreUrlKey: common.GenerateUUId(common.UUID_TYPE_OSCORE_URL),
		ApiProvider:     "",
		ApiUrl:          "",
		Price:           "",
		ApiState:        int32(1),
		ApiDesc:         "",
		Specifications:  1,
		Popularity:      0,
		Delay:           0,
		SuccessRate:     0,
		InvokeFrequency: 0,
	}

	assert.Nil(t, TestDB.InsertApiBasicInfo(nil, []*tables.ApiBasicInfo{info}))
	_, err := TestDB.QueryApiBasicInfoByOscoreUrlKey(nil, info.ApiOscoreUrlKey, info.ApiState)
	assert.Nil(t, err)
	assert.Nil(t, TestDB.ClearApiBasicDB())
}
