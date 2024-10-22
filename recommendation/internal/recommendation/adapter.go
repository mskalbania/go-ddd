package recommendation

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/govalues/decimal"
	"net/http"
	"time"
)

const partnershipPath = "/partnership"

// This is adapter pattern or anti-corruption layer,
// bridge between our bounded context (recommendation) and external bounded context (partnership)
// Our domain model is not polluted with external concerns.
type partnershipAvailabilityService struct {
	client  *http.Client
	baseUrl string
}

func NewPartnershipAvailability(client *http.Client, baseUrl string) AvailabilityService {
	return &partnershipAvailabilityService{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (p *partnershipAvailabilityService) Get(ctx context.Context, from time.Time, to time.Time, location string) ([]Option, error) {
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")
	rq, err := http.NewRequest(http.MethodGet, fmt.Sprintf(partnershipPath+"?from=%s&to=%s&location=%s", fromStr, toStr, location), nil)
	if err != nil {
		return nil, err
	}
	rsp, err := p.client.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("partnership service returned %d", rsp.StatusCode)
	}
	var availableHotels availableHotels
	if err := json.NewDecoder(rsp.Body).Decode(&availableHotels); err != nil {
		return nil, err
	}
	options := make([]Option, 0, len(availableHotels.Hotels))
	for _, hotel := range availableHotels.Hotels {
		price, err := decimal.NewFromFloat64(hotel.PricePerNightUSD)
		if err != nil {
			return nil, err
		}
		options = append(options, Option{
			HotelName: hotel.Name,
			Location:  location,
			TripPrice: price,
		})
	}
	return options, nil
}

type availableHotels struct {
	Hotels []hotel `json:"availableHotels"`
}

type hotel struct {
	Name             string  `json:"name"`
	PricePerNightUSD float64 `json:"priceInUSDPerNight"`
}
