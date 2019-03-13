package models

import (
	bsonPrimitive "github.com/mongodb/mongo-go-driver/bson/primitive"
	commonModels "github.com/drockdriod/gatewayscope/models/common"
)

// Binding from JSON
type User struct {
	Id		 	bsonPrimitive.ObjectID `bson:"_id" json:"-"`
	Account		commonModels.Account
	Meta	 	interface{}	`bson:"meta" json:"meta,omitempty"`
	ClientId	bsonPrimitive.ObjectID `bson:"clientId" json:"-"`
}