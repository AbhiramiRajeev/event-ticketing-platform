package gateway

import (
	"net/http"

	userpb "github.com/AbhiramiRajeev/event-ticketing-platform/proto/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userClient userpb.UserServiceClient
}

func NewUserHandler(userClient userpb.UserServiceClient) *UserHandler {
	return &UserHandler{
		userClient: userClient,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request userpb.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.userClient.CreateUser(
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

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	response, err := h.userClient.GetUser(
		c.Request.Context(),
		&userpb.GetUserRequest{
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

func (h *UserHandler) GetUsers(c *gin.Context) {
	response, err := h.userClient.GetUsers(
		c.Request.Context(),
		&userpb.GetUsersRequest{},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
