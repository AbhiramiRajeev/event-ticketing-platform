package gateway

import (
	"net/http"

	"github.com/gin-gonic/gin"
	registrationpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/registration"
)

type RegistrationHandler struct {
	registrationClient registrationpb.RegistrationServiceClient
}

func NewRegistrationHandler(
	registrationClient registrationpb.RegistrationServiceClient,
) *RegistrationHandler {
	return &RegistrationHandler{
		registrationClient: registrationClient,
	}
}

func (h *RegistrationHandler) CreateRegistration(c *gin.Context) {
	var request registrationpb.CreateRegistrationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.registrationClient.CreateRegistration(
		c.Request.Context(),
		&request,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *RegistrationHandler) GetRegistration(c *gin.Context) {
	id := c.Param("id")

	response, err := h.registrationClient.GetRegistration(
		c.Request.Context(),
		&registrationpb.GetRegistrationRequest{
			Id: id,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *RegistrationHandler) GetRegistrations(c *gin.Context) {
	userID := c.Param("userID")

	response, err := h.registrationClient.GetRegistrations(
		c.Request.Context(),
		&registrationpb.GetRegistrationsRequest{
			UserId: userID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *RegistrationHandler) CancelRegistration(c *gin.Context) {
	id := c.Param("id")

	response, err := h.registrationClient.CancelRegistration(
		c.Request.Context(),
		&registrationpb.CancelRegistrationRequest{
			Id: id,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}