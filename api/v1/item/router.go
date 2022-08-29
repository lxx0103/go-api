package item

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/items", NewItem)
	g.GET("/items", GetItemList)
	g.PUT("/items/:id", UpdateItem)
	g.GET("/items/:id", GetItemByID)
	g.DELETE("/items/:id", DeleteItem)

	g.GET("/barcodes", GetBarcodeList)
	// g.GET("/barcodes/:id", GetBarcodeByID)
	g.GET("/barcodes/:code", GetBarcodeByCode)
	g.PUT("/barcodes/:id", UpdateBarcode)
	g.POST("/barcodes", NewBarcode)
	g.DELETE("/barcodes/:id", DeleteBarcode)
}
