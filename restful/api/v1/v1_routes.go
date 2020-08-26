package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/oscore/middleware/cors"
	"github.com/ontio/oscore/restful/api/v1/api"
	"github.com/ontio/oscore/restful/api/v1/order"
	"github.com/ontio/oscore/restful/api/v1/process"
)

func RoutesV1(parent *gin.RouterGroup) {
	v1Route := parent.Group("/v1")
	v1Route.Use(cors.Cors())
	order.RoutesOrder(v1Route)
	api.RoutesApiList(v1Route)
	process.RoutesDataProcess(v1Route)

	v1Route.POST("/data_source/:oscoreUrlKey/:apiKey", HandleDataSourceReq)
}
