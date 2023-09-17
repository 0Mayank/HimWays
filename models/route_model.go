package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Route struct {
	Number   string               `bson:"number" json:"number" validate:"required"`
	BusRoute []primitive.ObjectID `bson:"route"  json:"route"  validate:"required"`
}
