package dao

import (
	"database/sql"
	"fmt"
	"github.com/ontio/oscore/common"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestApiDB_TmpInsert(t *testing.T) {
	TestDB, err := NewOscoreApiDB(oscoreconfig.DefDBConfigMap[oscoreconfig.NETWORK_ID_TRAVIS_NET])
	assert.Nil(t, err)
	tx, err := TestDB.DB.Beginx()
	assert.Nil(t, err)

	TestDB.ClearAll()

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
		Author:              "steven",
	}

	// insert.
	err = TestDB.InsertApiBasicInfo(tx, []*tables.ApiBasicInfo{info2})
	assert.Nil(t, err)

	// try query with tx.
	infoResult, err := TestDB.QueryApiBasicInfoByOscoreUrlKey(tx, info2.ApiOscoreUrlKey, info2.ApiState)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	infoResult, err = TestDB.QueryApiBasicInfoByApiId(tx, infoResult.ApiId, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)
	assert.Equal(t, infoResult.Author, info2.Author)

	// try query with db.
	infoResult, err = TestDB.QueryApiBasicInfoByOscoreUrlKey(nil, info2.ApiOscoreUrlKey, info2.ApiState)
	assert.Equal(t, err, sql.ErrNoRows)

	err = tx.Commit()
	assert.Nil(t, err)

	// try query with db again.
	infoResult, err = TestDB.QueryApiBasicInfoByOscoreUrlKey(nil, info2.ApiOscoreUrlKey, info2.ApiState)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	infoResult, err = TestDB.QueryApiBasicInfoByApiId(nil, infoResult.ApiId, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	assert.Equal(t, infoResult.ApplicationScenario, info2.ApplicationScenario)

	l := 11
	infos := make([]*tables.ApiBasicInfo, l)
	for i := 0; i < len(infos); i++ {
		info := &tables.ApiBasicInfo{
			Icon:                "",
			Title:               "mytestasd" + strconv.Itoa(i),
			ApiProvider:         common.GenerateUUId(1),
			ApiOscoreUrlKey:     common.GenerateUUId(1),
			ApiUrl:              "",
			Price:               "",
			ApiDesc:             "",
			Specifications:      1,
			ApiState:            tables.API_STATE_BUILTIN,
			Popularity:          0,
			Delay:               0,
			SuccessRate:         0,
			InvokeFrequency:     0,
			ApplicationScenario: common.GenerateUUId(1),
		}

		info.ApiProvider = common.GenerateUUId(1)
		infos[i] = info
	}
	err = TestDB.InsertApiBasicInfo(nil, infos)
	assert.Nil(t, err)

	for i := 0; i < len(infos); i++ {
		infoResult, err := TestDB.QueryApiBasicInfoByOscoreUrlKey(nil, infos[i].ApiOscoreUrlKey, infos[i].ApiState)
		assert.Nil(t, err)
		assert.Equal(t, infoResult.ApplicationScenario, infos[i].ApplicationScenario)
	}

	tx, err = TestDB.DB.Beginx()
	assert.Nil(t, err)
	// test SearchApi
	res, err := TestDB.SearchApi(tx)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(res["newest"]))
	assert.Equal(t, 10, len(res["hottest"]))

	// test RequestParam
	infoResult, err = TestDB.QueryApiBasicInfoByOscoreUrlKey(tx, info2.ApiOscoreUrlKey, info2.ApiState)
	assert.Nil(t, err)
	params := &tables.RequestParam{
		ApiId:      infoResult.ApiId,
		ParamName:  "",
		Required:   true,
		ParamType:  "",
		ParamWhere: tables.URL_PARAM_RESTFUL,
		Note:       "zzzz",
		ValueDesc:  "ValueDesc",
	}

	assert.Nil(t, TestDB.InsertRequestParam(tx, []*tables.RequestParam{params}))
	paramResult, err := TestDB.QueryRequestParamByApiId(tx, infoResult.ApiId, 0)
	assert.Nil(t, err)
	assert.Equal(t, params.ValueDesc, paramResult[0].ValueDesc)
	assert.Equal(t, params.Required, paramResult[0].Required)

	//TestDB.InsertTag

	err = tx.Commit()
	assert.Nil(t, err)

	count, err := TestDB.QueryApiBasicInfoCount(nil, tables.API_STATE_BUILTIN)
	assert.Nil(t, err)
	assert.Equal(t, uint64(12), count)

	err = TestDB.ClearRequestParamDB()
	err = TestDB.ClearApiBasicDB()
	assert.Nil(t, err)
	fmt.Printf("TestApiDB_TmpInsert done")
}
