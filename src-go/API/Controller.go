package API

import (
	"Mesh_Mesh/API/Endpoints"
	"Mesh_Mesh/ZeroMQ"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pebbe/zmq4"
	"log"
	"net/http"
)

func HandleServer() (*gin.Engine, *zmq4.Socket) {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-User-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(leakBucket())
	r.Use(authMiddleware())
	r.Use(userIDMiddleware())

	ZMQHandler, err := ZeroMQ.HandleConnection()
	if err != nil {
		log.Fatal("Failed to connect ZeroMQ:", err)
	}

	r.GET("/processResume", func(a *gin.Context) {
		Endpoints.ResumeProcess(a, ZMQHandler)
	})
	r.GET("/processCover", func(a *gin.Context) {
		Endpoints.CoverLetterProcess(a, ZMQHandler)
	})
	r.GET("/processVideo", func(a *gin.Context) {
		Endpoints.VideoProcess(a, ZMQHandler)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r, ZMQHandler
}

func CloseServer(a *gin.Engine, c *zmq4.Socket) {
	a = nil
	ZeroMQ.CloseConnection(c)
	c = nil
}
