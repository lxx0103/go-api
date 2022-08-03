package setting

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 单位列表
// @Id 301
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "单位名称"
// @Success 200 object response.ListRes{data=[]UnitResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units [GET]
func GetUnitList(c *gin.Context) {
	var filter UnitFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetUnitList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建单位
// @Id 302
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param unit_info body UnitNew true "单位信息"
// @Success 200 object response.SuccessRes{data=UnitResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units [POST]
func NewUnit(c *gin.Context) {
	var info UnitNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewUnit(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新单位
// @Id 303
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "单位ID"
// @Param unit_info body UnitNew true "单位信息"
// @Success 200 object response.SuccessRes{data=Unit} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units/:id [PUT]
func UpdateUnit(c *gin.Context) {
	var uri UnitID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info UnitNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateUnit(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取单位
// @Id 304
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "单位ID"
// @Success 200 object response.SuccessRes{data=UnitResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units/:id [GET]
func GetUnitByID(c *gin.Context) {
	var uri UnitID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	unit, err := settingService.GetUnitByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, unit)

}

// @Summary 根据ID删除单位
// @Id 305
// @Tags 单位管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "单位ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /units/:id [DELETE]
func DeleteUnit(c *gin.Context) {
	var uri UnitID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteUnit(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 生产商列表
// @Id 306
// @Tags 生产商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "生产商名称"
// @Success 200 object response.ListRes{data=[]ManufacturerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /manufacturers [GET]
func GetManufacturerList(c *gin.Context) {
	var filter ManufacturerFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetManufacturerList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建生产商
// @Id 307
// @Tags 生产商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param manufacturer_info body ManufacturerNew true "生产商信息"
// @Success 200 object response.SuccessRes{data=ManufacturerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /manufacturers [POST]
func NewManufacturer(c *gin.Context) {
	var info ManufacturerNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewManufacturer(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新生产商
// @Id 308
// @Tags 生产商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "生产商ID"
// @Param manufacturer_info body ManufacturerNew true "生产商信息"
// @Success 200 object response.SuccessRes{data=Manufacturer} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /manufacturers/:id [PUT]
func UpdateManufacturer(c *gin.Context) {
	var uri ManufacturerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info ManufacturerNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateManufacturer(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取生产商
// @Id 309
// @Tags 生产商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "生产商ID"
// @Success 200 object response.SuccessRes{data=ManufacturerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /manufacturers/:id [GET]
func GetManufacturerByID(c *gin.Context) {
	var uri ManufacturerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	manufacturer, err := settingService.GetManufacturerByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, manufacturer)

}

// @Summary 根据ID删除生产商
// @Id 310
// @Tags 生产商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "生产商ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /manufacturers/:id [DELETE]
func DeleteManufacturer(c *gin.Context) {
	var uri ManufacturerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteManufacturer(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 品牌列表
// @Id 311
// @Tags 品牌管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "品牌名称"
// @Success 200 object response.ListRes{data=[]BrandResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands [GET]
func GetBrandList(c *gin.Context) {
	var filter BrandFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetBrandList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建品牌
// @Id 312
// @Tags 品牌管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param brand_info body BrandNew true "品牌信息"
// @Success 200 object response.SuccessRes{data=BrandResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands [POST]
func NewBrand(c *gin.Context) {
	var info BrandNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewBrand(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新品牌
// @Id 313
// @Tags 品牌管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "品牌ID"
// @Param brand_info body BrandNew true "品牌信息"
// @Success 200 object response.SuccessRes{data=Brand} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands/:id [PUT]
func UpdateBrand(c *gin.Context) {
	var uri BrandID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info BrandNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateBrand(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取品牌
// @Id 314
// @Tags 品牌管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "品牌ID"
// @Success 200 object response.SuccessRes{data=BrandResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands/:id [GET]
func GetBrandByID(c *gin.Context) {
	var uri BrandID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	brand, err := settingService.GetBrandByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, brand)

}

// @Summary 根据ID删除品牌
// @Id 315
// @Tags 品牌管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "品牌ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /brands/:id [DELETE]
func DeleteBrand(c *gin.Context) {
	var uri BrandID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteBrand(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
