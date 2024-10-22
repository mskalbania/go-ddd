package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// This is implementation of partnership dependency for our recommendation system
func main() {
	g := gin.Default()
	g.GET("/partnership", func(c *gin.Context) {
		location := c.Query("location")
		if location == "" {
			c.JSON(400, gin.H{"error": "location is required"})
			return
		}
		random := rand.Intn(10) + 1
		if random <= 3 {
			c.JSON(500, gin.H{"error": "partnership service is down"})
			return
		}
		c.JSON(200, availableHotels{
			Hotels: []hotel{
				{"Hotel A", 100},
				{"Hotel B", 200},
				{"Hotel C", 300},
			},
		})
	})
	http.ListenAndServe("localhost:3333", g)
}

type availableHotels struct {
	Hotels []hotel `json:"availableHotels"`
}

type hotel struct {
	Name             string  `json:"name"`
	PricePerNightUSD float64 `json:"priceInUSDPerNight"`
}
