package warehouse

import "time"

type WarehouseModel struct {
	ID        int        `gorm:"primaryKey;autoIncrement;column:id"`
	Code      string     `gorm:"column:code"`
	Name      string     `gorm:"column:name"`
	Location  string     `gorm:"column:location"`
	IsActive  int16      `gorm:"column:is_active"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
}

func (WarehouseModel) TableName() string {
	return "warehouse"
}
