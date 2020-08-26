package api

import "github.com/gin-gonic/gin"

func RoutesApiList(parent *gin.RouterGroup) {
	apiRouteGroup := parent.Group("/apilist")

	apiRouteGroup.POST("/searchApiByKey", SearchApiByKey)
	apiRouteGroup.POST("/searchApiByCategory", SearchApiByCategoryId)

	apiRouteGroup.GET("/getBasicApiInfoByPage/:pageNum/:pageSize", GetBasicApiInfoByPage)
	apiRouteGroup.GET("/getApiDetailByApiId/:apiId", GetApiDetailByApiId)
	apiRouteGroup.GET("/searchApi", SearchApi) //get newest
	apiRouteGroup.GET("/searchFreeApi/:pageNum/:pageSize", SearchFreeApi) //get newest
}
