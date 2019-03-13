package common

import (
	bsonPrimitive "github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/bson"
    "github.com/drockdriod/gatewayscope/db"
)

// Binding from JSON
type Account struct {
	Id		 bsonPrimitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Name     string `bson:"name" json:"name"`
	Email 	 string `bson:"email" json:"email"`
	Password string `json:"password" bson:"-"`
	HashPassword string `bson:"hashpassword"`
}

/**
 * Updates the account with a new token and blacklists all old tokens
 */
func RegenerateClient(id bsonPrimitive.ObjectID, tokenString string) error{
	_, err := db.UpdateObj("accounts", bson.M{"_id": id}, bson.D{
			{"$set", bson.M{ 
				"token": tokenString,
			},
		} } )

	if err != nil{
		return err
	}

	_, err = db.UpdateObj("api_tokens", bson.D{{"accountId", id}, {"blacklisted", false}}, bson.D{
			{"$set", bson.M{ 
				"blacklisted": true,
			},
		} } )


	_, err, _ = db.InsertObj("api_tokens", bson.M{
		"token": tokenString,
		"blacklisted": false,
		"accountId": id,
	})
	
	if err != nil{
		return err
	}

	return nil
}