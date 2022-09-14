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

	g.POST("/pickingorders/:id/pick", NewPickingFromLocation)
	g.POST("/pickingorders/:id/picked", MarkPickingorderPicked)
	g.POST("/pickingorders/:id/unpicked", MarkPickingorderUnPicked)
	g.DELETE("/pickingorders/:id", DeletePickingorder)

	g.POST("/salesorders/:id/packages", NewPackage)
	g.GET("/packages", GetPackageList)
	g.GET("/packages/:id/items", GetPackageItemList)
	g.DELETE("/packages/:id", DeletePackage)

	g.POST("/shippingorders", BatchShippingorder)
	g.GET("/shippingorders", GetShippingorderList)
	g.GET("/shippingorders/:id/items", GetShippingorderItemList)
	g.GET("/shippingorders/:id/details", GetShippingorderDetailList)
	g.DELETE("/shippingorders/:id", DeleteShippingorder)

	g.GET("/requisitions", GetRequisitionList)

	g.POST("/salesorders/:id/invoices", NewInvoice)
	g.GET("/invoices", GetInvoiceList)
	g.GET("/invoices/:id/items", GetInvoiceItemList)
	g.PUT("/invoices/:id", UpdateInvoice)
	g.DELETE("/invoices/:id", DeleteInvoice)

}
