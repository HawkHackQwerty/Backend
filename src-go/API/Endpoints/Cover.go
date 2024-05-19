package Endpoints

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pebbe/zmq4"
	"io"
	"net/http"
)

func CoverLetterProcess(c *gin.Context, socket *zmq4.Socket) {
	// Retrieve user ID from the Gin context, set by the middleware
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID is required"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file: " + err.Error()})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file: " + err.Error()})
		return
	}
	defer src.Close()

	var fileBuf bytes.Buffer
	if _, err = io.Copy(&fileBuf, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	// Create a structured message including the user ID and file data
	message := struct {
		UserID   string `json:"userID"`
		FileData []byte `json:"fileData"`
	}{
		UserID:   userID.(string),
		FileData: fileBuf.Bytes(),
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode message"})
		return
	}

	if _, err := socket.SendBytes(msgBytes, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data over ZeroMQ: " + err.Error()})
		return
	}

	// Receive the response from Python and return it in the HTTP response
	reply, err := socket.RecvBytes(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive reply from ZeroMQ: " + err.Error()})
		return
	}

	c.String(http.StatusOK, string(reply))
}
