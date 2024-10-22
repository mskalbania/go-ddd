package recommendation

import (
	"context"
	"fmt"
	"github.com/govalues/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type RecommendationSuite struct {
	suite.Suite
	service          Service
	availabilityMock *mockAvailabilityService
}

func TestRecommendationSuite(t *testing.T) {
	suite.Run(t, new(RecommendationSuite))
}

func (s *RecommendationSuite) SetupTest() {
	s.availabilityMock = &mockAvailabilityService{}
	s.service = New(s.availabilityMock)
}

func (s *RecommendationSuite) TestShouldRejectZeroBudget() {
	//when
	_, err := s.service.Get(context.Background(), time.Now(), time.Now().AddDate(0, 0, 1), "", decimal.Zero)

	//then
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "budget can't be 0")
}

func (s *RecommendationSuite) TestShouldRejectNegativeDurationBetweenDates() {
	//when
	start, _ := time.Parse(time.RFC3339, "2021-01-05T11:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2021-01-01T00:22:00Z")
	_, err := s.service.Get(context.Background(), start, end, "", decimal.Hundred)

	//then
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "trip duration can't be lower than 0")
}

func (s *RecommendationSuite) TestShouldRejectDurationLowerThanOneDay() {
	//when
	start, _ := time.Parse(time.RFC3339, "2021-01-05T11:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2021-01-05T15:59:00Z")
	_, err := s.service.Get(context.Background(), start, end, "", decimal.Hundred)

	//then
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "trip duration can't be lower than 1 day")

}

func (s *RecommendationSuite) TestShouldFailWhenAvailabilityFails() {
	//given
	s.availabilityMock.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]Option{}, fmt.Errorf("err")).Once()

	//when
	_, err := s.service.Get(context.Background(), time.Now(), time.Now().AddDate(0, 0, 1), "", decimal.Hundred)

	//then
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "failed to get availability options: err")
	s.availabilityMock.AssertExpectations(s.T())
}

func (s *RecommendationSuite) TestShouldFailWhenNoOptionsAvailable() {
	//given
	s.availabilityMock.On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]Option{}, nil).Once()

	//when
	_, err := s.service.Get(context.Background(), time.Now(), time.Now().AddDate(0, 0, 1), "", decimal.Hundred)

	//then
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "no available recommendation found")
	s.availabilityMock.AssertExpectations(s.T())
}

func (s *RecommendationSuite) TestShouldReturnCheapestOptionWithinBudget() {
	//given
	options := []Option{
		{HotelName: "Hotel A", Location: "PL", TripPrice: decimal.MustParse("10.23")},
		{HotelName: "Hotel B", Location: "UK", TripPrice: decimal.MustParse("5.10")},
		{HotelName: "Hotel C", Location: "US", TripPrice: decimal.MustParse("100")},
	}
	start := time.Now().Truncate(24 * time.Hour)
	end := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	s.availabilityMock.On("Get", mock.Anything, start, end, "UK").Return(options, nil).Once()

	//when
	recommendation, err := s.service.Get(context.Background(), start, end, "UK", decimal.MustParse("50"))

	//then
	require.NoError(s.T(), err)
	require.Equal(s.T(), "Hotel B", recommendation.HotelName)
	require.Equal(s.T(), "UK", recommendation.Location)
	require.Equal(s.T(), decimal.MustParse("5.10"), recommendation.TripPrice)
	require.Equal(s.T(), start, recommendation.TripStart)
	require.Equal(s.T(), end, recommendation.TripEnd)
	s.availabilityMock.AssertExpectations(s.T())
}

func (s *RecommendationSuite) TestShouldReturnErrorWhenNoOptionFoundWithinBudget() {
	//given
	options := []Option{
		{HotelName: "Hotel A", Location: "PL", TripPrice: decimal.MustParse("10.23")},
		{HotelName: "Hotel B", Location: "UK", TripPrice: decimal.MustParse("5.10")},
	}
	start := time.Now().Truncate(24 * time.Hour)
	end := time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour)
	s.availabilityMock.On("Get", context.Background(), start, end, mock.Anything).Return(options, nil).Once()

	//when
	_, err := s.service.Get(context.Background(), start, end, "UK", decimal.MustParse("5"))

	//then
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "no available recommendation found within budget")
	s.availabilityMock.AssertExpectations(s.T())
}

type mockAvailabilityService struct {
	mock.Mock
}

func (m *mockAvailabilityService) Get(ctx context.Context, from time.Time, to time.Time, location string) ([]Option, error) {
	args := m.Called(ctx, from, to, location)
	return args.Get(0).([]Option), args.Error(1)
}
