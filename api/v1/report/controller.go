package report

import (
	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

// @Summary 销售订单报告
// @Id 901
// @Tags 报告管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param date_from query string false "开始日期"
// @Param date_to query string false "结束日期"
// @Param customer_id query string false "顾客ID"
// @Success 200 object response.SuccessRes{data=SalesReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesreports [GET]
func GetSalesReport(c *gin.Context) {
	var filter SalesReportFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	reportService := NewReportService()
	res, err := reportService.GetSalesReport(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}

// @Summary 采购订单报告
// @Id 902
// @Tags 报告管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param date_from query string false "开始日期"
// @Param date_to query string false "结束日期"
// @Param vendor_id query string false "生产商ID"
// @Success 200 object response.SuccessRes{data=PurchaseReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /salesreports [GET]
func GetPurchaseReport(c *gin.Context) {
	var filter PurchaseReportFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	reportService := NewReportService()
	res, err := reportService.GetPurchaseReport(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}

// @Summary 库存调整报告
// @Id 903
// @Tags 报告管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param date_from query string false "开始日期"
// @Param date_to query string false "结束日期"
// @Param item_id query string false "商品ID"
// @Success 200 object response.SuccessRes{data=[]AdjustmentReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /adjustmentreports [GET]
func GetAdjustmentReport(c *gin.Context) {
	var filter AdjustmentReportFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	reportService := NewReportService()
	res, err := reportService.GetAdjustmentReport(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}

// @Summary 商品报告
// @Id 904
// @Tags 报告管理
// @version 1.0
// @Accept application/json
// @Produce application/json
// @Param item_id query string false "商品ID"
// @Success 200 object response.SuccessRes{data=[]AdjustmentReportResponse} 成功
// @Failure 400 object response.ErrorRes 内部错误
// @Router /adjustmentreports [GET]
func GetItemReport(c *gin.Context) {
	var filter ItemReportFilter
	err := c.ShouldBindQuery(&filter)
	if err != nil {
		response.ResponseError(c, "BindingError", err)
		return
	}
	claims := c.MustGet("claims").(*service.CustomClaims)
	filter.OrganizationID = claims.OrganizationID
	reportService := NewReportService()
	res, err := reportService.GetItemReport(filter)
	if err != nil {
		response.ResponseError(c, "DatabaseError", err)
		return
	}
	response.Response(c, res)
}
