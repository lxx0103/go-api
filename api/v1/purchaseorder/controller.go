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

// @Summary 根据ID删除收货单
// @Id 412
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "收货单ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchasereceives/:id [DELETE]
func DeletePurchasereceive(c *gin.Context) {
	var uri PurchasereceiveID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	err := purchaseorderService.DeletePurchasereceive(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 新建Bill
// @Id 413
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param bill body BillNew true "采购单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /purchaseorders/:id/bills [POST]
func NewBill(c *gin.Context) {
	var uri PurchaseorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var bill BillNew
	if err := c.ShouldBindJSON(&bill); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	bill.OrganizationID = claims.OrganizationID
	bill.User = claims.UserName
	bill.Email = claims.Email
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.NewBill(uri.ID, bill)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary Bill列表
// @Id 414
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param bill_number query string false "bill编码"
// @Param purchaseorder_id query string false "销售订单ID"
// @Success 200 object response.ListRes{data=[]BillResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bills [GET]
func GetBillList(c *gin.Context) {
	var filter BillFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	count, list, err := purchaseorderService.GetBillList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary bill产品列表
// @Id 415
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "billID"
// @Success 200 object response.ListRes{data=[]BillItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bills/:id/items [GET]
func GetBillItemList(c *gin.Context) {
	var uri BillID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	list, err := purchaseorderService.GetBillItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 根据ID更新Bill
// @Id 416
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "billID"
// @Param bill_info body BillNew true "bill信息"
// @Success 200 object response.SuccessRes{data=Bill} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bills/:id [PUT]
func UpdateBill(c *gin.Context) {
	var uri BillID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info BillNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	info.Email = claims.Email
	info.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.UpdateBill(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除bill
// @Id 417
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "billID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bills/:id [DELETE]
func DeleteBill(c *gin.Context) {
	var uri BillID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	err := purchaseorderService.DeleteBill(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 新建Payment
// @Id 418
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param bill body PaymentMadeNew true "付款信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bills/:id/payments [POST]
func NewPayment(c *gin.Context) {
	var uri BillID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var paymentMade PaymentMadeNew
	if err := c.ShouldBindJSON(&paymentMade); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	paymentMade.OrganizationID = claims.OrganizationID
	paymentMade.User = claims.UserName
	paymentMade.Email = claims.Email
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.NewPaymentMade(uri.ID, paymentMade)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 获取Bill已付款金额
// @Id 419
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "billID"
// @Success 200 object response.SuccessRes{data=float64} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bills/:id/paid [GET]
func GeBillPaymentMade(c *gin.Context) {
	var uri BillID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	res, err := purchaseorderService.GeBillPaymentMade(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}

// @Summary Payment列表
// @Id 420
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param payment_received_number query string false "付款单编码"
// @Param bill_id query string false "billID"
// @Param payment_method_id query string false "付款方式ID"
// @Success 200 object response.ListRes{data=[]PaymentMadeResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentmades [GET]
func GetPaymentList(c *gin.Context) {
	var filter PaymentMadeFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	count, list, err := purchaseorderService.GetPaymentMadeList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 根据ID更新Payment
// @Id 421
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "paymentID"
// @Param payment_info body PaymentMadeNew true "payment信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentmades/:id [PUT]
func UpdatePayment(c *gin.Context) {
	var uri PaymentMadeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info PaymentMadeNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	info.Email = claims.Email
	info.OrganizationID = claims.OrganizationID
	purchaseorderService := NewPurchaseorderService()
	new, err := purchaseorderService.UpdatePaymentMade(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除payment
// @Id 422
// @Tags 采购单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "paymentID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /paymentmades/:id [DELETE]
func DeletePayment(c *gin.Context) {
	var uri PaymentMadeID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	purchaseorderService := NewPurchaseorderService()
	err := purchaseorderService.DeletePaymentMade(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
