package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LicensePlate string             `json:"licensePlate" bson:"licensePlate"`
	Category     string             `json:"category" bson:"category"`
	Datetime     time.Time          `json:"datetime" bson:"datetime"`
	SpotNumber   int32              `json:"spotNumber" bson:"spotNumber"`
}
