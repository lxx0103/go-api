package item

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 新建商品
// @Id 201
// @Tags 商品管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param item_info body ItemNew true "商品信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /items [POST]
func NewItem(c *gin.Context) {
	var item ItemNew
	if err := c.ShouldBindJSON(&item); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	item.OrganizationID = claims.OrganizationID
	item.User = claims.UserName
	item.Email = claims.Email
	itemService := NewItemService()
	new, err := itemService.NewItem(item)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 商品列表
// @Id 202
// @Tags 商品管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "商品名称"
// @Success 200 object response.ListRes{data=[]ItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /items [GET]
func GetItemList(c *gin.Context) {
	var filter ItemFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	itemService := NewItemService()
	count, list, err := itemService.GetItemList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 根据ID更新商品
// @Id 203
// @Tags 商品管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "商品ID"
// @Param item_info body ItemNew true "商品信息"
// @Success 200 object response.SuccessRes{data=Item} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /items/:id [PUT]
func UpdateItem(c *gin.Context) {
	var uri ItemID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ItemNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	info.Email = claims.Email
	info.OrganizationID = claims.OrganizationID
	itemService := NewItemService()
	new, err := itemService.UpdateItem(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取商品
// @Id 204
// @Tags 商品管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "商品ID"
// @Success 200 object response.SuccessRes{data=ItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /items/:id [GET]
func GetItemByID(c *gin.Context) {
	var uri ItemID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	itemService := NewItemService()
	item, err := itemService.GetItemByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, item)

}

// @Summary 根据ID删除商品
// @Id 205
// @Tags 商品管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "商品ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /items/:id [DELETE]
func DeleteItem(c *gin.Context) {
	var uri ItemID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	itemService := NewItemService()
	err := itemService.DeleteItem(uri.ID, claims.OrganizationID, claims.Email, claims.UserName)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 条码列表
// @Id 206
// @Tags 条码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param code query string false "条码编码"
// @Param sku query string false "SKU"
// @Success 200 object response.ListRes{data=[]BarcodeResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /barcodes [GET]
func GetBarcodeList(c *gin.Context) {
	var filter BarcodeFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	itemService := NewItemService()
	count, list, err := itemService.GetBarcodeList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建条码
// @Id 207
// @Tags 条码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param barcode_info body BarcodeNew true "条码信息"
// @Success 200 object response.SuccessRes{data=Barcode} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /barcodes [POST]
func NewBarcode(c *gin.Context) {
	var barcode BarcodeNew
	if err := c.ShouldBindJSON(&barcode); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	barcode.User = claims.Email
	barcode.OrganizationID = claims.OrganizationID
	itemService := NewItemService()
	new, err := itemService.NewBarcode(barcode)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取条码
// @Id 208
// @Tags 条码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "条码ID"
// @Success 200 object response.SuccessRes{data=Barcode} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /barcodes/:id [GET]
func GetBarcodeByID(c *gin.Context) {
	var uri BarcodeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	itemService := NewItemService()
	barcode, err := itemService.GetBarcodeByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, barcode)

}

// @Summary 根据ID更新条码
// @Id 209
// @Tags 条码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "条码ID"
// @Param barcode_info body BarcodeNew true "条码信息"
// @Success 200 object response.SuccessRes{data=Barcode} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /barcodes/:id [PUT]
func UpdateBarcode(c *gin.Context) {
	var uri BarcodeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var barcode BarcodeNew
	if err := c.ShouldBindJSON(&barcode); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	barcode.User = claims.Email
	barcode.OrganizationID = claims.OrganizationID
	itemService := NewItemService()
	new, err := itemService.UpdateBarcode(uri.ID, barcode)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除条码
// @Id 210
// @Tags 条码管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "条码ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /barcodes/:id [DELETE]
func DeleteBarcode(c *gin.Context) {
	var uri BarcodeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	itemService := NewItemService()
	err := itemService.DeleteBarcode(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
