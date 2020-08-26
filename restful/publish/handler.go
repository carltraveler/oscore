package publish

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/oscore/common"
	"github.com/ontio/oscore/core"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/oscoreconfig"
	"github.com/ontio/oscore/restful/api/common"
	"strconv"
	"strings"
)

type UrlParams struct {
	Name  string
	Type  int32
	Index uint32
}

func GetSellerMarketApi(c *gin.Context) {
	arr, err := common.ParseGetParamByParamName(c, "ontId", "pageNum", "pageSize")
	if len(arr) != 3 {
		log.Errorf("[GetSellerMarketApi] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("param should be 3 num.")))
		return
	}
	OntId := arr[0]
	pageNum, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Errorf("[GetSellerMarketApi] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(arr[2])
	if err != nil {
		log.Errorf("[GetSellerMarketApi] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	log.Debugf("GetSellerMarketApi: %d, %d, %s", pageNum, pageSize, OntId)

	infos, err := dao.DefOscoreApiDB.QueryApiBasicInfoByOntId(nil, OntId, []int32{tables.API_STATE_BUILTIN}, uint32(pageNum), uint32(pageSize))
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	respinfos, err := core.BuildApiAttatchMent(infos)
	if err != nil {
		log.Debugf("GetSellerMarketApi: N.0 %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	count, err := dao.DefOscoreApiDB.QueryApiBasicInfoOntIdCount(nil, OntId, []int32{tables.API_STATE_BUILTIN})
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	var scoreNum float64
	log.Debugf("GetSellerMarketApi: Y.0")
	score, err := dao.DefOscoreApiDB.QueryOntScoreByOntId(nil, OntId)
	log.Debugf("GetSellerMarketApi: Y.1 : %s", score)
	if err != nil && !dao.IsErrNoRows(err) {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	} else if dao.IsErrNoRows(err) {
		scoreNum = 0
	} else {
		scoreNum = float64(score.TotalScore) / float64(score.TotalCommentNum)
	}

	value, err := strconv.ParseFloat(fmt.Sprintf("%.1f", scoreNum), 64)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	res := map[string]interface{}{
		"count": count,
		"list":  respinfos,
		"score": value,
	}

	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetPulishApi(c *gin.Context) {
	arr, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	pageNum, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Errorf("[GetALLPublishPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Errorf("[GetALLPublishPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no ontid")))
		return
	}
	OntId := ontid.(string)
	log.Debugf("GetPulishApi: %d, %d, %s", pageNum, pageSize, OntId)

	infos, err := dao.DefOscoreApiDB.QueryApiBasicInfoByOntId(nil, OntId, []int32{tables.API_STATE_PUBLISH, tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER, tables.API_STATE_DELETE}, uint32(pageNum), uint32(pageSize))
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	count, err := dao.DefOscoreApiDB.QueryApiBasicInfoOntIdCount(nil, OntId, []int32{tables.API_STATE_PUBLISH, tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER, tables.API_STATE_DELETE})
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	res := common.PageResult{
		List:  infos,
		Count: count,
	}

	common.WriteResponse(c, common.ResponseSuccess(res))
}

func PublishAPIHandle(c *gin.Context) {
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no ontid")))
		return
	}
	OntId := ontid.(string)

	userId, ok := c.Get(oscoreconfig.Key_UserId)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no userId")))
		return
	}

	author, ok := c.Get(oscoreconfig.JWTAud)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no author")))
		return
	}
	Author := author.(string)

	param := &core.PublishAPI{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	err = core.PublishAPIHandleCore(param, OntId, Author, core.ACCESS_MUST_OK, userId.(string))
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func VerifyAPIHandle(c *gin.Context) {
	res, err := common.ParseGetParamByParamName(c, "apiId", "oscoreUrlKey")
	if err != nil {
		log.Errorf("[VerifyAPIHandle] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if len(res) != 2 {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("need pass apiId and oscoreUrlKey.")))
		return
	}
	apiId, err := strconv.Atoi(res[0])
	if err != nil {
		log.Errorf("[VerifyAPIHandle] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	oscoreUrlKey := res[1]
	if oscoreUrlKey == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("oscoreUrlKey can not empty.")))
		return
	}
	err = dao.DefOscoreApiDB.ApiBasicUpateApiState(nil, tables.API_STATE_BUILTIN, uint32(apiId), oscoreUrlKey)
	if err != nil {
		log.Errorf("[VerifyAPIHandle] ApiBasicUpateApiState error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func DelPulishApi(c *gin.Context) {
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no ontid")))
		return
	}
	OntId := ontid.(string)

	res, err := common.ParseGetParamByParamName(c, "apiId")
	if err != nil {
		log.Errorf("[DelPulishApi] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if len(res) != 1 {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("need pass apiId.")))
		return
	}

	apiId, err := strconv.Atoi(res[0])
	if err != nil {
		log.Errorf("[DelPulishApi] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	err = core.DefOscoreApi.DelPulishApiCore(OntId, uint32(apiId))
	if err != nil {
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func GetApiDetailByApiIdApiState(c *gin.Context) {
	res, err := common.ParseGetParamByParamName(c, "apiId", "apiState")
	if err != nil {
		log.Errorf("[GetApiDetailByApiIdApiState] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	apiId, err := strconv.Atoi(res[0])
	if err != nil {
		log.Errorf("[GetApiDetailByApiIdApiState] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	apiState, err := strconv.Atoi(res[1])
	if err != nil {
		log.Errorf("[GetApiDetailByApiIdApiState] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	info, err := core.DefOscoreApi.QueryApiDetailInfoByApiId(uint32(apiId), []int32{int32(apiState)})
	if err != nil {
		log.Errorf("[GetApiDetailByApiId] QueryApiDetailInfoByApiId error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(info))
}

func AdminTestAPIKey(c *gin.Context) {
	var params []*tables.RequestParam

	apiId := c.Param("apiId")
	if apiId == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("apiId can not empty.")))
		return
	}

	id, err := strconv.Atoi(apiId)
	if err != nil {
		log.Errorf("[AdminTestAPIKey] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	err = common.ParsePostParam(c, &params)
	if err != nil {
		log.Errorf("[AdminTestAPIKey] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	data, err := core.DefOscoreApi.AdminTestApi(params, uint32(id))
	if err != nil {
		log.Errorf("[AdminTestAPIKey] failed: %s", err.Error())
		res := make(map[string]string)
		res["errorDesc"] = err.Error()
		bs, _ := json.Marshal(res)
		common.WriteResponse(c, common.ResponseSuccess(string(bs)))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(string(data)))
}

func GetALLPublishPage(c *gin.Context) {
	arr, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	pageNum, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Errorf("[GetALLPublishPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(arr[1])
	if err != nil {
		log.Errorf("[GetALLPublishPage] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	log.Debugf("pageNum:%d, pageSize %d", pageNum, pageSize)
	infos, err := dao.DefOscoreApiDB.QueryApiBasicInfoByPage(uint32(pageNum), uint32(pageSize), tables.API_STATE_PUBLISH)
	if err != nil {
		log.Errorf("[GetALLPublishPage] QueryApiBasicInfoByPage error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	count, err := dao.DefOscoreApiDB.QueryApiBasicInfoCount(nil, tables.API_STATE_PUBLISH)
	if err != nil {
		log.Errorf("[GetALLPublishPage] QueryApiBasicInfoByPage error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	res := common.PageResult{
		List:  infos,
		Count: count,
	}

	common.WriteResponse(c, common.ResponseSuccess(res))
}

func AdminGenerateTestKey(c *gin.Context) {
	params := &common2.AdminGenerateTestKeyParam{}
	err := common.ParsePostParam(c, params)
	if err != nil {
		log.Errorf("[AdminGenerateTestKey] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[AdminGenerateTestKey] ontId is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	testKey, err := core.DefOscoreApi.GenerateApiTestKey(params.ApiId, ontId.(string), tables.API_STATE_PUBLISH)
	if err != nil || testKey == nil {
		log.Errorf("[AdminGenerateTestKey] GenerateApiTestKey failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(testKey))
}

func ParseUrl(url string) ([]UrlParams, error) {
	params := make([]UrlParams, 0)
	key := url
	count := uint32(0)
	var i, j, k int
	var queryArgHandled bool
	for {
		i = strings.IndexAny(key, "{")
		if i != -1 {
			if queryArgHandled {
				return nil, errors.New("url errror")
			}
			key = key[i:]
			i = strings.IndexAny(key, "{")
			k = strings.IndexAny(key, "}")
			j = strings.IndexAny(key, "/")
			if k == -1 || (j != -1 && k+1 != j) {
				return nil, errors.New("url error")
			}
			p := key[i+1 : k]
			if i+1 == k {
				return nil, errors.New("url error")
			}
			params = append(params, UrlParams{
				Name: p,
				Type: tables.URL_PARAM_RESTFUL,
			})
			count += 1
			key = key[k:]
		} else {
			if !queryArgHandled {
				i = strings.IndexAny(key, "?")
				if i == -1 {
					break
				}

				k = strings.IndexAny(key, "&")
				if k == -1 {
					k = len(key)
				}

				p := key[i+1 : k]
				if i+1 == k {
					return nil, errors.New("url error")
				}
				params = append(params, UrlParams{
					Name: p,
					Type: tables.URL_PARAM_QUERY,
				})
				count += 1
				key = key[k:]
			}
			queryArgHandled = true

			i = strings.IndexAny(key, "&")
			if i == -1 {
				break
			}

			k = strings.IndexAny(key[i+1:], "&")
			if k == -1 {
				k = len(key)
			} else {
				k += 1
			}

			if i+1 == k {
				return nil, errors.New("url error")
			}
			p := key[i+1 : k]
			params = append(params, UrlParams{
				Name: p,
				Type: tables.URL_PARAM_QUERY,
			})
			count += 1
			key = key[k:]
		}
	}

	return params, nil
}
