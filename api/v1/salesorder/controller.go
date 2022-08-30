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
