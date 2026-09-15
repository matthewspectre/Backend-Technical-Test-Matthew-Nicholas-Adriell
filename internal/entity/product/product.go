package product

import "time"

type Product struct {
	ID        int        `json:"id"`
	SKU       string     `json:"sku"`
	Name      string     `json:"name"`
	Unit      string     `json:"unit"`
	IsActive  int16      `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
