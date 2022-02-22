package http

import (
	"github.com/CollActionteam/collaction_backend/internal/contact"
	"github.com/CollActionteam/collaction_backend/internal/models"
	"github.com/CollActionteam/collaction_backend/pkg/handler"
	"github.com/CollActionteam/collaction_backend/pkg/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ContactHandler struct {
	service contact.Service
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
		c.JSON(http.StatusBadRequest, handler.Response{
			Status: handler.StatusFail,
			Data:   gin.H{"error": err.Error()},
		})
		return
	}

	if err := h.service.SendEmail(c.Request.Context(), request); err != nil {
		c.JSON(http.StatusInternalServerError, handler.Response{
			Status: handler.StatusFail,
			Data:   gin.H{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, handler.Response{Status: handler.StatusSuccess, Data: nil})
}
