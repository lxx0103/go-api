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
	g.GET("/purchasereceives", GetPurchasereceiveList)
	g.GET("/purchasereceives/:id/items", GetPurchasereceiveItemList)
	g.GET("/purchasereceives/:id/details", GetPurchasereceiveDetailList)
	g.DELETE("/purchasereceives/:id", DeletePurchasereceive)

	g.POST("/purchaseorders/:id/bills", NewBill)
	g.GET("/bills", GetBillList)
	g.GET("/bills/:id/items", GetBillItemList)
	g.PUT("/bills/:id", UpdateBill)
	g.DELETE("/bills/:id", DeleteBill)

	g.POST("/bills/:id/payments", NewPayment)
	g.GET("/bills/:id/paid", GeBillPaymentMade)
	g.GET("/paymentmades", GetPaymentList)
	g.PUT("/paymentmades/:id", UpdatePayment)
	g.DELETE("/paymentmades/:id", DeletePayment)

}
