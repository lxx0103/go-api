package salesorder

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 新建销售单
// @Id 601
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param salesorder_info body SalesorderNew true "销售单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders [POST]
func NewSalesorder(c *gin.Context) {
	var salesorder SalesorderNew
	if err := c.ShouldBindJSON(&salesorder); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorder.OrganizationID = claims.OrganizationID
	salesorder.User = claims.UserName
	salesorder.Email = claims.Email
	salesorderService := NewSalesorderService()
	new, err := salesorderService.NewSalesorder(salesorder)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 销售单列表
// @Id 602
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "销售单名称"
// @Success 200 object response.ListRes{data=[]SalesorderResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders [GET]
func GetSalesorderList(c *gin.Context) {
	var filter SalesorderFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	salesorderService := NewSalesorderService()
	count, list, err := salesorderService.GetSalesorderList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 根据ID更新销售单
// @Id 603
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "销售单ID"
// @Param salesorder_info body SalesorderNew true "销售单信息"
// @Success 200 object response.SuccessRes{data=Salesorder} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id [PUT]
func UpdateSalesorder(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info SalesorderNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	info.Email = claims.Email
	info.OrganizationID = claims.OrganizationID
	salesorderService := NewSalesorderService()
	new, err := salesorderService.UpdateSalesorder(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取销售单
// @Id 604
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "销售单ID"
// @Success 200 object response.SuccessRes{data=SalesorderResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id [GET]
func GetSalesorderByID(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	salesorder, err := salesorderService.GetSalesorderByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, salesorder)

}

// @Summary 根据ID删除销售单
// @Id 605
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "销售单ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id [DELETE]
func DeleteSalesorder(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.DeleteSalesorder(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 销售单产品列表
// @Id 606
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "销售单ID"
// @Success 200 object response.ListRes{data=[]SalesorderItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id/items [GET]
func GetSalesorderItemList(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetSalesorderItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 更新销售单为ISSUED
// @Id 607
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "销售单ID"
// @Success 200 object response.ListRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id/confirmed [POST]
func ConfirmSalesorder(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.ConfirmSalesorder(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新建拣货单
// @Id 608
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param pickingorder body PickingorderNew true "销售单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id/pickings [POST]
func NewPickingorder(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var pickingorder PickingorderNew
	if err := c.ShouldBindJSON(&pickingorder); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	pickingorder.OrganizationID = claims.OrganizationID
	pickingorder.User = claims.UserName
	pickingorder.Email = claims.Email
	salesorderService := NewSalesorderService()
	new, err := salesorderService.NewPickingorder(uri.ID, pickingorder)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 拣货单列表
// @Id 609
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param pickingorder_number query string false "拣货单编码"
// @Param salesorder_id query string false "销售订单ID"
// @Success 200 object response.ListRes{data=[]PickingorderResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders [GET]
func GetPickingorderList(c *gin.Context) {
	var filter PickingorderFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	salesorderService := NewSalesorderService()
	count, list, err := salesorderService.GetPickingorderList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 拣货单产品列表
// @Id 610
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "拣货单ID"
// @Success 200 object response.ListRes{data=[]PickingorderItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders/:id/items [GET]
func GetPickingorderItemList(c *gin.Context) {
	var uri PickingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetPickingorderItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 拣货单详情列表
// @Id 611
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "拣货单ID"
// @Success 200 object response.ListRes{data=[]PickingorderDetailResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders/:id/details [GET]
func GetPickingorderDetailList(c *gin.Context) {
	var uri PickingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetPickingorderDetailList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 批量拣货
// @Id 612
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param purchasereceive_info body PickingorderBatch true "销售单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders [POST]
func BatchPickingorder(c *gin.Context) {
	var batchPickingorder PickingorderBatch
	if err := c.ShouldBindJSON(&batchPickingorder); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	batchPickingorder.OrganizationID = claims.OrganizationID
	batchPickingorder.User = claims.UserName
	batchPickingorder.Email = claims.Email
	salesorderService := NewSalesorderService()
	new, err := salesorderService.BatchPickingorder(batchPickingorder)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 从货位拣货
// @Id 613
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param purchasereceive_info body PickingFromLocationNew true "销售单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders/:id/pick [POST]
func NewPickingFromLocation(c *gin.Context) {
	var uri PickingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var pickingFromLocation PickingFromLocationNew
	if err := c.ShouldBindJSON(&pickingFromLocation); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	pickingFromLocation.OrganizationID = claims.OrganizationID
	pickingFromLocation.User = claims.UserName
	pickingFromLocation.Email = claims.Email
	salesorderService := NewSalesorderService()
	new, err := salesorderService.NewPickingFromLocation(uri.ID, pickingFromLocation)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 标记拣货单为已拣货
// @Id 614
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders/:id/picked [POST]
func MarkPickingorderPicked(c *gin.Context) {
	var uri PickingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.UpdatePickingorderPicked(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 新建包裹
// @Id 615
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param package body PackageNew true "销售单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesorders/:id/packages [POST]
func NewPackage(c *gin.Context) {
	var uri SalesorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var packageNew PackageNew
	if err := c.ShouldBindJSON(&packageNew); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	packageNew.OrganizationID = claims.OrganizationID
	packageNew.User = claims.UserName
	packageNew.Email = claims.Email
	salesorderService := NewSalesorderService()
	new, err := salesorderService.NewPackage(uri.ID, packageNew)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 包裹列表
// @Id 616
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param package_number query string false "包裹编码"
// @Param salesorder_id query string false "销售订单ID"
// @Success 200 object response.ListRes{data=[]PackageResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /packages [GET]
func GetPackageList(c *gin.Context) {
	var filter PackageFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	salesorderService := NewSalesorderService()
	count, list, err := salesorderService.GetPackageList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 拣货单产品列表
// @Id 617
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "拣货单ID"
// @Success 200 object response.ListRes{data=[]PackageItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /packages/:id/items [GET]
func GetPackageItemList(c *gin.Context) {
	var uri PackageID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetPackageItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 批量发货
// @Id 618
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param shippingorder_info body ShippingorderBatch true "销售单信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shippingorders [POST]
func BatchShippingorder(c *gin.Context) {
	var batchShippingorder ShippingorderBatch
	if err := c.ShouldBindJSON(&batchShippingorder); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	batchShippingorder.OrganizationID = claims.OrganizationID
	batchShippingorder.User = claims.UserName
	batchShippingorder.Email = claims.Email
	salesorderService := NewSalesorderService()
	new, err := salesorderService.BatchShippingorder(batchShippingorder)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 发货单列表
// @Id 619
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param shippingorder_number query string false "发货单编码"
// @Param package_id query string false "包裹ID"
// @Success 200 object response.ListRes{data=[]ShippingorderResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shippingorders [GET]
func GetShippingorderList(c *gin.Context) {
	var filter ShippingorderFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	salesorderService := NewSalesorderService()
	count, list, err := salesorderService.GetShippingorderList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 发货单产品列表
// @Id 620
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "拣货单ID"
// @Success 200 object response.ListRes{data=[]ShippingorderItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shippingorders/:id/items [GET]
func GetShippingorderItemList(c *gin.Context) {
	var uri ShippingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetShippingorderItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 发货单详情列表
// @Id 621
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "发货单ID"
// @Success 200 object response.ListRes{data=[]ShippingorderItemResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shippingorders/:id/details [GET]
func GetShippingorderDetailList(c *gin.Context) {
	var uri ShippingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetShippingorderItemList(uri.ID, claims.OrganizationID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 采购需求列表
// @Id 622
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param start_date query string false "开始时间"
// @Param end_date query string false "结束时间"
// @Param target_day query int false "目标时间"
// @Success 200 object response.ListRes{data=[]RequsitionResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /requisitions [GET]
func GetRequisitionList(c *gin.Context) {
	var filter RequsitionFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	salesorderService := NewSalesorderService()
	list, err := salesorderService.GetRequisitionList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, list)
}

// @Summary 根据ID删除发货单
// @Id 623
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "发货单ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /shippingorders/:id [DELETE]
func DeleteShippingorder(c *gin.Context) {
	var uri ShippingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.DeleteShippingorder(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 根据ID删除包裹
// @Id 624
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "包裹ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /packages/:id [DELETE]
func DeletePackage(c *gin.Context) {
	var uri PackageID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.DeletePackage(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 标记拣货单为未拣货
// @Id 625
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders/:id/unpicked [POST]
func MarkPickingorderUnPicked(c *gin.Context) {
	var uri PickingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.UpdatePickingorderUnPicked(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "ok")
}

// @Summary 根据ID删除拣货单
// @Id 626
// @Tags 销售单管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "拣货单ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /pickingorders/:id [DELETE]
func DeletePickingorder(c *gin.Context) {
	var uri PickingorderID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	salesorderService := NewSalesorderService()
	err := salesorderService.DeletePickingorder(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
