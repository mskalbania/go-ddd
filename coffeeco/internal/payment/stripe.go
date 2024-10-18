package payment

import (
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/client"
)

/*
Infrastructure layer
*/

type StripeService struct {
	stripeClient *client.API
}

func NewStripeService(apiKey string) (*StripeService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("cannot create stripe client - API key empty")
	}
	c := new(client.API)
	c.Init(apiKey, nil)
	return &StripeService{stripeClient: c}, nil
}

func (s *StripeService) ChargeCard(ctx context.Context, amount money.Money, cardToken string) error {
	params := &stripe.ChargeParams{
		Amount:   stripe.Int64(amount.Amount()),
		Currency: stripe.String(amount.Currency().Code),
		Source: &stripe.PaymentSourceSourceParams{
			Token: stripe.String(cardToken),
		},
	}
	_, err := s.stripeClient.Charges.New(params)
	if err != nil {
		return fmt.Errorf("unable to create charge: %w", err)
	}
	return nil
}
