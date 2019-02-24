package common

// Binding from JSON
type Register struct {
	Name     string `bson:"name" json:"name"`
	Email 	 string `bson:"email" json:"email"`
	Password string `json:"password" bson:"-"`
	HashPassword string `bson:"hashpassword"`
}