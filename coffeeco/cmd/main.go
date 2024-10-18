package main

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/loyalty"
	"coffeeco/internal/payment"
	"coffeeco/internal/purchase"
	"coffeeco/internal/store"
	"context"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"log"
)

func main() {

	//stores
	storeRepository := store.NewInMemoryRepository(map[uuid.UUID]store.Store{
		uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"): {
			ID:                  uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479"),
			Location:            "New York",
			ProductsForSale:     nil,
			DiscountForProducts: floatPtr(50),
		},
	})
	storeService := store.NewService(storeRepository)

	//purchases
	cardService, _ := payment.NewStripeService("secret-key")
	cashService := payment.NewNotifyingCashPayment()
	purchaseRepository := purchase.NewInMemoryRepository(make(map[uuid.UUID]purchase.Purchase))
	purchaseService := purchase.NewService(cardService, cashService, purchaseRepository, storeService)

	p := purchase.Purchase{
		Store:        store.Store{ID: uuid.MustParse("f47ac10b-58cc-0372-8567-0e02b2c3d479")},
		PaymentMeans: payment.COFFEEBUX,
		Products: []coffeeco.Product{
			{
				ItemName:  "Latte",
				BasePrice: decimal.MustParse("3.50"),
			},
		},
	}
	l := loyalty.CoffeeBux{
		ID:                                    uuid.New(),
		FreeDrinksAvailable:                   2,
		RemainingDrinkPurchasesUntilFreeDrink: 9,
	}

	err := purchaseService.CompletePurchase(context.Background(), p, &l)
	if err != nil {
		log.Fatalf("err completing purchase: %v", err)
	}

	log.Printf("purchase completed successfully, remaining free drinks: %d", l.FreeDrinksAvailable)

	p.PaymentMeans = payment.CASH
	err = purchaseService.CompletePurchase(context.Background(), p, nil)
	if err != nil {
		log.Fatalf("err completing purchase: %v", err)
	}
}

func floatPtr(f float32) *float32 {
	return &f
}
