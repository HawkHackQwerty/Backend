package Endpoints

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pebbe/zmq4"
	"net/http"
)

type RequestBody struct {
	UserID    string `json:"userID"` // Include UserID in the struct
	StringOne string `json:"stringOne"`
	StringTwo string `json:"stringTwo"`
}

func JobInfoProcess(c *gin.Context, socket *zmq4.Socket) {
	var requestBody RequestBody

	// Retrieve user ID from the Gin context, set by the middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is required"})
		return
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	// Assign userID from the context to the request body before sending
	requestBody.UserID = userID.(string)

	data, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode request data"})
		return
	}

	// Send the serialized data over ZeroMQ without waiting for a response
	if _, err := socket.SendBytes(data, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data over ZeroMQ: " + err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
