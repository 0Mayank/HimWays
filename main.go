package main

import (
	"github.com/0Mayank/himways/configs"
	"github.com/0Mayank/himways/routes"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func main() {
	r := gin.Default()

	configs.ConnectDB()

	// routes
	routes.BusRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
