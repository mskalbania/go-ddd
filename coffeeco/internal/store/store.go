package store

import (
	coffeeco "coffeeco/internal"
	"github.com/google/uuid"
)

// Store is entity, belongs to store domain
type Store struct {
	ID              uuid.UUID
	Location        string
	ProductsForSale []coffeeco.Product
}
