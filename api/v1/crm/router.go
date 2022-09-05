package crm

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/leads", NewLead)
	g.GET("/leads", GetLeadList)
	g.PUT("/leads/:id", UpdateLead)
	g.GET("/leads/:id", GetLeadByID)
	g.DELETE("/leads/:id", DeleteLead)
	g.POST("/leads/:id/convert", ConvertLead)

}
