package http

import (
	"codea-it/internal/app/adapters/outbound"
	"codea-it/internal/app/domain/services"
	"github.com/gin-gonic/gin"
)

const base_url = "/api"

func RegisterRoutes(r *gin.Engine) {
	apiClient := outbound.NewKrakenAPIClient()
	repository := outbound.NewSQLiteRepository()
	service := services.NewLtpService(apiClient, repository)
	handler := NewHandler(service)

	r.GET(base_url+"/v1/ltp", handler.HandleRequest)
}

type Handler struct {
	service services.LtpService
}

func NewHandler(service services.LtpService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) HandleRequest(c *gin.Context) {
	result, err := h.service.ExtractTickerData()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}
