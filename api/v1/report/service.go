package report

import (
	"errors"
	"go-api/api/v1/setting"
	"go-api/core/database"
)

type reportService struct {
}

func NewReportService() *reportService {
	return &reportService{}
}

//sales

func (s *reportService) GetSalesReport(filter SalesReportFilter) (*SalesReportResponse, error) {
	db := database.RDB()
	query := NewReportQuery(db)
	settingQuery := setting.NewSettingQuery(db)
	invoices, err := query.GetSalesReport(filter)
	if err != nil {
		return nil, err
	}
	var res SalesReportResponse
	var customers []CustomerSalesReportResponse
	count := 0
	total := 0.0
	taxTotal := 0.0
	for _, invoice := range *invoices {
		customerExist := false
		for idx, customer := range customers {
			if invoice.CustomerID == customer.CustomerID {
				customers[idx].Invoices = append(customer.Invoices, invoice)
				customers[idx].Total += invoice.Total
				customers[idx].TaxTotal += invoice.TaxTotal
				customers[idx].InvoiceCount += 1
				count += 1
				total += invoice.Total
				taxTotal += invoice.TaxTotal
				customerExist = true
				break
			}
		}
		if !customerExist {
			var newCustomer CustomerSalesReportResponse
			newCustomer.CustomerID = invoice.CustomerID
			customerInfo, err := settingQuery.GetCustomerByID(filter.OrganizationID, invoice.CustomerID)
			if err != nil {
				msg := "get customer info error"
				return nil, errors.New(msg)
			}
			newCustomer.CustomerName = customerInfo.Name
			newCustomer.InvoiceCount = 1
			newCustomer.Invoices = append(newCustomer.Invoices, invoice)
			newCustomer.Total = invoice.Total
			newCustomer.TaxTotal = invoice.TaxTotal
			customers = append(customers, newCustomer)
			count += 1
			total += invoice.Total
			taxTotal += invoice.TaxTotal
		}
	}
	res.Count = count
	res.TaxTotal = taxTotal
	res.Total = total
	res.CustomerReports = customers
	return &res, err
}

//purchase

func (s *reportService) GetPurchaseReport(filter PurchaseReportFilter) (*PurchaseReportResponse, error) {
	db := database.RDB()
	query := NewReportQuery(db)
	settingQuery := setting.NewSettingQuery(db)
	bills, err := query.GetPurchaseReport(filter)
	if err != nil {
		return nil, err
	}
	var res PurchaseReportResponse
	var vendors []VendorPurchaseReportResponse
	count := 0
	total := 0.0
	taxTotal := 0.0
	for _, invoice := range *bills {
		vendorExist := false
		for idx, vendor := range vendors {
			if invoice.VendorID == vendor.VendorID {
				vendors[idx].Bills = append(vendor.Bills, invoice)
				vendors[idx].Total += invoice.Total
				vendors[idx].TaxTotal += invoice.TaxTotal
				vendors[idx].BillCount += 1
				count += 1
				total += invoice.Total
				taxTotal += invoice.TaxTotal
				vendorExist = true
				break
			}
		}
		if !vendorExist {
			var newVendor VendorPurchaseReportResponse
			newVendor.VendorID = invoice.VendorID
			vendorInfo, err := settingQuery.GetVendorByID(filter.OrganizationID, invoice.VendorID)
			if err != nil {
				msg := "get vendor info error"
				return nil, errors.New(msg)
			}
			newVendor.VendorName = vendorInfo.Name
			newVendor.BillCount = 1
			newVendor.Bills = append(newVendor.Bills, invoice)
			newVendor.Total = invoice.Total
			newVendor.TaxTotal = invoice.TaxTotal
			vendors = append(vendors, newVendor)
			count += 1
			total += invoice.Total
			taxTotal += invoice.TaxTotal
		}
	}
	res.Count = count
	res.TaxTotal = taxTotal
	res.Total = total
	res.VendorReports = vendors
	return &res, err
}

//adjustment

func (s *reportService) GetAdjustmentReport(filter AdjustmentReportFilter) (*[]AdjustmentReportResponse, error) {
	db := database.RDB()
	query := NewReportQuery(db)
	res, err := query.GetAdjustmentReport(filter)
	return res, err
}

//item

func (s *reportService) GetItemReport(filter ItemReportFilter) (*[]ItemReportResponse, error) {
	db := database.RDB()
	query := NewReportQuery(db)
	res, err := query.GetItemReport(filter)
	return res, err
}
