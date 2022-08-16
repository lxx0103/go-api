package purchaseorder

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/purchaseorders", NewPurchaseorder)
	g.GET("/purchaseorders", GetPurchaseorderList)
	g.PUT("/purchaseorders/:id", UpdatePurchaseorder)
	g.GET("/purchaseorders/:id", GetPurchaseorderByID)
	g.DELETE("/purchaseorders/:id", DeletePurchaseorder)
	g.GET("/purchaseorders/:id/items", GetPurchaseorderItemList)
	g.POST("/purchaseorders/:id/issued", IssuePurchaseorder)

	g.POST("/purchaseorders/:id/receives", NewPurchasereceive)

}
