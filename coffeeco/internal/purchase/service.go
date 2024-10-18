package purchase

import (
	"coffeeco/internal/loyalty"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"context"
	"errors"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type CardChargeService interface {
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}

type Repository interface {
	Store(ctx context.Context, purchase Purchase) error
}

type Service interface {
	CompletePurchase(ctx context.Context, purchase Purchase, coffeeBuxCard *loyalty.CoffeeBux) error
}

// serviceImpl is the example of domain service that orchestrates the purchase process
type serviceImpl struct {
	cardService  CardChargeService
	purchaseRepo Repository
	storeService store.Service
}

func NewService(cardService CardChargeService, repository Repository, storeService store.Service) Service {
	return &serviceImpl{
		cardService:  cardService,
		purchaseRepo: repository,
		storeService: storeService,
	}
}

// CompletePurchase coffeeBuxCard is optional since customer might or might not have loyalty card
func (s *serviceImpl) CompletePurchase(ctx context.Context, purchase Purchase, coffeeBuxCard *loyalty.CoffeeBux) error {
	purchase, err := purchase.ValidateAndEnrich()
	if err != nil {
		return err
	}

	discount, err := s.calculateStoreSpecificDiscount(purchase.Store.ID)
	if err != nil {
		return fmt.Errorf("err calculating store specific discount: %w", err)
	}
	purchase = purchase.ApplyDiscount(discount)

	switch purchase.PaymentMeans {
	case payment.CARD:
		if err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.CardToken); err != nil {
			return fmt.Errorf("err charging card: %w", err)
		}
	case payment.CASH:
		return fmt.Errorf("unsupported operation") //TODO
	case payment.COFFEEBUX:
		if coffeeBuxCard == nil {
			return fmt.Errorf("no coffeebux card presented")
		}
		updatedCard, err := coffeeBuxCard.Pay(purchase.Products)
		if err != nil {
			return fmt.Errorf("err paying with coffeebux: %w", err)
		}
		*coffeeBuxCard = updatedCard
	default:
		return fmt.Errorf("unexpected payment type")
	}

	if err := s.purchaseRepo.Store(ctx, purchase); err != nil {
		return fmt.Errorf("err storing purchase %w", err)
	}

	//this is a potential bug that needs to be consulted with domain experts,
	//if you pay with coffeebux card, should you get a stamp?
	if coffeeBuxCard != nil {
		*coffeeBuxCard = coffeeBuxCard.AddStamp()
	}
	return nil
}

func (s *serviceImpl) calculateStoreSpecificDiscount(storeId uuid.UUID) (float32, error) {
	discount, err := s.storeService.GetStoreSpecificDiscount(storeId)
	if err != nil {
		if errors.Is(err, store.ErrNoDiscount) {
			discount = 0
		} else {
			return 0, fmt.Errorf("err getting store specific discount: %w", err)
		}
	}
	return discount, nil
}
