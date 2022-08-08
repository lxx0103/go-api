package purchaseorder

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 新建采购单
// @Id 401
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param purchaseorder_info body PurchaseorderNew true "采购单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders [POST]
func NewPurchaseorder(c *gin.Context) {
	var purchaseorder PurchaseorderNew
	if err := c.ShouldBindJSON(&purchaseorder); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorder.OrganizationID = claims.OrganizationID
	purchaseorder.User = claims.Email
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.NewPurchaseorder(purchaseorder)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 采购单列表
// @Id 402
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "采购单名称"
// @Success 200 object response.ListRes{data=[]PurchaseorderResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders [GET]
func GetPurchaseorderList(c *gin.Context) {
	var filter PurchaseorderFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	count, list, err := purchaseorderService.GetPurchaseorderList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 根据ID更新采购单
// @Id 403
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "采购单ID"
// @Param purchaseorder_info body PurchaseorderNew true "采购单信息"
// @Success 200 object response.SuccessRes{data=Purchaseorder} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id [PUT]
func UpdatePurchaseorder(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info PurchaseorderNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.UpdatePurchaseorder(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取采购单
// @Id 404
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "采购单ID"
// @Success 200 object response.SuccessRes{data=PurchaseorderResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id [GET]
func GetPurchaseorderByID(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	purchaseorder, err := purchaseorderService.GetPurchaseorderByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, purchaseorder)

}

// @Summary 根据ID删除采购单
// @Id 405
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "采购单ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id [DELETE]
func DeletePurchaseorder(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	err := purchaseorderService.DeletePurchaseorder(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 条码列表
// @Id 406
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
	purchaseorderService := NewPurchaseorderService()
	count, list, err := purchaseorderService.GetBarcodeList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建条码
// @Id 407
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
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.NewBarcode(barcode)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取条码
// @Id 408
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
	purchaseorderService := NewPurchaseorderService()
	barcode, err := purchaseorderService.GetBarcodeByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, barcode)

}

// @Summary 根据ID更新条码
// @Id 409
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
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.UpdateBarcode(uri.ID, barcode)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除条码
// @Id 410
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
	purchaseorderService := NewPurchaseorderService()
	err := purchaseorderService.DeleteBarcode(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}