package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bus struct {
	Id      primitive.ObjectID `json:"_id"`
	Plate   string             `json:"plate"  validate:"required"`
	Number  string             `json:"number" validate:"required"`
	BusType string             `json:"type"   validate:"required"`
	Route   []int              `json:"route"  validate:"required"`
}
