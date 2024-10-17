package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
	"time"
)

// Purchase is immutable entity since we are referencing purchases by ID in our domain
type Purchase struct {
	id             uuid.UUID
	total          money.Money
	timeOfPurchase time.Time

	Store        store.Store
	PaymentMeans payment.Means
	Products     []coffeeco.Product
	CardToken    *string //might not be available when cash
}

// ValidateAndEnrich this logic belongs entirely to purchase
func (p Purchase) ValidateAndEnrich() (Purchase, error) {
	if len(p.Products) == 0 {
		return Purchase{}, fmt.Errorf("no products to purchase")
	}
	total := money.New(0, money.USD)
	for _, pr := range p.Products {
		total, _ = total.Add(&pr.BasePrice)
	}
	if p.total.IsZero() {
		return Purchase{}, fmt.Errorf("zero total not allowed")
	}
	if p.PaymentMeans == payment.CARD && p.CardToken == nil {
		return Purchase{}, fmt.Errorf("selected card payment but token is nil")
	}
	return Purchase{
		id:             uuid.New(),
		total:          *total,
		timeOfPurchase: time.Now().UTC(),
		Store:          p.Store,
		PaymentMeans:   p.PaymentMeans,
		Products:       p.Products,
		CardToken:      p.CardToken,
	}, nil
}
