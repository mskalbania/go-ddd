package subscription

import (
	coffeeco "coffeeco/internal"
	"fmt"
	"github.com/google/uuid"
)

type Subscription struct {
	ID            uuid.UUID
	coffeeLover   coffeeco.CoffeeLover
	freeDrinkLeft bool
}

func (s Subscription) Pay(products []coffeeco.Product) (Subscription, error) {
	if !s.freeDrinkLeft {
		return Subscription{}, fmt.Errorf("no more free drink available within sub")
	}
	if len(products) > 1 {
		return Subscription{}, fmt.Errorf("sub only allows for 1 free drink")
	}
	s.freeDrinkLeft = false
	return s, nil
}
