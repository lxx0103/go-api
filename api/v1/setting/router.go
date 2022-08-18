package setting

import "github.com/gin-gonic/gin"

func AuthRouter(g *gin.RouterGroup) {
	g.POST("/units", NewUnit)
	g.GET("/units", GetUnitList)
	g.PUT("/units/:id", UpdateUnit)
	g.GET("/units/:id", GetUnitByID)
	g.DELETE("/units/:id", DeleteUnit)

	g.POST("/manufacturers", NewManufacturer)
	g.GET("/manufacturers", GetManufacturerList)
	g.PUT("/manufacturers/:id", UpdateManufacturer)
	g.GET("/manufacturers/:id", GetManufacturerByID)
	g.DELETE("/manufacturers/:id", DeleteManufacturer)

	g.POST("/brands", NewBrand)
	g.GET("/brands", GetBrandList)
	g.PUT("/brands/:id", UpdateBrand)
	g.GET("/brands/:id", GetBrandByID)
	g.DELETE("/brands/:id", DeleteBrand)

	g.POST("/vendors", NewVendor)
	g.GET("/vendors", GetVendorList)
	g.PUT("/vendors/:id", UpdateVendor)
	g.GET("/vendors/:id", GetVendorByID)
	g.DELETE("/vendors/:id", DeleteVendor)

	g.POST("/taxes", NewTax)
	g.GET("/taxes", GetTaxList)
	g.PUT("/taxes/:id", UpdateTax)
	g.GET("/taxes/:id", GetTaxByID)
	g.DELETE("/taxes/:id", DeleteTax)
}
