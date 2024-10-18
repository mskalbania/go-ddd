package store

import (
	"fmt"
	"github.com/google/uuid"
)

var (
	ErrStoreNotFound = fmt.Errorf("store not found")
)

type inMemoryRepository struct {
	stores map[uuid.UUID]Store
}

func NewInMemoryRepository(stores map[uuid.UUID]Store) Repository {
	return &inMemoryRepository{
		stores: stores,
	}
}

func (r *inMemoryRepository) GetStore(storeId uuid.UUID) (Store, error) {
	store, ok := r.stores[storeId]
	if !ok {
		return Store{}, ErrStoreNotFound
	}
	return store, nil
}
