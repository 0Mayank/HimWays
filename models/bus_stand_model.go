package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BusStand struct {
	Id        primitive.ObjectID `bson:"_id"       json:"_id"`
	Name      string             `bson:"name"      json:"name"      validate:"required"`
	Latitude  string             `bson:"latitude"  json:"latitude"  validate:"required"`
	Longitude string             `bson:"longitude" json:"longitude" validate:"required"`
}
