package supplier

import (
	"context"
	"errors"
	"strings"
	"time"

	entity "be_evindo/internal/entity/supplier"
	repo "be_evindo/internal/repository/supplier"
)

type Usecase interface {
	Create(ctx context.Context, data *entity.Supplier) error
	FindAll(ctx context.Context) ([]*entity.Supplier, error)
	FindByID(ctx context.Context, id int) (*entity.Supplier, error)
	Update(ctx context.Context, id int, updates map[string]interface{}) error
	Delete(ctx context.Context, id int) error
}

type usecase struct {
	repo repo.SupplierRepository
}

func NewUsecase(repository repo.SupplierRepository) Usecase {
	return &usecase{repo: repository}
}

func (usecase *usecase) Create(ctx context.Context, data *entity.Supplier) error {
	if data == nil {
		return errors.New("supplier is required")
	}
	data.CompanyName = strings.TrimSpace(data.CompanyName)
	data.Name = strings.TrimSpace(data.Name)
	data.Email = strings.ToLower(strings.TrimSpace(data.Email))
	data.Phone = strings.TrimSpace(data.Phone)
	data.Address = strings.TrimSpace(data.Address)
	if data.CompanyName == "" || data.Name == "" || data.Email == "" {
		return errors.New("company_name, name, and email are required")
	}
	return usecase.repo.Create(ctx, data)
}

func (usecase *usecase) FindAll(ctx context.Context) ([]*entity.Supplier, error) {
	return usecase.repo.FindAll(ctx)
}

func (usecase *usecase) FindByID(ctx context.Context, id int) (*entity.Supplier, error) {
	if id <= 0 {
		return nil, errors.New("supplier id must be greater than zero")
	}
	return usecase.repo.FindByID(ctx, id)
}

func (usecase *usecase) Update(ctx context.Context, id int, updates map[string]interface{}) error {
	if id <= 0 {
		return errors.New("supplier id must be greater than zero")
	}
	if len(updates) == 0 {
		return errors.New("at least one field is required")
	}
	updates["updated_at"] = time.Now()
	for _, field := range []string{"company_name", "name", "email", "phone", "address"} {
		if value, exists := updates[field]; exists {
			text, ok := value.(string)
			if !ok || strings.TrimSpace(text) == "" {
				return errors.New(field + " must not be empty")
			}
			updates[field] = strings.TrimSpace(text)
		}
	}
	if value, exists := updates["email"]; exists {
		updates["email"] = strings.ToLower(value.(string))
	}
	return usecase.repo.Update(ctx, id, updates)
}

func (usecase *usecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("supplier id must be greater than zero")
	}
	return usecase.repo.Delete(ctx, id)
}
