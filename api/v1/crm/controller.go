package crm

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 新建线索
// @Id 801
// @Tags 客户资源管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param lead_info body LeadNew true "线索信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /leads [POST]
func NewLead(c *gin.Context) {
	var lead LeadNew
	if err := c.ShouldBindJSON(&lead); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	lead.OrganizationID = claims.OrganizationID
	lead.User = claims.UserName
	lead.Email = claims.Email
	crmService := NewCrmService()
	new, err := crmService.NewLead(lead)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 线索列表
// @Id 802
// @Tags 客户资源管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param company query string false "公司名称"
// @Success 200 object response.ListRes{data=[]LeadResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /leads [GET]
func GetLeadList(c *gin.Context) {
	var filter LeadFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	crmService := NewCrmService()
	count, list, err := crmService.GetLeadList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 根据ID更新线索
// @Id 803
// @Tags 客户资源管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "线索ID"
// @Param lead_info body LeadNew true "客户管理信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /leads/:id [PUT]
func UpdateLead(c *gin.Context) {
	var uri LeadID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info LeadNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	info.Email = claims.Email
	info.OrganizationID = claims.OrganizationID
	crmService := NewCrmService()
	new, err := crmService.UpdateLead(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}

// @Summary 根据ID获取线索
// @Id 804
// @Tags 客户资源管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "线索ID"
// @Success 200 object response.SuccessRes{data=LeadResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /leads/:id [GET]
func GetLeadByID(c *gin.Context) {
	var uri LeadID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	crmService := NewCrmService()
	crm, err := crmService.GetLeadByID(claims.OrganizationID, uri.ID)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, crm)

}

// @Summary 根据ID删除线索
// @Id 805
// @Tags 客户资源管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "线索ID"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /leads/:id [DELETE]
func DeleteLead(c *gin.Context) {
	var uri LeadID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	crmService := NewCrmService()
	err := crmService.DeleteLead(uri.ID, claims.OrganizationID, claims.UserName, claims.Email)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, "OK")
}

// @Summary 转化线索为顾客/供应商
// @Id 806
// @Tags 客户资源管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param id path string true "线索ID"
// @Param lead_info body LeadConvertNew true "客户管理信息"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /leads/:id/convert [POST]
func ConvertLead(c *gin.Context) {
	var uri LeadID
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	var info LeadConvertNew
	if err := c.ShouldBindJSON(&info); err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	info.User = claims.UserName
	info.Email = claims.Email
	info.OrganizationID = claims.OrganizationID
	crmService := NewCrmService()
	new, err := crmService.ConvertLead(uri.ID, info)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, new)
}
