package warehouse

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 货架列表
// @Id 501
// @Tags 货架管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param code query string false "货架编码"
// @Success 200 object response.ListRes{data=[]BayResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bays [GET]
func GetBayList(c *gin.Context) {
	var filter BayFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	warehouseService := NewWarehouseService()
	count, list, err := warehouseService.GetBayList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建货架
// @Id 502
// @Tags 货架管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param bay_info body BayNew true "货架信息"
// @Success 200 object response.SuccessRes{data=Bay} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bays [POST]
func NewBay(c *gin.Context) {
	var bay BayNew
	if err := c.ShouldBindJSON(&bay); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	bay.User = claims.Email
	bay.OrganizationID = claims.OrganizationID
	warehouseService := NewWarehouseService()
	new, err := warehouseService.NewBay(bay)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取货架
// @Id 503
// @Tags 货架管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "货架ID"
// @Success 200 object response.SuccessRes{data=Bay} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bays/:id [GET]
func GetBayByID(c *gin.Context) {
	var uri BayID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	warehouseService := NewWarehouseService()
	bay, err := warehouseService.GetBayByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, bay)

}

// @Summary 根据ID更新货架
// @Id 504
// @Tags 货架管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "货架ID"
// @Param bay_info body BayNew true "货架信息"
// @Success 200 object response.SuccessRes{data=Bay} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bays/:id [PUT]
func UpdateBay(c *gin.Context) {
	var uri BayID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var bay BayNew
	if err := c.ShouldBindJSON(&bay); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	bay.User = claims.Email
	bay.OrganizationID = claims.OrganizationID
	warehouseService := NewWarehouseService()
	new, err := warehouseService.UpdateBay(uri.ID, bay)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除货架
// @Id 505
// @Tags 货架管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "货架ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /bays/:id [DELETE]
func DeleteBay(c *gin.Context) {
	var uri BayID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	warehouseService := NewWarehouseService()
	err := warehouseService.DeleteBay(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 货位列表
// @Id 506
// @Tags 货位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数"
// @Param code query string false "货位编码"
// @Param bay_id query string false "货架ID"
// @Success 200 object response.ListRes{data=[]BayResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /locations [GET]
func GetLocationList(c *gin.Context) {
	var filter LocationFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	warehouseService := NewWarehouseService()
	count, list, err := warehouseService.GetLocationList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建货位
// @Id 507
// @Tags 货位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param location_info body LocationNew true "货位信息"
// @Success 200 object response.SuccessRes{data=Location} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /locations [POST]
func NewLocation(c *gin.Context) {
	var location LocationNew
	if err := c.ShouldBindJSON(&location); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	location.User = claims.Email
	location.OrganizationID = claims.OrganizationID
	warehouseService := NewWarehouseService()
	new, err := warehouseService.NewLocation(location)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据Code获取货位
// @Id 508
// @Tags 货位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param code path string true "货位code"
// @Success 200 object response.SuccessRes{data=Location} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /locations/:code [GET]
func GetLocationByCode(c *gin.Context) {
	var uri LocationCode
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	warehouseService := NewWarehouseService()
	location, err := warehouseService.GetLocationByCode(claims.OrganizationID, uri.Code)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, location)

}

// @Summary 根据ID更新货位
// @Id 509
// @Tags 货位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "货位ID"
// @Param location_info body LocationNew true "货位信息"
// @Success 200 object response.SuccessRes{data=Location} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /locations/:id [PUT]
func UpdateLocation(c *gin.Context) {
	var uri LocationID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var location LocationNew
	if err := c.ShouldBindJSON(&location); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	location.User = claims.Email
	location.OrganizationID = claims.OrganizationID
	warehouseService := NewWarehouseService()
	new, err := warehouseService.UpdateLocation(uri.ID, location)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID删除货位
// @Id 510
// @Tags 货位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "货位ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /locations/:id [DELETE]
func DeleteLocation(c *gin.Context) {
	var uri LocationID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	warehouseService := NewWarehouseService()
	err := warehouseService.DeleteLocation(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
