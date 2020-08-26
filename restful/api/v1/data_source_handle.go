package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/oscore/common"
	"github.com/ontio/oscore/core"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	"io/ioutil"
	"net/http"
)

func HandleDataSourceReq(c *gin.Context) {
	log.Debugf("HandleDataSourceReq")
	oscoreUrlKey := c.Param("oscoreUrlKey")
	if !common2.IsOscoreUrlKey(oscoreUrlKey) {
		c.String(http.StatusInternalServerError, "oscore url error.")
		return
	}
	apiKey := c.Param("apiKey")
	if !common2.IsApiKey(apiKey) && !common2.IsTestKey(apiKey) {
		c.String(http.StatusInternalServerError, "apikey false.")
		return
	}

	info, err := dao.DefOscoreApiDB.QueryApiBasicInfoByOscoreUrlKeyStates(nil, oscoreUrlKey, []int32{tables.API_STATE_BUILTIN, tables.API_STATE_DISABLE_ORDER})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	referParams, err := dao.DefOscoreApiDB.QueryRequestParamByApiId(nil, info.ApiId, tables.REQUEST_PARAM_TAG_PARAM)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	bodymapsTmps := make(map[string]interface{})
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	err = json.Unmarshal(paramsBs, &bodymapsTmps)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	for _, p := range referParams {
		if _, ok := bodymapsTmps[p.ParamName]; ok {
			if p.ParamWhere == tables.URL_PARAM_BODY && p.ParamTag == tables.REQUEST_PARAM_TAG_PARAM {
				vstr, err := json.Marshal(bodymapsTmps[p.ParamName])
				p.ValueDesc = string(vstr)
				if err != nil {
					log.Debugf("HandleDataSourceReq %s marshal error", p.ParamName)
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
			} else {
				vstr, ok := bodymapsTmps[p.ParamName].(string)
				if !ok {
					log.Debugf("HandleDataSourceReq %s value must be string", p.ParamName)
					c.String(http.StatusInternalServerError, err.Error())
					return
				}
				p.ValueDesc = string(vstr)
				log.Debugf("paramName %s : %s", p.ParamName, p.ValueDesc)
			}
		} else {
			p.ValueDesc = ""
		}
	}

	data, err := core.HandleDataSourceReqCore(nil, oscoreUrlKey, referParams, apiKey, false)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, string(data))
}
