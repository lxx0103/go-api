package common

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 历史记录列表
// @Id 701
// @Tags 历史记录管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param page_id query int true "页码"
// @Param page_size query int true "每页行数（5/10/15/20）"
// @Param reference_id query string true "相关ID"
// @Param history_type query string true "历史类型"
// @Success 200 object response.ListRes{data=[]HistoryResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /historys [GET]
func GetHistoryList(c *gin.Context) {
	var filter HistoryFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	commonService := NewCommonService()
	count, list, err := commonService.GetHistoryList(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.ResponseList(c, filter.PageID, filter.PageSize, count, list)
}

// @Summary 获取下一个编码
// @Id 702
// @Tags 通用
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param number_type query string true "编码类型"
// @Success 200 object response.SuccessRes{data=string} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /nextnumber [GET]
func GetNextNumber(c *gin.Context) {
	var filter NumberFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	commonService := NewCommonService()
	res, err := commonService.GetNextNumber(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}
