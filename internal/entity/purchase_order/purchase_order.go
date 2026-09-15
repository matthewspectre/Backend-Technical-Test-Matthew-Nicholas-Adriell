package purchase_order

import "time"

type PurchaseOrder struct {
	ID                int64                `json:"id"`
	PONumber          string               `json:"po_number"`
	PurchaseRequestID int64                `json:"purchase_request_id"`
	SupplierID        int                  `json:"supplier_id"`
	SupplierName      string               `json:"supplier_name"`
	WarehouseID       int                  `json:"warehouse_id"`
	WarehouseName     string               `json:"warehouse_name"`
	Status            string               `json:"status"`
	Items             []*PurchaseOrderItem `json:"items"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         *time.Time           `json:"updated_at"`
}

type PurchaseOrderItem struct {
	ID               int64  `json:"id"`
	PurchaseOrderID  int64  `json:"purchase_order_id"`
	ProductID        int    `json:"product_id"`
	ProductName      string `json:"product_name"`
	Unit             string `json:"unit"`
	OrderedQuantity  int    `json:"ordered_quantity"`
	ReceivedQuantity int    `json:"received_quantity"`
}
