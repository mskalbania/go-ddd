package recommendation

import (
	"context"
	"fmt"
	"github.com/govalues/decimal"
	"time"
)

type Recommendation struct {
	TripStart time.Time
	TripEnd   time.Time
	HotelName string
	Location  string
	TripPrice decimal.Decimal
}

type Option struct {
	HotelName string
	Location  string
	TripPrice decimal.Decimal
}

// Service from recommendation returns cheapest available trip option based on availability provided by other service
type Service interface {
	Get(ctx context.Context, from time.Time, to time.Time, location string, budget decimal.Decimal) (Recommendation, error)
}

// AvailabilityService is infrastructure service example, domain cares about contract not coupled with impl
type AvailabilityService interface {
	Get(ctx context.Context, from time.Time, to time.Time, location string) ([]Option, error)
}

type serviceImpl struct {
	availabilityService AvailabilityService
}

func New(availabilityService AvailabilityService) Service {
	return &serviceImpl{
		availabilityService: availabilityService,
	}
}

func (s *serviceImpl) Get(ctx context.Context, from time.Time, to time.Time, location string, budget decimal.Decimal) (Recommendation, error) {
	if budget.Equal(decimal.Zero) {
		return Recommendation{}, fmt.Errorf("budget can't be 0")
	}
	from = from.Truncate(24 * time.Hour)
	to = to.Truncate(24 * time.Hour)
	duration := to.Sub(from)
	if int64(duration) < 0 {
		return Recommendation{}, fmt.Errorf("trip duration can't be lower than 0")
	}
	nights := int64(duration.Hours() / 24)
	if nights < 1 {
		return Recommendation{}, fmt.Errorf("trip duration can't be lower than 1 day")
	}
	options, err := s.availabilityService.Get(ctx, from, to, location)
	if err != nil {
		return Recommendation{}, fmt.Errorf("failed to get availability options: %w", err)
	}
	if len(options) == 0 {
		return Recommendation{}, fmt.Errorf("no available recommendation found")
	}
	cheapest := Option{
		TripPrice: decimal.MustNew(9999999999, 0),
	}
	anyFound := false
	for _, option := range options {
		cost, _ := option.TripPrice.Mul(decimal.MustNew(nights, 0))
		if cost.Less(cheapest.TripPrice) && (cost.Less(budget) || cost.Equal(budget)) {
			cheapest = option
			anyFound = true
		}
	}
	if !anyFound {
		return Recommendation{}, fmt.Errorf("no available recommendation found within budget")
	}
	return Recommendation{
		TripStart: from,
		TripEnd:   to,
		HotelName: cheapest.HotelName,
		Location:  cheapest.Location,
		TripPrice: cheapest.TripPrice,
	}, nil
}
