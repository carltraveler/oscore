package process

import (
	"github.com/gin-gonic/gin"
)

func RoutesDataProcess(parent *gin.RouterGroup) {
	publishG := parent.Group("/dataprocess")
	//publishG.Use(jwt.JWT())
	publishG.GET("/GetWetherForcastInfo/:toolid", GetWetherForcastInfo)
	publishG.GET("/GetLocation/:country", GetLocation)
	publishG.GET("/GetAllToolBox", GetAllToolBox)

	publishG.POST("/searchToolBoxByKey", searchToolBoxByKey)
	publishG.POST("/searchToolBoxByCategory", searchToolBoxByCategory)
}
