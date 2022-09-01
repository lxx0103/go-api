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
// @Param unit_type query string false "单位类型（weight/length/custom)"
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

// @Summary 供应商列表
// @Id 316
// @Tags 供应商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "供应商名称"
// @Success 200 object response.ListRes{data=[]VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors [GET]
func GetVendorList(c *gin.Context) {
	var filter VendorFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetVendorList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建供应商
// @Id 317
// @Tags 供应商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param vendor_info body VendorNew true "供应商信息"
// @Success 200 object response.SuccessRes{data=VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors [POST]
func NewVendor(c *gin.Context) {
	var info VendorNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewVendor(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新供应商
// @Id 318
// @Tags 供应商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "供应商ID"
// @Param vendor_info body VendorNew true "供应商信息"
// @Success 200 object response.SuccessRes{data=Vendor} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [PUT]
func UpdateVendor(c *gin.Context) {
	var uri VendorID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info VendorNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateVendor(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取供应商
// @Id 319
// @Tags 供应商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "供应商ID"
// @Success 200 object response.SuccessRes{data=VendorResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [GET]
func GetVendorByID(c *gin.Context) {
	var uri VendorID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	vendor, err := settingService.GetVendorByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, vendor)

}

// @Summary 根据ID删除供应商
// @Id 320
// @Tags 供应商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "供应商ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /vendors/:id [DELETE]
func DeleteVendor(c *gin.Context) {
	var uri VendorID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteVendor(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 税率列表
// @Id 321
// @Tags 税率管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "税率名称"
// @Success 200 object response.ListRes{data=[]TaxResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /taxes [GET]
func GetTaxList(c *gin.Context) {
	var filter TaxFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetTaxList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建税率
// @Id 322
// @Tags 税率管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param tax_info body TaxNew true "税率信息"
// @Success 200 object response.SuccessRes{data=TaxResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /taxes [POST]
func NewTax(c *gin.Context) {
	var info TaxNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewTax(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新税率
// @Id 323
// @Tags 税率管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "税率ID"
// @Param tax_info body TaxNew true "税率信息"
// @Success 200 object response.SuccessRes{data=Tax} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /taxes/:id [PUT]
func UpdateTax(c *gin.Context) {
	var uri TaxID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info TaxNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateTax(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取税率
// @Id 324
// @Tags 税率管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "税率ID"
// @Success 200 object response.SuccessRes{data=TaxResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /taxes/:id [GET]
func GetTaxByID(c *gin.Context) {
	var uri TaxID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	tax, err := settingService.GetTaxByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, tax)

}

// @Summary 根据ID删除税率
// @Id 325
// @Tags 税率管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "税率ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /taxes/:id [DELETE]
func DeleteTax(c *gin.Context) {
	var uri TaxID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteTax(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 客户列表
// @Id 326
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "客户名称"
// @Success 200 object response.ListRes{data=[]CustomerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /customers [GET]
func GetCustomerList(c *gin.Context) {
	var filter CustomerFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetCustomerList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建客户
// @Id 327
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param customer_info body CustomerNew true "客户信息"
// @Success 200 object response.SuccessRes{data=CustomerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /customers [POST]
func NewCustomer(c *gin.Context) {
	var info CustomerNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewCustomer(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新客户
// @Id 328
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Param customer_info body CustomerNew true "客户信息"
// @Success 200 object response.SuccessRes{data=Customer} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /customers/:id [PUT]
func UpdateCustomer(c *gin.Context) {
	var uri CustomerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info CustomerNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateCustomer(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取客户
// @Id 329
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Success 200 object response.SuccessRes{data=CustomerResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /customers/:id [GET]
func GetCustomerByID(c *gin.Context) {
	var uri CustomerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	customer, err := settingService.GetCustomerByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, customer)

}

// @Summary 根据ID删除客户
// @Id 330
// @Tags 客户管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "客户ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /customers/:id [DELETE]
func DeleteCustomer(c *gin.Context) {
	var uri CustomerID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteCustomer(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 物流商列表
// @Id 331
// @Tags 物流商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param name query string false "物流商名称"
// @Success 200 object response.ListRes{data=[]CarrierResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /carriers [GET]
func GetCarrierList(c *gin.Context) {
	var filter CarrierFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	count, list, err := settingService.GetCarrierList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 新建物流商
// @Id 332
// @Tags 物流商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param carrier_info body CarrierNew true "物流商信息"
// @Success 200 object response.SuccessRes{data=CarrierResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /carriers [POST]
func NewCarrier(c *gin.Context) {
	var info CarrierNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.NewCarrier(info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID更新物流商
// @Id 333
// @Tags 物流商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "物流商ID"
// @Param carrier_info body CarrierNew true "物流商信息"
// @Success 200 object response.SuccessRes{data=Carrier} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /carriers/:id [PUT]
func UpdateCarrier(c *gin.Context) {
	var uri CarrierID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info CarrierNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.Email
	info.OrganizationID = claims.OrganizationID
	settingService := NewSettingService()
	new, err := settingService.UpdateCarrier(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取物流商
// @Id 334
// @Tags 物流商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "物流商ID"
// @Success 200 object response.SuccessRes{data=CarrierResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /carriers/:id [GET]
func GetCarrierByID(c *gin.Context) {
	var uri CarrierID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	carrier, err := settingService.GetCarrierByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, carrier)

}

// @Summary 根据ID删除物流商
// @Id 335
// @Tags 物流商管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path int true "物流商ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /carriers/:id [DELETE]
func DeleteCarrier(c *gin.Context) {
	var uri CarrierID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	settingService := NewSettingService()
	err := settingService.DeleteCarrier(uri.ID, claims.OrganizationID, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}
