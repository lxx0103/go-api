package salesorder

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/salesorders", NewSalesorder)
	g.GET("/salesorders", GetSalesorderList)
	g.PUT("/salesorders/:id", UpdateSalesorder)
	g.GET("/salesorders/:id", GetSalesorderByID)
	g.DELETE("/salesorders/:id", DeleteSalesorder)
	g.GET("/salesorders/:id/items", GetSalesorderItemList)
	g.POST("/salesorders/:id/confirmed", ConfirmSalesorder)

	g.POST("/salesorders/:id/pickings", NewPickingorder)
	g.POST("/pickingorders", BatchPickingorder)
	g.GET("/pickingorders", GetPickingorderList)
	g.GET("/pickingorders/:id/items", GetPickingorderItemList)
	g.GET("/pickingorders/:id/details", GetPickingorderDetailList)

}
