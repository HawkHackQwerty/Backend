package ZeroMQ

import (
	"fmt"
	"github.com/pebbe/zmq4"
)

// HandleConnection establishes a new ZeroMQ REQ socket and returns it.
func InitReqSockets(ports []string) ([]*zmq4.Socket, error) {
	context, err := zmq4.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create ZMQ context: %v", err)
	}

	var sockets []*zmq4.Socket
	for _, port := range ports {
		socket, err := context.NewSocket(zmq4.REQ)
		if err != nil {
			// Properly close all sockets created before error
			for _, sock := range sockets {
				sock.Close()
			}
			return nil, fmt.Errorf("failed to create REQ socket for port %s: %v", port, err)
		}
		err = socket.Connect("tcp://python:" + port)
		if err != nil {
			socket.Close()
			// Properly close all sockets created before error
			for _, sock := range sockets {
				sock.Close()
			}
			return nil, fmt.Errorf("failed to connect REQ socket to port %s: %v", port, err)
		}
		sockets = append(sockets, socket)
	}

	return sockets, nil
}
