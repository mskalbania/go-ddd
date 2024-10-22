package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-retryablehttp"
	"microservice/internal/api"
	"microservice/internal/recommendation"
)

func main() {
	c := retryablehttp.NewClient()
	c.RetryMax = 10
	partnershipAdapter := recommendation.NewPartnershipAvailability(c.StandardClient(), "localhost:3333")
	recommendationService := recommendation.New(partnershipAdapter)

	//this is essentially open host example - exposing bounded context to the outside world
	g := gin.Default()
	g.Handle("GET", "/recommendation", api.GetRecommendationHandler(recommendationService))
}
