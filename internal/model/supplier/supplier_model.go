package supplier

import "time"

type SupplierModel struct {
	ID          int        `gorm:"primaryKey;autoIncrement;column:id"`
	CompanyName string     `gorm:"column:company_name"`
	Name        string     `gorm:"column:name"`
	Email       string     `gorm:"column:email"`
	Phone       string     `gorm:"column:phone"`
	Address     string     `gorm:"column:address"`
	IsActive    int16      `gorm:"column:is_active"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at"`
}

func (SupplierModel) TableName() string {
	return "supplier"
}
