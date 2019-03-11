package common

// Binding from JSON
type Login struct {
	Email 	 string `bson:"email" json:"email"`
	Password string `json:"password" bson:"-"`
}