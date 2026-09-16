package product

import "time"

type ProductModel struct {
	ID        int        `gorm:"primaryKey;autoIncrement;column:id"`
	SKU       string     `gorm:"column:sku"`
	Name      string     `gorm:"column:name"`
	Unit      string     `gorm:"column:unit"`
	IsActive  int16      `gorm:"column:is_active"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (ProductModel) TableName() string {
	return "product"
}
