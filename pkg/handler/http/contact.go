package http

import (
	"context"
	"github.com/CollActionteam/collaction_backend/internal/contact"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/pkg/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContactService interface {
	SendEmail(ctx context.Context, data models.EmailContactRequest) error
}

type ContactHandler struct {
	service ContactService
}

func NewContactHandler() *ContactHandler {
	emailRepository := repository.NewEmail("userName", "password")
	configManager := repository.NewConfigManager()
	return &ContactHandler{service: contact.NewContactService(emailRepository, configManager, gin.Mode())}
}

func (h *ContactHandler) Register(router *gin.Engine) {
	router.POST("/contact", h.EmailContact)
}

func (h *ContactHandler) EmailContact(c *gin.Context) {
	var request models.EmailContactRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := h.service.SendEmail(c.Request.Context(), request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message sent successfully"})
}
