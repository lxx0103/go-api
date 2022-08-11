package warehouse

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.GET("/bays", GetBayList)
	g.GET("/bays/:id", GetBayByID)
	g.PUT("/bays/:id", UpdateBay)
	g.POST("/bays", NewBay)
	g.DELETE("/bays/:id", DeleteBay)

	g.GET("/locations", GetLocationList)
	g.GET("/locations/:id", GetLocationByID)
	g.PUT("/locations/:id", UpdateLocation)
	g.POST("/locations", NewLocation)
	g.DELETE("/locations/:id", DeleteLocation)
}
