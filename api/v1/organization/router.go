package organization

import "github.com/gin-gonic/gin"

func Routers(g *gin.RouterGroup) {
	g.POST("/organizations", NewOrganization)
}
