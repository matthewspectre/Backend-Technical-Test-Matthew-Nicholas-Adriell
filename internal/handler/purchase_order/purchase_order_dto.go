package purchase_order

type PurchaseOrderCreateRequest struct {
	SupplierID int `json:"supplier_id"`
}

type PurchaseOrderStatusUpdateRequest struct {
	Status string `json:"status"`
}

type PurchaseOrderResponse struct {
	ID                int64                       `json:"id"`
	PONumber          string                      `json:"po_number"`
	PurchaseRequestID int64                       `json:"purchase_request_id"`
	SupplierID        int                         `json:"supplier_id"`
	SupplierName      string                      `json:"supplier_name"`
	WarehouseID       int                         `json:"warehouse_id"`
	WarehouseName     string                      `json:"warehouse_name"`
	Status            string                      `json:"status"`
	Items             []PurchaseOrderItemResponse `json:"items"`
	CreatedAt         string                      `json:"created_at"`
	UpdatedAt         *string                     `json:"updated_at"`
}

type PurchaseOrderItemResponse struct {
	ID               int64  `json:"id"`
	ProductID        int    `json:"product_id"`
	ProductName      string `json:"product_name"`
	Unit             string `json:"unit"`
	OrderedQuantity  int    `json:"ordered_quantity"`
	ReceivedQuantity int    `json:"received_quantity"`
}

type PurchaseOrderListResponse struct {
	Data      []PurchaseOrderResponse `json:"data"`
	TotalData int                     `json:"total_data"`
}
