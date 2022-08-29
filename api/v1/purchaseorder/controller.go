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
	purchaseorder.User = claims.UserName
	purchaseorder.Email = claims.Email
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
// @Param purchaseorder_number query string false "采购单号码"
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
	info.User = claims.UserName
	info.Email = claims.Email
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
	err := purchaseorderService.DeletePurchaseorder(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 采购单产品列表
// @Id 406
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "采购单ID"
// @Success 200 object response.ListRes{data=[]PurchaseorderItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id/items [GET]
func GetPurchaseorderItemList(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	list, err := purchaseorderService.GetPurchaseorderItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 更新采购单为ISSUED
// @Id 407
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "采购单ID"
// @Success 200 object response.ListRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id/issued [POST]
func IssuePurchaseorder(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	err := purchaseorderService.IssuePurchaseorder(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新建收货单
// @Id 408
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param purchasereceive_info body PurchasereceiveNew true "采购单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id/receives [POST]
func NewPurchasereceive(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var purchasereceive PurchasereceiveNew
	if err := c.ShouldBindJSON(&purchasereceive); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchasereceive.OrganizationID = claims.OrganizationID
	purchasereceive.User = claims.UserName
	purchasereceive.Email = claims.Email
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.NewPurchasereceive(uri.ID, purchasereceive)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 收货单列表
// @Id 409
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param purchasereceive_number query string false "采购单号码"
// @Success 200 object response.ListRes{data=[]PurchasereceiveResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchasereceives [GET]
func GetPurchasereceiveList(c *gin.Context) {
	var filter PurchasereceiveFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	count, list, err := purchaseorderService.GetPurchasereceiveList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 收货单商品列表
// @Id 410
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "收货单ID"
// @Success 200 object response.ListRes{data=[]PurchasereceiveItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchasereceives/:id/items [GET]
func GetPurchasereceiveItemList(c *gin.Context) {
	var uri PurchasereceiveID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	list, err := purchaseorderService.GetPurchaseReceiveItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 收货单详情列表
// @Id 411
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "收货单ID"
// @Success 200 object response.ListRes{data=[]PurchasereceiveDetailResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchasereceives/:id/details [GET]
func GetPurchasereceiveDetailList(c *gin.Context) {
	var uri PurchasereceiveID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	list, err := purchaseorderService.GetPurchaseReceiveDetailList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}
