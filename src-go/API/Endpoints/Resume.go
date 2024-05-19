package Endpoints

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pebbe/zmq4"
	"io"
	"net/http"
)

func ResumeProcess(c *gin.Context, socket *zmq4.Socket) {
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

	// Receive the multipart response from Python
	parts, err := socket.RecvMessageBytes(0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive reply from ZeroMQ: " + err.Error()})
		return
	}

	if len(parts) != 3 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unexpected number of response parts received"})
		return
	}

	// Return the three strings as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"messageOne":   string(parts[0]),
		"messageTwo":   string(parts[1]),
		"messageThree": string(parts[2]),
	})
}
