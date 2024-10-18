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
	if total.IsZero() {
		return Purchase{}, fmt.Errorf("zero total not allowed")
	}
	if p.PaymentMeans == payment.CARD && p.CardToken == nil {
		return Purchase{}, fmt.Errorf("selected card payment but token is nil")
	}
	p.id = uuid.New()
	p.timeOfPurchase = time.Now().UTC()
	p.total = *total
	return p, nil
}

func (p Purchase) ApplyDiscount(value float32) Purchase {
	if value > 0 {
		p.total = *p.total.Multiply(int64(1 - value))
	}
	return p
}
