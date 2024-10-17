package coffeeco

import "github.com/google/uuid"

// CoffeeLover is entity, it is put inside internal since it is required by all other domain code
type CoffeeLover struct {
	ID           uuid.UUID
	FirstName    string
	LastName     string
	EmailAddress string
}
