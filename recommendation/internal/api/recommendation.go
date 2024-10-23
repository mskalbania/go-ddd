package api

import (
	"github.com/gin-gonic/gin"
	"github.com/govalues/decimal"
	"microservice/internal/recommendation"
	"time"
)

func GetRecommendationHandler(service recommendation.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		from, err := time.Parse("2006-01-02", c.Query("from"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid from date"})
			return
		}
		to, err := time.Parse("2006-01-02", c.Query("to"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid to date"})
			return
		}
		location := c.Query("location")
		if location == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": "location is required"})
			return
		}
		budget, err := decimal.Parse(c.Query("budget"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "invalid budget"})
			return
		}
		r, err := service.Get(c, from, to, location, budget)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(500, gin.H{"error": "internal error"})
			return
		}
		c.JSON(200, recommendationResponse{
			Location:  r.Location,
			HotelName: r.HotelName,
			From:      r.TripStart.Format("2006-01-02"),
			To:        r.TripEnd.Format("2006-01-02"),
			Total:     budget.Rescale(2).String(),
		})
	}
}

type recommendationResponse struct {
	Location  string `json:"location"`
	HotelName string `json:"hotel_name"`
	From      string `json:"from"`
	To        string `json:"to"`
	Total     string `json:"total"`
}
