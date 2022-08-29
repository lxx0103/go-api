package common

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.GET("/historys", GetHistoryList)
	g.GET("/nextnumber", GetNextNumber)
}
