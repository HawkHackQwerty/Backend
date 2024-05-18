package main

import (
	"Mesh_Mesh/Endpoints"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	router, zmqHandler := Endpoints.HandleServer()

	go func() {
		// Listening on port 8080
		if err := router.Run(":8080"); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cleanup resources
	Endpoints.CloseServer(router, zmqHandler)
	log.Println("Server gracefully stopped")
}
