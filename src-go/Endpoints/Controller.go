package Endpoints

import (
	"Mesh_Mesh/ZeroMQ"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pebbe/zmq4"
	"go.uber.org/ratelimit"
	"log"
	"net/http"
	"time"
)

var limit = ratelimit.New(100) // 100 requests per second

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(c *gin.Context) {
		now := limit.Take()
		log.Printf("Time since last request: %v", now.Sub(prev))
		prev = now
		c.Next()
	}
}

func HandleServer() (*gin.Engine, *zmq4.Socket) {
	r := gin.Default()

	// Configuring CORS middleware
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Applying the rate limiting middleware
	r.Use(leakBucket())

	ZMQHandler, err := ZeroMQ.HandleConnection()
	if err != nil {
		log.Fatal("Failed to connect ZeroMQ:", err)
	}

	// Authentication routes
	r.POST("/signup", SignUp)
	r.POST("/signin", SignIn)
	r.GET("/logout", Logout)

	// Example endpoint to test rate limiting
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	return r, ZMQHandler
}

func CloseServer(a *gin.Engine, c *zmq4.Socket) {
	a = nil
	ZeroMQ.CloseConnection(c)
	c = nil
}
