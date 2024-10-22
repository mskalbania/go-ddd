package api

import (
	"github.com/gin-gonic/gin"
	"microservice/internal/recommendation"
)

func GetRecommendationHandler(service recommendation.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		//TODO
	}
}
