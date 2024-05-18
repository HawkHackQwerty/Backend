package main

import (
	"fmt"
	"github.com/pebbe/zmq4"
	"time"
)

func main() {
	context, _ := zmq4.NewContext()
	defer context.Term()

	// Function to handle sending and receiving messages
	go func() {
		socket, _ := context.NewSocket(zmq4.REQ)
		defer socket.Close()
		socket.Connect("tcp://localhost:5555")

		for i := 0; i < 10; i++ {
			msg := fmt.Sprintf("Hello %d", i)
			socket.Send(msg, 0)
			reply, _ := socket.Recv(0)
			fmt.Println("Received", reply)
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(15 * time.Second)
}
