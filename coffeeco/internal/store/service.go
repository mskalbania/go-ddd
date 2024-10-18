package store

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrNoDiscount = fmt.Errorf("no discount available")
)

type Repository interface {
	GetStore(storeId uuid.UUID) (Store, error)
}

type Service interface {
	GetStoreSpecificDiscount(storeId uuid.UUID) (float32, error)
}

type serviceImpl struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &serviceImpl{repository: repository}
}

func (s *serviceImpl) GetStoreSpecificDiscount(storeId uuid.UUID) (float32, error) {
	store, err := s.repository.GetStore(storeId)
	if err != nil {
		return 0, fmt.Errorf("failed to get store: %w", err)
	}
	if store.DiscountForProducts == nil {
		return 0, ErrNoDiscount
	}
	return *store.DiscountForProducts, nil
}
