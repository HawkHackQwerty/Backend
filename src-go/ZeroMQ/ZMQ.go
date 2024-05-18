package ZeroMQ

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pebbe/zmq4"
)

// Private Constants _____________________________________________________________________________________________________________

const (
	address = "tcp://python:5555"
)

// End of Private Constants __________________________________________________________________________________________________________

// Public Types _____________________________________________________________________________________________________________

type JobInformation struct {
	JobTitle       string
	JobDescription string
}

type JobFeedback struct {
	ResumeFeedback      string
	Score               uint8
	UpdatedCoverLetter  string
	LinkedInColdMessage string
}

// End of Public Types __________________________________________________________________________________________________________

// Public Functions _____________________________________________________________________________________________________________

// HandleConnection establishes a new ZeroMQ REQ socket and returns it.
func HandleConnection() (*zmq4.Socket, error) {
	socket, err := zmq4.NewSocket(zmq4.REQ)
	if err != nil {
		return nil, err
	}
	err = socket.Connect(address)
	if err != nil {
		return nil, err
	}
	return socket, nil
}

// CloseConnection closes the given ZeroMQ socket.
func CloseConnection(socket *zmq4.Socket) {
	socket.Close()
}

// SendJobInformation sends job information to the server and receives feedback
func SendJobInformation(socket *zmq4.Socket, job JobInformation) (JobFeedback, error) {
	// Set up Input into byte buffers
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, int32(len(job.JobTitle))); err != nil {
		return JobFeedback{}, fmt.Errorf("failed to write JobTitle length: %v", err)
	}
	buf.WriteString(job.JobTitle)

	if err := binary.Write(buf, binary.LittleEndian, int32(len(job.JobDescription))); err != nil {
		return JobFeedback{}, fmt.Errorf("failed to write JobDescription length: %v", err)
	}
	buf.WriteString(job.JobDescription)

	// Send Information
	_, err := socket.SendBytes(buf.Bytes(), 0)
	if err != nil {
		return JobFeedback{}, err
	}

	// Receive and deserialize JobFeedback
	reply, err := socket.RecvBytes(0)
	if err != nil {
		return JobFeedback{}, err
	}
	replyBuf := bytes.NewReader(reply)
	var feedback JobFeedback
	var score uint8
	binary.Read(replyBuf, binary.LittleEndian, &feedback.ResumeFeedback)
	binary.Read(replyBuf, binary.LittleEndian, &score)
	feedback.Score = score
	binary.Read(replyBuf, binary.LittleEndian, &feedback.UpdatedCoverLetter)
	binary.Read(replyBuf, binary.LittleEndian, &feedback.LinkedInColdMessage)

	return feedback, nil
}

// End of Public Functions __________________________________________________________________________________________________________
