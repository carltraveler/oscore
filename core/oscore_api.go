package core

import (
	"encoding/json"
	"errors"
	"time"

	"fmt"

	"github.com/ontio/ontology/common/log"
	"github.com/ontio/oscore/common"
	"github.com/ontio/oscore/core/freq"
	"github.com/ontio/oscore/core/http"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
)

var DefOscoreApi *OscoreApi

type OscoreApi struct {
	OscoreOrder *OscoreOrder
	Cache       *freq.DBCache
}

func NewOscoreApi() *OscoreApi {
	http.DefClient = http.NewClient()
	cache := freq.NewDBCache()
	return &OscoreApi{
		OscoreOrder: NewOscoreOrder(),
		Cache:       cache,
	}
}

func (this *OscoreApi) GenerateApiTestKey(apiId uint32, ontid string, apiState int32) (*tables.APIKey, error) {
	tx, errl := dao.DefOscoreApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("GenerateApiTestKey.0. %s", errl)
		return nil, errl
	}

	_, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiId(tx, apiId, apiState)
	if err != nil {
		errl = err
		log.Debugf("GenerateApiTestKey.1. %s", err)
		return nil, err
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	testKey, err := dao.DefOscoreApiDB.QueryApiTestKeyByOntidAndApiId(tx, ontid, apiId)
	if err != nil && !dao.IsErrNoRows(err) {
		errl = err
		log.Debugf("GenerateApiTestKey.3. %s", err)
		return nil, err
	} else if err == nil {
		return testKey, nil
	} else {
		apiKey := &tables.APIKey{
			ApiKey:       common.GenerateUUId(common.UUID_TYPE_TEST_API_KEY),
			ApiId:        apiId,
			RequestLimit: oscoreconfig.DefRequestLimit,
			UsedNum:      0,
			OntId:        ontid,
			ApiKeyType:   tables.API_KEY_TYPE_COUNT,
			CreateTime:   time.Now().Unix(),
		}
		err = dao.DefOscoreApiDB.InsertApiTestKey(tx, apiKey)
		if err != nil {
			errl = err
			return nil, err
		}

		err = tx.Commit()
		if err != nil {
			errl = err
			return nil, err
		}
		return apiKey, nil
	}
}

