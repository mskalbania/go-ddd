package payment

import (
	"context"
	"github.com/govalues/decimal"
	"log"
)

type NotifyingCashPaymentService struct{}

func NewNotifyingCashPayment() *NotifyingCashPaymentService {
	return &NotifyingCashPaymentService{}
}

func (n *NotifyingCashPaymentService) PayCash(ctx context.Context, amount decimal.Decimal) error {
	log.Printf("\n\nPaying cash: %.2f \n\n", amount)
	//emit event
	return nil
}
