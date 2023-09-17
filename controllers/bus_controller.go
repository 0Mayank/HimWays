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
	busCollection   *mongo.Collection = configs.GetCollection(configs.DB, "Bus")
	validate                          = validator.New()
	timeoutDuration                   = 10 * time.Second
)

func CreateBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)

		var bus models.Bus
		defer cancel()

		// validate request body
		if err := c.BindJSON(&bus); err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.BusResponse{
					Status:  http.StatusBadRequest,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
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
					Body:    map[string]interface{}{"data": validationErr.Error()},
				},
			)
			return
		}

		newBus := models.Bus{
			Id:      primitive.NewObjectID(),
			Number:  bus.Number,
			Plate:   bus.Plate,
			BusType: bus.BusType,
		}

		result, err := busCollection.InsertOne(ctx, newBus)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		c.JSON(
			http.StatusCreated,
			responses.BusResponse{
				Status:  http.StatusCreated,
				Message: "success",
				Body:    map[string]interface{}{"data": result},
			},
		)
	}
}

// TODO: search for bus which pass through this route??
func SearchBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
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

		cursor, err := busCollection.Find(ctx, filter)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
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
					Body:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.BusResponse{
				Status:  http.StatusOK,
				Message: "success",
				Body:    map[string]interface{}{"data": buses},
			},
		)
	}
}

func EditBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		busId := c.Param("busId")
		var bus models.Bus

		objId, _ := primitive.ObjectIDFromHex(busId)

		// validate the request body
		if err := c.BindJSON(&bus); err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.BusResponse{
					Status:  http.StatusBadRequest,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		update := bson.M{}

		if bus.Plate != "" {
			update["plate"] = bus.Plate
		}

		if bus.Number != "" {
			update["number"] = bus.Number
		}

		if bus.BusType != "" {
			update["type"] = bus.BusType
		}

		result, err := busCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		var updatedBus models.Bus
		if result.MatchedCount == 1 {
			err := busCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedBus)
			if err != nil {
				c.JSON(
					http.StatusInternalServerError,
					responses.BusResponse{
						Status:  http.StatusInternalServerError,
						Message: "error",
						Body:    map[string]interface{}{"data": err.Error()},
					},
				)
				return
			}
		}

		c.JSON(
			http.StatusOK,
			responses.BusResponse{
				Status:  http.StatusOK,
				Message: "success",
				Body:    map[string]interface{}{"data": updatedBus},
			},
		)
	}
}

func GetBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		busId := c.Param("busId")
		var bus models.Bus

		objId, _ := primitive.ObjectIDFromHex(busId)

		err := busCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&bus)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.BusResponse{
				Status:  http.StatusOK,
				Message: "success",
				Body:    map[string]interface{}{"data": bus},
			},
		)
	}
}

func DeleteBus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
		defer cancel()

		busId := c.Param("busId")
		objId, _ := primitive.ObjectIDFromHex(busId)

		result, err := busCollection.DeleteOne(ctx, bson.M{"_id": objId})
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.BusResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Body:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(
				http.StatusNotFound,
				responses.BusResponse{
					Status:  http.StatusNotFound,
					Message: "error",
					Body:    map[string]interface{}{"data": "Bus with specified id not found!"},
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.BusResponse{
				Status:  http.StatusOK,
				Message: "success",
				Body:    map[string]interface{}{"data": "Bus deleted successfully!"},
			},
		)
	}
}
