package purchase

import (
	"context"
)

type Repository interface {
	Store(ctx context.Context, purchase Purchase) error
	Ping(ctx context.Context) error
}
