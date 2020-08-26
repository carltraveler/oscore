package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ontio/ontology/common/log"
	common2 "github.com/ontio/oscore/common"
	"github.com/ontio/oscore/core"
	"github.com/ontio/oscore/models/tables"
	"github.com/ontio/oscore/restful/api/common"
	"github.com/ontio/oscore/oscoreconfig"
	"io/ioutil"
	"net/http"
	"strconv"
)

func TakeWetherForcastApiOrder(c *gin.Context) {
	param := &common2.WetherForcastRequest{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[TakeWetherForcastApiOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[TakeWetherForcastApiOrder] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	OntId := ontid.(string)
	userId, ok := c.Get(oscoreconfig.Key_UserId)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no userId")))
		return
	}

	res, err := core.DefOscoreApi.OscoreOrder.TakeWetherForcastApiOrder(param, OntId, userId.(string))
	if err != nil {
		log.Errorf("[TakeWetherForcastApiOrder] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func TakeOrder(c *gin.Context) {
	param := &common2.TakeOrderParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[TakeOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[TakeOrder] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	param.OntId = ontid.(string)
	userId, ok := c.Get(oscoreconfig.Key_UserId)
	if !ok {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("no userId")))
		return
	}

	res, err := core.DefOscoreApi.OscoreOrder.TakeOrder(param, userId.(string), tables.TAKE_ORDER_DEFAULT)
	if err != nil {
		log.Errorf("[TakeOrder] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func QueryAliPayResultResetful(c *gin.Context) {
	param := &common2.AliPayOderParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[QueryAliPayResultResetfu] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[QueryAliPayResultResetfu] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}

	resp, err := core.DefOscoreApi.OscoreOrder.QueryAliPayResult(param.OrderId, ontid.(string))
	if err != nil {
		log.Errorf("[QueryAliPayResultResetful] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(resp))
}

func AliPayOder(c *gin.Context) {
	param := &common2.AliPayOderParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[RenewOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[RenewOrder] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}

	resp, err := core.DefOscoreApi.OscoreOrder.RequestAliPay(param.OrderId, ontid.(string))

	if err != nil {
		log.Errorf("[RenewOrder] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(resp))
}

func RenewOrder(c *gin.Context) {
	param := &common2.RenewOrderParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[RenewOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[RenewOrder] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	resp, err := core.DefOscoreApi.OscoreOrder.RenewOrder(param, ontid.(string))

	if err != nil {
		log.Errorf("[RenewOrder] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}

	common.WriteResponse(c, common.ResponseSuccess(resp))
}

func QueryOrderByPage(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[QueryOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok || ontId == "" {
		log.Errorf("[QueryOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	log.Infof("[QueryOrderByPage] ontid:%s", ontId)
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	paseSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	orders, err := core.DefOscoreApi.OscoreOrder.QueryOrderByPage(pageNum, paseSize, ontId.(string))
	if err != nil {
		log.Errorf("[QueryOrderByPage] QueryOrderByPage failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(orders))
}

func QueryApiKeysByPage(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[QueryApiKeysByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok || ontId == "" {
		log.Errorf("[QueryApiKeysByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	log.Infof("[QueryOrderByPage] ontid:%s", ontId)
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryApiKeysByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	paseSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryApiKeysByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	res, err := core.DefOscoreApi.OscoreOrder.QueryApiKeyByPage(pageNum, paseSize, ontId.(string))
	if err != nil {
		log.Errorf("[QueryApiKeysByPage] QueryOrderByPage failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func QueryDataProcessOrderByPage(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok || ontId == "" {
		log.Errorf("[QueryDataProcessOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	log.Infof("[QueryDataProcessOrderByPage] ontid:%s", ontId)
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	paseSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	orders, err := core.DefOscoreApi.OscoreOrder.QueryDataProcessOrderByPage(pageNum, paseSize, ontId.(string))
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] QueryDataProcessOrderByPage failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(orders))
}

func QueryDataProcessResultByPage(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok || ontId == "" {
		log.Errorf("[QueryDataProcessOrderByPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	log.Infof("[QueryDataProcessOrderByPage] ontid:%s", ontId)
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	paseSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	orders, err := core.DefOscoreApi.OscoreOrder.QueryDataProcessResultByPage(pageNum, paseSize, ontId.(string))
	if err != nil {
		log.Errorf("[QueryDataProcessOrderByPage] QueryDataProcessOrderByPage failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(orders))
}

func GetApiOrderList(c *gin.Context) {
	param := &common2.OrderListRequest{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[GetOrderList] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[GetOrderList] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	res, err := core.DefOscoreApi.OscoreOrder.GetApiOrderList(param, ontid.(string))
	if err != nil {
		log.Errorf("[GetOrderList] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}

func CommentOrderApi(c *gin.Context) {
	param := &common2.CommentOrderApiRequest{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[CommentOrderApi] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[CommentOrderApi] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	err = core.DefOscoreApi.OscoreOrder.CommentOrderApi(param, ontid.(string))
	if err != nil {
		log.Errorf("[CommentOrderApi] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func DelCommentOrderByID(c *gin.Context) {
	param := &common2.DelCommentOrderRequest{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[CommentOrderApi] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontid, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[CommentOrderApi] ontid is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("ontid is nil")))
		return
	}
	err = core.DefOscoreApi.OscoreOrder.DelCommentOrderByID(param, ontid.(string))
	if err != nil {
		log.Errorf("[CommentOrderApi] TakeOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func GetCommentsByApiId(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize", "apiId")
	if err != nil {
		log.Errorf("[GetCommentsByApiId] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	apiId, err := strconv.Atoi(params[2])
	if err != nil {
		log.Errorf("[GetCommentsByApiId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	comment, err := core.DefOscoreApi.OscoreOrder.GetCommentsByApiId(pageNum, pageSize, uint32(apiId))
	if err != nil {
		log.Errorf("[GetCommentsByApiId] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(comment))
}

func GetCommentsPage(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize")
	if err != nil {
		log.Errorf("[GetCommentsPage] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[GetCommentsPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[GetCommentsPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	comment, err := core.DefOscoreApi.OscoreOrder.GetCommentPage(pageNum, pageSize)
	if err != nil {
		log.Errorf("[GetCommentsPage] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(comment))
}

func GetCommentsById(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "commentId")
	if err != nil {
		log.Errorf("[GetCommentsById] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	commentId, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[GetCommentsById] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	comment, err := core.DefOscoreApi.OscoreOrder.GetCommentsById(uint32(commentId))
	if err != nil {
		log.Errorf("[GetCommentsById] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(comment))
}

func DelCommentsById(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "commentId")
	if err != nil {
		log.Errorf("[DelCommentsById] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	commentId, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[DelCommentsById] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	err = core.DefOscoreApi.OscoreOrder.DelCommentsById(uint32(commentId))
	if err != nil {
		log.Errorf("[DelCommentsById] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func GetCommentsByToolBoxId(c *gin.Context) {
	params, err := common.ParseGetParamByParamName(c, "pageNum", "pageSize", "toolBoxId")
	if err != nil {
		log.Errorf("[GetCommentsByApiId] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageNum, err := strconv.Atoi(params[0])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	pageSize, err := strconv.Atoi(params[1])
	if err != nil {
		log.Errorf("[QueryOrderByPage] Atoi failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	toolBoxId, err := strconv.Atoi(params[2])
	if err != nil {
		log.Errorf("[GetCommentsByApiId] ParseGetParam error: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	comment, err := core.DefOscoreApi.OscoreOrder.GetCommentsByToolBoxId(pageNum, pageSize, uint32(toolBoxId))
	if err != nil {
		log.Errorf("[GetCommentsByApiId] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(comment))
}

func GetOrderDetailById(c *gin.Context) {
	orderId := c.Param("orderId")
	if orderId == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("orderId can no empty")))
		return

	}
	orders, err := core.DefOscoreApi.OscoreOrder.GetOrderDetailById(orderId)
	if err != nil {
		log.Errorf("[GetOrderDetailById] failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(orders))
}

func GenerateTestKey(c *gin.Context) {
	params := &common2.GenerateTestKeyParam{}
	err := common.ParsePostParam(c, params)
	if err != nil {
		log.Errorf("[GenerateTestKey] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	ontId, ok := c.Get(oscoreconfig.Key_OntId)
	if !ok {
		log.Errorf("[GenerateTestKey] ontId is nil: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	testKey, err := core.DefOscoreApi.GenerateApiTestKey(params.ApiId, ontId.(string), tables.API_STATE_BUILTIN)
	if err != nil || testKey == nil {
		log.Errorf("[GenerateTestKey] GenerateApiTestKey failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(testKey))
}

func TestAPIKey(c *gin.Context) {
	var params []*tables.RequestParam
	err := common.ParsePostParam(c, &params)
	if err != nil {
		log.Errorf("[GenerateTestKey] ParseGetParamByParamName failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}

	apiKey := c.Param("apiKey")
	if apiKey == "" {
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, errors.New("apikey is nil")))
		return
	}

	data, err := core.DefOscoreApi.TestApiKey(params, apiKey)
	if err != nil {
		log.Errorf("[TestAPIKey] TestApiKey failed: %s", err.Error())
		res := make(map[string]string)
		res["errorDesc"] = err.Error()
		bs, _ := json.Marshal(res)
		common.WriteResponse(c, common.ResponseSuccess(string(bs)))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(string(data)))
}

func CancelOrder(c *gin.Context) {
	param := &common2.OrderIdParam{}
	err := common.ParsePostParam(c, param)
	if err != nil {
		log.Errorf("[CancelOrder] ParsePostParam failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, err))
		return
	}
	if param == nil || param.OrderId == "" {
		log.Errorf("[CancelOrder] param is nil failed")
		common.WriteResponse(c, common.ResponseFailed(common.PARA_ERROR, fmt.Errorf("param is nil")))
		return
	}
	err = core.DefOscoreApi.OscoreOrder.CancelOrder(param.OrderId)
	if err != nil {
		log.Errorf("[CancelOrder] CancelOrder failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	log.Infof("[CancelOrder] orderId:%s", param.OrderId)
	common.WriteResponse(c, common.ResponseSuccess(nil))
}

func SendTxAli(c *gin.Context) {
	log.Debugf("SendTxAli CallBack start.")
	param := &core.AliPayCallBackArg{}
	paramsBs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Errorf("[SendTxAli] ParsePostParam failed: %s", err)
		writeOntoResponse(c, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}

	err = json.Unmarshal(paramsBs, param)
	if err != nil {
		log.Errorf("[SendTxAli] ParsePostParam failed: %s", err)
		writeOntoResponse(c, common.ResponseFailedOnto(common.PARA_ERROR, err))
		return
	}
	err = core.SendTxAliCore(param)
	if err != nil {
		log.Errorf("[SendTxAli] failed: %s", err)
		writeOntoResponse(c, common.ResponseFailedOnto(common.INTER_ERROR, err))
		return
	}
	writeOntoResponse(c, common.ResponseSuccessOnto())
}

func writeOntoResponse(c *gin.Context, param map[string]interface{}) {
	c.JSON(http.StatusOK, param)
}

func GetTxResult(c *gin.Context) {
	orderId := c.Param("orderId")
	res, err := core.DefOscoreApi.OscoreOrder.GetTxResult(orderId)
	if err != nil {
		log.Errorf("[GetTxResult] QueryOrderByOrderId failed: %s", err)
		common.WriteResponse(c, common.ResponseFailed(common.INTER_ERROR, err))
		return
	}
	common.WriteResponse(c, common.ResponseSuccess(res))
}
