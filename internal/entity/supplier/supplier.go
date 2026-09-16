package supplier

import "time"

type Supplier struct {
	ID          int        `json:"id"`
	CompanyName string     `json:"company_name"`
	Name        string     `json:"name"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Address     string     `json:"address"`
	IsActive    int16      `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
