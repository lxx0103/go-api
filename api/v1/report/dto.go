package report

type SalesReportFilter struct {
	DateFrom       string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo         string `form:"date_to" binding:"required,datetime=2006-01-02"`
	CustomerID     string `form:"customer_id" binding:"omitempty,max=64"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
}

type SalesReportResponse struct {
	CustomerReports []CustomerSalesReportResponse `json:"customer_reports"`
	Count           int                           `json:"count"`
	Total           float64                       `json:"total"`
	TaxTotal        float64                       `json:"tax_total"`
}

type CustomerSalesReportResponse struct {
	CustomerID   string                  `json:"customer_id"`
	CustomerName string                  `json:"customer_name"`
	Invoices     []InvoiceReportResponse `json:"invoices"`
	InvoiceCount int                     `json:"invoice_count"`
	Total        float64                 `json:"total"`
	TaxTotal     float64                 `json:"tax_total"`
}

type InvoiceReportResponse struct {
	InvoiceNumber string  `json:"invoice_number" db:"invoice_number"`
	InvoiceDate   string  `json:"invoice_date" db:"invoice_date"`
	CustomerID    string  `json:"customer_id" db:"customer_id"`
	Status        int     `json:"status" db:"status"`
	Total         float64 `json:"total" db:"total"`
	TaxTotal      float64 `json:"tax_total" db:"tax_total"`
}

// purchase report
type PurchaseReportFilter struct {
	DateFrom       string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo         string `form:"date_to" binding:"required,datetime=2006-01-02"`
	VendorID       string `form:"vendor_id" binding:"omitempty,max=64"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
}

type PurchaseReportResponse struct {
	VendorReports []VendorPurchaseReportResponse `json:"vendor_reports"`
	Count         int                            `json:"count"`
	Total         float64                        `json:"total"`
	TaxTotal      float64                        `json:"tax_total"`
}

type VendorPurchaseReportResponse struct {
	VendorID   string               `json:"vendor_id"`
	VendorName string               `json:"vendor_name"`
	Bills      []BillReportResponse `json:"bills"`
	BillCount  int                  `json:"bill_count"`
	Total      float64              `json:"total"`
	TaxTotal   float64              `json:"tax_total"`
}

type BillReportResponse struct {
	BillNumber string  `json:"bill_number" db:"bill_number"`
	BillDate   string  `json:"bill_date" db:"bill_date"`
	VendorID   string  `json:"vendor_id" db:"vendor_id"`
	Status     int     `json:"status" db:"status"`
	Total      float64 `json:"total" db:"total"`
	TaxTotal   float64 `json:"tax_total" db:"tax_total"`
}

// adjustment report
type AdjustmentReportFilter struct {
	DateFrom       string `form:"date_from" binding:"required,datetime=2006-01-02"`
	DateTo         string `form:"date_to" binding:"required,datetime=2006-01-02"`
	ItemID         string `form:"item_id" binding:"omitempty,max=64"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
}

type AdjustmentReportResponse struct {
	LocationCode         string `db:"location_code" json:"location_code"`
	ItemName             string `db:"item_name" json:"item_name"`
	SKU                  string `db:"sku" json:"sku"`
	Quantity             int    `db:"quantity" json:"quantity"`
	OriginalQuantiy      int    `db:"original_quantity" json:"original_quantity"`
	NewQuantiy           int    `db:"new_quantity" json:"new_quantity"`
	AdjustmentDate       string `db:"adjustment_date" json:"adjustment_date"`
	AdjustmentReasonName string `db:"adjustment_reason_name" json:"adjustment_reason_name"`
}

// item report
type ItemReportFilter struct {
	ItemID         string `form:"item_id" binding:"omitempty,max=64"`
	OrganizationID string `json:"organiztion_id" swaggerignore:"true"`
}

type ItemReportResponse struct {
	ItemName     string  `db:"item_name" json:"item_name"`
	SKU          string  `db:"sku" json:"sku"`
	Unit         string  `db:"unit" json:"unit"`
	StockOnHand  int     `db:"stock_on_hand" json:"stock_on_hand"`
	SellingPrice float64 `db:"selling_price" json:"selling_price"`
	CostPrice    float64 `db:"cost_price" json:"cost_price"`
}
