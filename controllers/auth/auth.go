package auth


import (
	"github.com/gin-gonic/gin"
    "github.com/drockdriod/gatewayscope/db"
    commonModels "github.com/drockdriod/gatewayscope/models/common"
    "log"
    // "encoding/json"
    // "context"
    // "fmt"
    "net/http"
    "github.com/mongodb/mongo-go-driver/bson"
	bsonPrimitive "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/drockdriod/gatewayscope/utils/crypto/bcrypt"
    "github.com/drockdriod/gatewayscope/utils/crypto/jwt"
)

func GetAccounts(c *gin.Context){
	accounts := db.GetItems("accounts", bson.D{})

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,	
	})
}

func apiTokenGenerate(id bsonPrimitive.ObjectID) string{
	tokenString, err := jwt.GenerateToken(id.Hex())

	if(err != nil){
		log.Fatal("error jwt")
		log.Println(err)
	}

	db.UpdateObj("accounts", bson.M{"_id": id}, bson.M{"token": tokenString})

	db.InsertObj("api_tokens", bson.M{
		"token": tokenString,
		"blacklisted": false,
		"accountId": id,
	})

	return tokenString
}


func Register(c *gin.Context) {

	var jsonBody commonModels.Register

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }

	jsonBody.HashPassword = bcrypt.HashAndSalt(jsonBody.Password)
	log.Println(jsonBody)

	res, err, objectId := db.InsertObj("accounts", jsonBody)

	log.Println(res)


	if err != nil { 
		c.AbortWithError(400, err)
		return 
	}
	
	tokenString := apiTokenGenerate(objectId)


	c.JSON(http.StatusOK, gin.H{
		"message": "Account registered",
		"token": tokenString,
	})

}




