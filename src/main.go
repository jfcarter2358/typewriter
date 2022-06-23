// main.go

package main

import (
	"fmt"
	"log"
	"typewriter/config"
	"typewriter/repo"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	config.LoadConfig()

	repo.DownloadContents()

	routerPort := fmt.Sprintf(":%v", config.Config.Port)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Initialize the routes
	initializeRoutes()

	// router.RunTLS(routerPort, "./certs/cert.crt", "./certs/cert.key")

	log.Printf("Running on port: %v", config.Config.Port)

	// Start serving the application
	router.Run(routerPort)
}
