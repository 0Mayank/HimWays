package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Bus struct {
	Plate  string `json:"plate"`
	Number string `json:"name"`
	Route  []int  `json:"route"`
}

func GetBus(c *gin.Context) {
	var busParams struct {
		Plate   string `form:"plate"`
		Number  string `form:"number"`
		BusType string `form:"type"`
	}

	c.Bind(&busParams)

	coll := client.Database("HimWays").Collection("Bus")
	filter := bson.D{}

	if busParams.Plate != "" {
		filter = append(filter, bson.E{"plate", busParams.Plate})
	}

	if busParams.Number != "" {
		filter = append(filter, bson.E{"number", busParams.Number})
	}

	if busParams.BusType != "" {
		filter = append(filter, bson.E{"type", busParams.BusType})
	}

	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}

	var results []Bus
	if err := cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, map[string]interface{}{"data": results})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		log.Fatal("No DB MONGODB_URI found.")
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()

	r.GET("/bus", GetBus)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
