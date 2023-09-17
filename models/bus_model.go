package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Bus struct {
	Id      primitive.ObjectID `bson:"_id"    json:"_id"`
	Plate   string             `bson:"plate"  json:"plate"  validate:"required"`
	Number  string             `bson:"number" json:"number" validate:"required"`
	BusType string             `bson:"type"   json:"type"   validate:"required"`
	Route   []int              `bson:"route"  json:"route"  validate:"required"`
}
