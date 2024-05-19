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
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to get file: %s", err.Error())
		return
	}

	// Open the file for reading
	src, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to open file: %s", err.Error())
		return
	}
	defer src.Close()

	// Read file content into a buffer
	var fileBuf bytes.Buffer
	if _, err = io.Copy(&fileBuf, src); err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %s", err.Error())
		return
	}

	// Send file content over ZeroMQ
	if _, err := socket.SendBytes(fileBuf.Bytes(), 0); err != nil {
		c.String(http.StatusInternalServerError, "Failed to send file over ZeroMQ: %s", err.Error())
		return
	}

	// receive and process a reply from the server
	reply, err := socket.RecvBytes(0)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to receive reply from ZeroMQ: %s", err.Error())
		return
	}

	var response struct {
		MessageOne string `json:"messageOne"`
	}
	if err := json.Unmarshal(reply, &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to parse server reply: " + err.Error(),
		})
		return
	}

	// Return the three strings as a JSON response
	c.JSON(http.StatusOK, gin.H{
		"messageOne": response.MessageOne,
	})
}
