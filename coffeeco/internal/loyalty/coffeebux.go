package loyalty

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/store"
	"fmt"
	"github.com/google/uuid"
)

type CoffeeBux struct {
	ID                                    uuid.UUID
	store                                 store.Store
	coffeeLover                           coffeeco.CoffeeLover
	FreeDrinksAvailable                   int
	RemainingDrinkPurchasesUntilFreeDrink int
}

func (b CoffeeBux) AddStamp() CoffeeBux {
	if b.RemainingDrinkPurchasesUntilFreeDrink == 1 {
		b.FreeDrinksAvailable++
		b.RemainingDrinkPurchasesUntilFreeDrink = 10
	} else {
		b.RemainingDrinkPurchasesUntilFreeDrink--
	}
	return b
}

func (b CoffeeBux) Pay(products []coffeeco.Product) (CoffeeBux, error) {
	if len(products) == 0 {
		return CoffeeBux{}, fmt.Errorf("no products to purchase")
	}
	if len(products) > b.FreeDrinksAvailable {
		return CoffeeBux{}, fmt.Errorf("not enough free drinks available")
	}
	b.FreeDrinksAvailable -= len(products)
	return b, nil
}
