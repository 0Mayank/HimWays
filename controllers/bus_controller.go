package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/0Mayank/himways/configs"
	"github.com/0Mayank/himways/models"
	"github.com/0Mayank/himways/responses"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	busCollection *mongo.Collection = configs.GetCollection(configs.DB, "Bus")
	validate                        = validator.New()
)

func CreateBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var bus models.Bus
		defer cancel()

		// validate request body
		if err := c.BindJSON(&bus); err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.BusResponse{
					Status:  http.StatusBadRequest,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		// validate required fields
		if validationErr := validate.Struct(&bus); validationErr != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.BusResponse{
					Status:  http.StatusBadRequest,
					Message: "validation error",
					Data:    map[string]interface{}{"data": validationErr.Error()},
				},
			)
		}

		newBus := models.Bus{
			Id:     primitive.NewObjectID(),
			Number: bus.Number,
			Plate:  bus.Plate,
			Route:  bus.Route,
		}

		result, err := busCollection.InsertOne(ctx, newBus)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		c.JSON(
			http.StatusCreated,
			responses.BusResponse{
				Status:  http.StatusCreated,
				Message: "success",
				Data:    map[string]interface{}{"data": result},
			},
		)
	}
}

func SearchBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var buses []models.Bus

		var busParams struct {
			Plate   string `form:"plate"`
			Number  string `form:"number"`
			BusType string `form:"type"`
		}

		c.Bind(&busParams)

		filter := bson.M{}

		if busParams.Plate != "" {
			filter["plate"] = busParams.Plate
		}

		if busParams.Number != "" {
			filter["number"] = busParams.Number
		}

		if busParams.BusType != "" {
			filter["type"] = busParams.BusType
		}

		cursor, err := busCollection.Find(context.TODO(), filter)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		if err := cursor.All(ctx, &buses); err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.BusResponse{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"data": buses},
			},
		)
	}
}
