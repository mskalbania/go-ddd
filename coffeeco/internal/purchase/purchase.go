package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"fmt"
	"github.com/google/uuid"
	"github.com/govalues/decimal"
	"time"
)

// Purchase is immutable entity since we are referencing purchases by ID in our domain
type Purchase struct {
	id             uuid.UUID
	total          decimal.Decimal
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
	total := decimal.Zero
	for _, pr := range p.Products {
		total, _ = total.Add(pr.BasePrice)
	}
	if total.IsZero() {
		return Purchase{}, fmt.Errorf("zero total not allowed")
	}
	if p.PaymentMeans == payment.CARD && p.CardToken == nil {
		return Purchase{}, fmt.Errorf("selected card payment but token is nil")
	}
	p.id = uuid.New()
	p.timeOfPurchase = time.Now().UTC()
	p.total = total
	return p, nil
}

func (p Purchase) ApplyDiscount(value float32) (Purchase, error) {
	if value > 0 {
		multiplier, err := decimal.NewFromFloat64(float64((100 - value) / 100))
		if err != nil {
			return Purchase{}, err
		}
		discounted, err := p.total.Mul(multiplier)
		if err != nil {
			return Purchase{}, err
		}
		p.total = discounted
	}
	return p, nil
}
