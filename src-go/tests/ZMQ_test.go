package tests

import (
	"Mesh_Mesh/ZeroMQ"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZMQConnection(t *testing.T) {
	socket, err := ZeroMQ.HandleConnection()
	defer ZeroMQ.CloseConnection(socket)
	assert.NoError(t, err, "Failed to establish a connection")
}

func TestZMQAsynchronousCall(t *testing.T) {
	socket, err := ZeroMQ.HandleConnection()
	if err != nil {
		t.Fatalf("Failed to establish a connection: %v", err)
	}
	defer ZeroMQ.CloseConnection(socket)

	// Create a JobInformation object
	jobInfo := ZeroMQ.JobInformation{
		JobTitle:       "Software Engineer",
		JobDescription: "Develop software applications",
	}

	// Send job information and expect feedback
	feedback, err := ZeroMQ.SendJobInformation(socket, jobInfo)
	assert.NoError(t, err, "Failed to send job information and receive feedback")
	assert.NotEmpty(t, feedback, "Feedback should not be empty")
}
