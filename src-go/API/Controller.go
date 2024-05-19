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

func HandleServer() (*gin.Engine, []*zmq4.Socket) {
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

	OpenPorts := []string{"5550", "5551", "5552", "5553"}
	reqSockets, err := ZeroMQ.InitReqSockets(OpenPorts)
	if err != nil {
		log.Fatalf("Failed to establish ZeroMQ REQ sockets: %v", err)
	}

	r.GET("/processResume", func(a *gin.Context) {
		Endpoints.ResumeProcess(a, reqSockets[0])
	})
	r.GET("/processCover", func(a *gin.Context) {
		Endpoints.CoverLetterProcess(a, reqSockets[1])
	})
	r.GET("/processVideo", func(a *gin.Context) {
		Endpoints.VideoProcess(a, reqSockets[2])
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/uploadJob", func(c *gin.Context) {
		Endpoints.JobInfoProcess(c, reqSockets[3])
	})

	return r, reqSockets
}

func CloseServer(a *gin.Engine, c []*zmq4.Socket) {
	a = nil
	defer func() {
		for _, socket := range c {
			socket.Close()
		}
	}()
	c = nil
}
