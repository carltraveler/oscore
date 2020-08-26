package order

import (
	"github.com/gin-gonic/gin"
	"github.com/ontio/oscore/middleware/jwt"
	"github.com/ontio/oscore/aksk"
)

func RoutesOrder(parent *gin.RouterGroup) {
	orderRouteGroup := parent.Group("/order")
	orderRouteGroup.GET("/GetCommentsByApiId/:pageNum/:pageSize/:apiId", GetCommentsByApiId)
	orderRouteGroup.GET("/GetCommentsByToolBoxId/:pageNum/:pageSize/:toolBoxId", GetCommentsByToolBoxId)
	orderRouteGroup.Use(jwt.JWT())
	orderRouteGroup.POST("/takeOrder", TakeOrder)
	orderRouteGroup.POST("/renewOrder", RenewOrder)
	orderRouteGroup.POST("/queryAliPayResult", QueryAliPayResultResetful)
	orderRouteGroup.POST("/aliPayOder", AliPayOder)
	orderRouteGroup.POST("/takeapiprocessOrder", TakeWetherForcastApiOrder)
	orderRouteGroup.POST("/cancelOrder", CancelOrder)
	orderRouteGroup.POST("/generateTestKey", GenerateTestKey)
	orderRouteGroup.POST("/testAPIKey/:apiKey", TestAPIKey)
	orderRouteGroup.GET("/queryOrderStatus/:orderId", GetTxResult)
	orderRouteGroup.GET("/queryOrderByPage/:pageNum/:pageSize", QueryOrderByPage)
	orderRouteGroup.GET("/queryApiKeysByPage/:pageNum/:pageSize", QueryApiKeysByPage)
	orderRouteGroup.GET("/QueryDataProcessOrderByPage/:pageNum/:pageSize", QueryDataProcessOrderByPage)
	orderRouteGroup.GET("/QueryDataProcessResultByPage/:pageNum/:pageSize", QueryDataProcessResultByPage)
	orderRouteGroup.GET("/GetOrderDetailById/:orderId", GetOrderDetailById)
	orderRouteGroup.POST("/DelCommentOrderByID", DelCommentOrderByID)
	orderRouteGroup.POST("/CommentOrderApi", CommentOrderApi)
	orderRouteGroup.POST("/GetApiOrderList", GetApiOrderList)

	orderRouteGroupNoToken := parent.Group("/ali")
	orderRouteGroupNoToken.Use(aksk.AkSk)
	orderRouteGroupNoToken.POST("/sendTxAli", SendTxAli)

	orderRouteGroupAdmin := parent.Group("/admin")
	orderRouteGroupAdmin.Use(jwt.JWTAdmin())
	orderRouteGroupAdmin.GET("/GetCommentsPage/:pageNum/:pageSize", GetCommentsPage)
	orderRouteGroupAdmin.GET("/GetCommentsById/:commentId", GetCommentsById)
	orderRouteGroupAdmin.GET("/DelCommentsById/:commentId", DelCommentsById)
}
