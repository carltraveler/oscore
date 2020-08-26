package process

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/oscore/common"
	"github.com/ontio/oscore/dao"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/restful/api/common"
	"strconv"
)

type AlgorithmObj struct {
	Algorithm *tables.Algorithm `json:"algorithm"`
	Env       []*tables.Env     `json:"env"`
}

type ApiSourceObj struct {
	ApiSource  *tables.ApiBasicInfo `json:"apiSource"`
	Algorithms []*AlgorithmObj      `json:"algorithm"`
}

type WetherForcastResponse struct {
	ToolBox      *tables.ToolBox
	ApiSourceObj []*ApiSourceObj        `json:"apiSourceObj"`
	ApiALL       []*tables.ApiBasicInfo `json:"apiALL"`
	AlgorithmALL []*tables.Algorithm    `json:"algorithmALL"`
	EnvAll       []*tables.Env          `json:"envAll"`
}

func GetLocation(c *gin.Context) {
	country := c.Param("country")

	res, err := dao.DefOscoreApiDB.QueryLocationOfCountryCity(nil, country)
	if err != nil {
		log.Errorf("[GetLocation]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetAllToolBox(c *gin.Context) {
	res, err := dao.DefOscoreApiDB.QueryToolBoxAll(nil)
	if err != nil {
		log.Errorf("[GetAllToolBox]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func GetWetherForcastInfo(c *gin.Context) {
	toolid := c.Param("toolid")
	if toolid == "" {
		log.Errorf("[GetWetherForcastInfo]: %s", errors.New("toolid can not empty."))
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("toolid can not empty.")))
		return
	}
	toolboxid, err := strconv.Atoi(toolid)
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	toolbox, err := dao.DefOscoreApiDB.QueryToolBoxById(nil, uint32(toolboxid))
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	algAll := make([]*tables.Algorithm, 0)
	envAll := make([]*tables.Env, 0)

	apis, err := dao.DefOscoreApiDB.QueryApiBasicInfoByApiTypeKind(nil, toolbox.Title, tables.API_KIND_DATA_PROCESS, tables.API_STATE_BUILTIN)
	if err != nil {
		log.Errorf("[GetWetherForcastInfo]: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	apisourceresp := make([]*ApiSourceObj, 0)
	for _, api := range apis {
		apialgs, err := dao.DefOscoreApiDB.QueryApiAlgorithmsByApiId(nil, api.ApiId)
		if err != nil {
			log.Errorf("[GetWetherForcastInfo]: %s", err)
			common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
			return
		}
		algsresp := make([]*AlgorithmObj, 0)
		for _, apialg := range apialgs {
			alg, err := dao.DefOscoreApiDB.QueryAlgorithmById(nil, apialg.AlgorithmId)
			if err != nil {
				log.Errorf("[GetWetherForcastInfo]: %s", err)
				common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
				return
			}

			algenvs, err := dao.DefOscoreApiDB.QueryAlgorithmEnvByAlgorithmId(nil, alg.Id)
			if err != nil {
				log.Errorf("[GetWetherForcastInfo]: %s", err)
				common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
				return
			}
			envsresp := make([]*tables.Env, 0)
			for _, env := range algenvs {
				env, err := dao.DefOscoreApiDB.QueryEnvById(nil, env.Id)
				if err != nil {
					log.Errorf("[GetWetherForcastInfo]: %s", err)
					common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
					return
				}
				envsresp = append(envsresp, env)
				envAll = append(envAll, env)
			}
			algsresp = append(algsresp, &AlgorithmObj{
				Algorithm: alg,
				Env:       envsresp,
			})
			algAll = append(algAll, alg)
		}
		apisourceresp = append(apisourceresp, &ApiSourceObj{
			ApiSource:  api,
			Algorithms: algsresp,
		})
	}

	res := WetherForcastResponse{
		ToolBox:      toolbox,
		ApiSourceObj: apisourceresp,
		ApiALL:       apis,
		EnvAll:       envAll,
		AlgorithmALL: algAll,
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func searchToolBoxByKey(c *gin.Context) {
	key := &common2.SearchApiByKey{}
	err := common.ParsePostParam(c, key)
	if err != nil {
		log.Errorf("[searchToolBoxByKey] ParsePostParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if key == nil || key.Key == "" {
		common.WriteResponse(c, common.ResponseSuccess(nil))
		return
	}
	//todo key.Key should not have sql statement
	infos, err := dao.DefOscoreApiDB.SearchToolBoxByKey(nil, key.Key)
	if err != nil {
		log.Errorf("[searchToolBoxByKey] error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(infos))
}

func searchToolBoxByCategory(c *gin.Context) {
	param := &common2.GetApiByCategoryId{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[SearchApiByCategoryId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	infos, err := dao.DefOscoreApiDB.QueryToolBoxByCategoryId(nil, param.CategoryId, param.PageNum, param.PageSize)
	if err != nil {
		log.Errorf("[SearchApiByCategoryId] SearchApiByKey error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	log.Debugf("SearchApiByCategoryId: num %d", len(infos))
	common.WriteResponse(c, common.ResponseSuccess(infos))
}