func (this *OscoreApi) AdminTestApi(params []*tables.RequestParam, apiId uint32) ([]byte, error) {
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiId != params[i+1].ApiId {
			return nil, errors.New("params should to the same api")
		}
		if params[i].Required && params[i].ValueDesc == "" {
			return nil, fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	info, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiId(nil, apiId, tables.API_STATE_PUBLISH)
	if err != nil {
		return nil, err
	}

	data, err := HandleDataSourceReqCore(nil, info.ApiOscoreUrlKey, params, "", true)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *OscoreApi) TestApiKey(params []*tables.RequestParam, apiKey string) ([]byte, error) {
	for i, _ := range params {
		if (i != len(params)-1) && params[i].ApiId != params[i+1].ApiId {
			return nil, errors.New("params should to the same api")
		}
		if params[i].Required && params[i].ValueDesc == "" {
			return nil, fmt.Errorf("param:%s is nil", params[i].ParamName)
		}
	}

	apiTestKey := apiKey

	key, err := dao.DefOscoreApiDB.QueryApiKeyByApiKey(nil, apiTestKey)
	if err != nil {
		return nil, err
	}

	apiId := key.ApiId
	info, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiIdInState(nil, apiId, []int32{tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER})
	if err != nil {
		return nil, err
	}

	data, err := HandleDataSourceReqCore(nil, info.ApiOscoreUrlKey, params, apiTestKey, false)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *OscoreApi) QueryBasicApiInfoByPage(pageNum, pageSize uint32, apiState int32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefOscoreApiDB.QueryApiBasicInfoByPage(start, pageSize, apiState)
}

const (
	NOTIFICATION_API_SHELF        int32 = 32
	NOTIFICATION_API_SELLER_SHELF int32 = 42
)

func (this *OscoreApi) DelPulishApiCore(OntId string, apiId uint32) error {
	keys, errl := dao.DefOscoreApiDB.QueryApiKeyByApiId(nil, apiId)
	if errl != nil {
		log.Debugf("DelPulishApiCore.N.0. %s", errl)
		return errl
	}

	tx, errl := dao.DefOscoreApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("DelPulishApiCore.N.1. %s", errl)
		return errl
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	var keycount int
	var delState int32

	for _, key := range keys {
		// key.OutDate default is 0.
		if ((key.ApiKeyType == tables.API_KEY_TYPE_COUNT && key.RequestLimit > key.UsedNum) || (key.ApiKeyType == tables.API_KEY_TYPE_DURATION && key.OutDate > time.Now().Unix())) && key.UserId != "" {
			notify := &tables.Notification{
				Id:      common.GenerateUUId(common.UUID_TYPE_RAW),
				UserId:  key.UserId,
				Type:    NOTIFICATION_API_SHELF,
				KeyWord: key.OrderId,
			}
			errl = dao.DefOscoreApiDB.InsertNotifycations(tx, []*tables.Notification{notify})
			if errl != nil {
				log.Debugf("DelPulishApiCore.N.2. %s", errl)
				return errl
			}

			keycount += 1
		}
	}

	if keycount != 0 {
		delState = tables.API_STATE_DISABLE_ORDER
	} else {
		delState = tables.API_STATE_DELETE
		info, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiId(tx, apiId, tables.API_STATE_IGNOR)
		if err != nil {
			log.Debugf("DelPulishApiCore.N.2.0. %s", errl)
			errl = err
			return errl
		}
		notify := &tables.Notification{
			Id:      common.GenerateUUId(common.UUID_TYPE_RAW),
			UserId:  info.UserId,
			Type:    NOTIFICATION_API_SELLER_SHELF,
			KeyWord: info.Title,
		}
		errl = dao.DefOscoreApiDB.InsertNotifycations(tx, []*tables.Notification{notify})
		if errl != nil {
			log.Debugf("DelPulishApiCore.N.2.1 %s", errl)
			return errl
		}

		errl = dao.DefOscoreApiDB.ApiBasicUpateApiNotifyDelete(tx, 1, uint32(apiId))
		if errl != nil {
			log.Debugf("DelPulishApiCore.N.2.2 %s", errl)
			return errl
		}
	}

	errl = dao.DefOscoreApiDB.ApiBasicUpateApiStateByOntIdApiId(tx, OntId, delState, uint32(apiId), time.Now().Unix())
	if errl != nil {
		log.Debugf("DelPulishApiCore.N.1. %s", errl)
		return errl
	}

	errl = tx.Commit()
	if errl != nil {
		log.Debugf("DelPulishApiCore.N.2. %s", errl)
		return errl
	}

	return nil
}

func (this *OscoreApi) QueryBasicApiInfoByCategory(id, pageNum, pageSize uint32) ([]*tables.ApiBasicInfo, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefOscoreApiDB.QueryApiBasicInfoByCategoryId(nil, id, start, pageSize)
}

func (this *OscoreApi) QueryApiDetailInfoByApiId(apiId uint32, apiState []int32) (*common.ApiDetailResponse, error) {
	basicInfo, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiIdInState(nil, apiId, apiState)
	if err != nil {
		return nil, err
	}

	requestParam, err := dao.DefOscoreApiDB.QueryRequestParamByApiId(nil, basicInfo.ApiId, tables.REQUEST_PARAM_TAG_PARAM)
	if err != nil {
		return nil, err
	}

	errorCode := make([]*common.PublishErrorCode, 0)
	if basicInfo.ErrorDesc != "" {
		err = json.Unmarshal([]byte(basicInfo.ErrorDesc), &errorCode)
		if err != nil {
			return nil, err
		}
	}

	sp, err := dao.DefOscoreApiDB.QuerySpecificationsByApiId(nil, basicInfo.ApiId)
	if err != nil {
		return nil, err
	}

	ResponseParamDescs := make([]*common.ResponseParamDesc, 0)
	if basicInfo.ResponseParam != "" {
		err = json.Unmarshal([]byte(basicInfo.ResponseParam), &ResponseParamDescs)
		if err != nil {
			return nil, err
		}
	}

	return &common.ApiDetailResponse{
		ApiId:               basicInfo.ApiId,
		Mark:                basicInfo.Mark,
		ResponseType:        basicInfo.RequestType,
		ResponseParam:       basicInfo.ResponseParam,
		ResponseExample:     basicInfo.ResponseExample,
		DataDesc:            basicInfo.DataDesc,
		DataSource:          basicInfo.DataSource,
		ApplicationScenario: basicInfo.ApplicationScenario,
		RequestParams:       requestParam,
		ErrorCodes:          errorCode,
		Specifications:      sp,
		ResponseParamDescs:  ResponseParamDescs,
		ApiBasicInfo:        basicInfo,
	}, nil
}

func (this *OscoreApi) SearchApiIdByCategoryId(categoryId, pageNum, pageSize uint32) ([]*common.ApiAttachMent, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize
	infos, err := dao.DefOscoreApiDB.QueryApiBasicInfoByCategoryId(nil, categoryId, start, pageSize)
	if err != nil {
		log.Debugf("SearchApiIdByCategoryId N.0 %s", err)
		return nil, err
	}

	return BuildApiAttatchMent(infos)
}

func (this *OscoreApi) SearchApi() (map[string][]*tables.ApiBasicInfo, error) {
	return dao.DefOscoreApiDB.SearchApi(nil)
}

func (this *OscoreApi) SearchFreeApi(pageNum, pageSize int) (map[string]interface{}, error) {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	start := (pageNum - 1) * pageSize
	return dao.DefOscoreApiDB.SearchFreeApi(nil, int(start), int(pageSize))
}

type PublishAPI struct {
	Name               string                     `json:"name"`
	Desc               string                     `json:"description"`
	RequestType        string                     `json:"requestType"`
	ApiProvider        string                     `json:"apiProvider"`
	DataSource         string                     `json:"dataSource"`
	ResponseExample    string                     `json:"responseExample"`
	Abstract           string                     `json:"abstract"`
	CryptoInfos        []tables.ApiCryptoInfo     `json:"cryptoInfos"`
	Tags               []tables.Tag               `json:"tags"`
	ErrorCodes         []common.PublishErrorCode  `json:"errorCodes"`
	Params             []tables.RequestParam      `json:"params"`
	ResponseParamDescs []common.ResponseParamDesc `json:"responseParamDesc"`
	Specs              []tables.Specifications    `json:"specifications"`
}

const (
	ACCESS_MUST_OK        int32 = 1
	ACCESS_BUT_IGNORE_ERR int32 = 2
	ACCESS_NO             int32 = 3
)

func PublishAPIHandleCore(param *PublishAPI, ontId, author string, accessMode int32, userId string) error {
	// handle error
	log.Debugf("PublishAPIHandleCore. start")
	errorDesc, err := json.Marshal(param.ErrorCodes)
	if err != nil {
		log.Debugf("PublishAPIHandleCore.0. %s", err)
		return err
	}

	if len(param.Tags) > 100 || len(param.Params) > 100 || len(param.Specs) > 100 || len(param.ErrorCodes) > 100 {
		log.Debugf("PublishAPIHandleCore.1. %s", err)
		return err
	}

	if len(param.Specs) == 0 {
		return errors.New("publish must fill a spec")
	}

	if len(param.Tags) == 0 {
		log.Debugf("PublishAPIHandleCore.1. tag can not empty")
		return errors.New("PublishAPIHandleCore.1. tag can not empty")
	}

	tags := make([]*tables.Tag, 0)

	for _, tag := range param.Tags {
		t, err := dao.DefOscoreApiDB.QueryTagByNameId(nil, tag.CategoryId, tag.Name)
		if err != nil {
			log.Debugf("PublishAPIHandleCore. Y.x.0 %v", tag)
			log.Debugf("PublishAPIHandleCore.2. %s", err)
			return err
		}
		tags = append(tags, t)
	}

	cat, err := dao.DefOscoreApiDB.QueryCategoryById(nil, tags[0].CategoryId)
	if err != nil {
		log.Debugf("PublishAPIHandleCore. Y.x.1 %v", cat)
		log.Debugf("PublishAPIHandleCore.2.0 %s", err)
		return err
	}

	if param.RequestType != "POST" && param.RequestType != "GET" {
		return errors.New("wrong RequestType type. only POST/GET")
	}

	rDescs, err := json.Marshal(param.ResponseParamDescs)
	if err != nil {
		log.Debugf("PublishAPIHandleCore.2.1 %s", err)
		return err
	}

	apibasic := &tables.ApiBasicInfo{
		Title:           param.Name,
		Icon:            cat.Icon,
		ApiProvider:     param.ApiProvider,
		ApiOscoreUrlKey: common.GenerateUUId(common.UUID_TYPE_OSCORE_URL),
		ApiDesc:         param.Desc,
		ApiState:        tables.API_STATE_PUBLISH,
		ErrorDesc:       string(errorDesc),
		RequestType:     param.RequestType,
		ResponseParam:   string(rDescs),
		ResponseExample: param.ResponseExample,
		DataSource:      param.DataSource,
		OntId:           ontId,
		UserId:          userId,
		ApiKind:         1,
		Author:          author,
		Price:           "0",
		Popularity:      100,
		Delay:           0,
		SuccessRate:     100,
		InvokeFrequency: 0,
		Abstract:        param.Abstract,
		CreateTime:      time.Now().Unix(),
	}
	port := fmt.Sprintf("%d", oscoreconfig.DefOscoreConfig.RestPort)
	apibasic.ApiUrl = oscoreconfig.DefOscoreConfig.OscoreHost + ":" + port + "/api/v1/data_source/" + apibasic.ApiOscoreUrlKey + "/:apikey"

	tx, errl := dao.DefOscoreApiDB.DB.Beginx()
	if errl != nil {
		log.Debugf("PublishAPIHandleCore.3. %s", err)
		return err
	}

	defer func() {
		if errl != nil {
			tx.Rollback()
		}
	}()

	log.Debugf("PublishAPIHandleCore. Y.0")
	err = dao.DefOscoreApiDB.InsertApiBasicInfo(tx, []*tables.ApiBasicInfo{apibasic})
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.4. %s", err)
		return err
	}

	info, err := dao.DefOscoreApiDB.QueryApiBasicInfoByOscoreUrlKey(tx, apibasic.ApiOscoreUrlKey, tables.API_STATE_PUBLISH)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.5. %s", err)
		return err
	}

	// cryptoinfo handle.
	if len(param.CryptoInfos) != 0 {
		for _, otherInfo := range param.CryptoInfos {
			otherInfo.ApiId = info.ApiId
			err = dao.DefOscoreApiDB.InsertApiCryptoInfo(tx, &otherInfo)
			if err != nil {
				errl = err
				log.Debugf("PublishAPIHandleCore 5.1 %s", err)
				return err
			}
		}
	}

	referCryptoInfos, err := dao.DefOscoreApiDB.QueryApiCryptoInfoByApiId(tx, info.ApiId)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.5.2. %s", err)
		return err
	}

	if len(referCryptoInfos) != len(param.CryptoInfos) {
		errl = fmt.Errorf("PublishAPIHandleCore 5.3 otherInfo insert err")
		log.Debugf("PublishAPIHandleCore 5.3 otherInfo insert err %d should be %d", len(referCryptoInfos), len(param.CryptoInfos))
		return errl
	}
	// tag handle

	for _, apiTag := range tags {
		tag := &tables.ApiTag{
			ApiId: info.ApiId,
			TagId: apiTag.Id,
			State: byte(1),
		}
		err = dao.DefOscoreApiDB.InsertApiTag(tx, tag)
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.6. %s", err)
			return err
		}
	}

	// handle param
	for _, p := range param.Params {
		p.ApiId = info.ApiId
		err := dao.DefOscoreApiDB.InsertRequestParam(tx, []*tables.RequestParam{&p})
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.7. %s", err)
			return err
		}
	}

	// spec handle.
	if len(param.Specs) == 0 {
		return errors.New("must fill a spec")
	}

	for _, s := range param.Specs {
		s.ApiId = info.ApiId
		err := dao.DefOscoreApiDB.InsertSpecifications(tx, []*tables.Specifications{&s})
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.8. %s", err)
			return err
		}
	}

	// min price handle.
	spec, err := dao.DefOscoreApiDB.GetMinPriceSpecOfApiId(tx, info.ApiId)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.8.1. %s", err)
		return err
	}

	err = dao.DefOscoreApiDB.UpdateApiBasicPrice(tx, info.ApiId, spec.Price)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.8.2. %s", err)
		return err
	}

	referParams, err := dao.DefOscoreApiDB.QueryRequestParamByApiId(tx, info.ApiId, tables.REQUEST_PARAM_TAG_PARAM)
	if err != nil {
		errl = err
		log.Debugf("PublishAPIHandleCore.9. %s", err)
		return err
	}

	var responseExample []byte
	var errinner error
	switch accessMode {
	case ACCESS_BUT_IGNORE_ERR:
		responseExample, errinner = HandleDataSourceReqCore(tx, info.ApiOscoreUrlKey, referParams, "", true)
		if errinner != nil {
			log.Debugf("PublishAPIHandleCore.10. %s", errinner)
		}
	case ACCESS_MUST_OK:
		responseExample, err = HandleDataSourceReqCore(tx, info.ApiOscoreUrlKey, referParams, "", true)
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.10. %s", err)
			return err
		}
	case ACCESS_NO:
	default:
		log.Debugf("PublishAPIHandleCore.10.0. %s", "error access mode")
		return errors.New("error access mode.")
	}

	if info.ResponseExample == "" {
		err = dao.DefOscoreApiDB.UpdateApiBasicResponseExample(tx, info.ApiId, string(responseExample))
		if err != nil {
			errl = err
			log.Debugf("PublishAPIHandleCore.10.1. %s", err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Debugf("PublishAPIHandleCore.11. %s", err)
		errl = err
		return err
	}

	return nil
}

func BuildApiAttatchMent(infos []*tables.ApiBasicInfo) ([]*common.ApiAttachMent, error) {
	res := make([]*common.ApiAttachMent, 0)
	for _, api := range infos {
		sps, err := dao.DefOscoreApiDB.GetMinPriceSpecOfApiId(nil, api.ApiId)
		if err != nil {
			return nil, err
		}
		t := &common.ApiAttachMent{
			api,
			sps.Id,
			sps.Price,
			sps.Amount,
			sps.EffectiveDuration,
			sps.SpecType,
		}
		res = append(res, t)
	}

	return res, nil
}

func SearchApiByKey(key string) ([]*common.ApiAttachMent, error) {
	infos, err := dao.DefOscoreApiDB.SearchApiByKey(key, tables.API_STATE_BUILTIN)
	if err != nil {
		return nil, err
	}

	return BuildApiAttatchMent(infos)
}
