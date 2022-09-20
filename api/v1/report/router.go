package report

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.GET("/salesreports", GetSalesReport)
	g.GET("/purchasereports", GetPurchaseReport)
	g.GET("/adjustmentreports", GetAdjustmentReport)
	g.GET("/itemreports", GetItemReport)

}
