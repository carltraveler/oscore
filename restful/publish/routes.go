package publish

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/oscore/middleware/jwt"
)

func RoutesPublish(parent *gin.RouterGroup) {
	publishG := parent.Group("/publish")
	publishG.GET("/getSellerMarketApi/:ontId/:pageNum/:pageSize", GetSellerMarketApi)
	publishG.Use(jwt.JWT())
	publishG.POST("/api", PublishAPIHandle)
	publishG.POST("/delpublishapi/:apiId", DelPulishApi)
	publishG.GET("/getpublishapi/:pageNum/:pageSize", GetPulishApi)

	publishGAmin := publishG.Group("/admin")
	publishGAmin.Use(jwt.JWTAdmin())
	publishGAmin.GET("/getallpublishapi/:pageNum/:pageSize", GetALLPublishPage)
	publishGAmin.GET("/getapidetailinfo/:apiId/:apiState", GetApiDetailByApiIdApiState)

	publishGAmin.POST("/admintest/:apiId", AdminTestAPIKey)
	publishGAmin.POST("/publish/:apiId/:oscoreUrlKey", VerifyAPIHandle)
}
