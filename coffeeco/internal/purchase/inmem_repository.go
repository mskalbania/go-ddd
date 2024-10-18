package purchase

import (
	"context"
	"github.com/google/uuid"
)

type inMemoryRepository struct {
	purchases map[uuid.UUID]Purchase
}

func NewInMemoryRepository(purchases map[uuid.UUID]Purchase) Repository {
	return &inMemoryRepository{
		purchases: purchases,
	}
}

func (i *inMemoryRepository) Store(ctx context.Context, purchase Purchase) error {
	i.purchases[purchase.id] = purchase
	return nil
}
