package gateway

import (
	"net/http"

	eventpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/event"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventHandler struct {
	eventClient eventpb.EventServiceClient
}

func NewEventHandler(eventClient eventpb.EventServiceClient) *EventHandler {
	return &EventHandler{
		eventClient: eventClient,
	}
}

func (h *EventHandler) GetEvents(c *gin.Context) {
	response, err := h.eventClient.GetEvents(
		c.Request.Context(),
		&eventpb.GetEventsRequest{},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EventHandler) GetEvent(c *gin.Context) {

	id := c.Param("id")
	response, err := h.eventClient.GetEvent(
		c.Request.Context(),
		&eventpb.GetEventRequest{
			Id: id,
		},
	)

	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *EventHandler) CreateEvent(c *gin.Context) {
	var request eventpb.CreateEventRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.eventClient.CreateEvent(
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

func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id := c.Param("id")

	var request eventpb.UpdateEventRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	request.Id = id

	response, err := h.eventClient.UpdateEvent(
		c.Request.Context(),
		&request,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id := c.Param("id")

	response, err := h.eventClient.DeleteEvent(
		c.Request.Context(),
		&eventpb.DeleteEventRequest{
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
