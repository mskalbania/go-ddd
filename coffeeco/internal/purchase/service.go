package purchase

import (
	"coffeeco/internal/payment"
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
)

type CardChargeService interface { //Infrastructure Service interface
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}

type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
}

// CompletePurchase value receiver since service is stateless and contains only reference values
func (s Service) CompletePurchase(ctx context.Context, purchase Purchase) error {
	purchase, err := purchase.ValidateAndEnrich()
	if err != nil {
		return err
	}
	switch purchase.PaymentMeans {
	case payment.CARD:
		if err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.CardToken); err != nil {
			return fmt.Errorf("err charging card: %w", err)
		}
	case payment.CASH:
		return fmt.Errorf("unsupported operation") //TODO
	case payment.COFFEEBUX:
		return fmt.Errorf("unsupported operation") //TODO
	default:
		return fmt.Errorf("unexpected payment type")
	}

	if err := s.purchaseRepo.Store(ctx, purchase); err != nil {
		return fmt.Errorf("err storing purchase %w", err)
	}
	return nil
}
