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
    "github.com/drockdriod/gatewayscope/utils/crypto/keygen"
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
		log.Println(err.Error())
		log.Fatal("error jwt")
	}

	err = commonModels.RegenerateClient(id, tokenString)

	if err != nil {
		log.Println("error regen")

		log.Fatal(err.Error())
	}

	return tokenString
}

func Login(c *gin.Context) {

	var jsonBody commonModels.Login
	var items interface{}
	var account commonModels.Account

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }


    items = db.FindOne("accounts", bson.M{
    	"email": jsonBody.Email,
    })

	body1, err := bson.Marshal(items)

	if err != nil {
		log.Println(err.Error())
	}

	bson.Unmarshal(body1, &account)

	log.Println(account)

	allowAccess := bcrypt.ComparePasswords(account.HashPassword, []byte(jsonBody.Password))

	if allowAccess {

		tokenString := apiTokenGenerate(account.Id)

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"clientId": account.Id.Hex(),
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
	}

}


func Register(c *gin.Context) {

	var jsonBody commonModels.Account

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }

	jsonBody.HashPassword = bcrypt.HashAndSalt(jsonBody.Password)

	res, err, clientId := db.InsertObj("accounts", jsonBody)

	log.Println(res)


	if err != nil { 
		c.AbortWithError(400, err)
		return 
	}

	err = keygen.GenerateKeysForClient(clientId.Hex())

	if err != nil {
		c.AbortWithError(400, err)
        return
	}
	
	tokenString := apiTokenGenerate(clientId)


	c.JSON(http.StatusOK, gin.H{
		"message": "Account registered",
		"token": tokenString,
		"clientId": clientId,
	})

}

func UserRegister(c *gin.Context) {
	var jsonBody commonModels.Account
	clientId := c.Query("clientId")

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }

	jsonBody.HashPassword = bcrypt.HashAndSalt(jsonBody.Password)

	res, err, userId := db.InsertObj("users", jsonBody)

	log.Println(res)
	

	tokenString, err := jwt.GenerateToken(clientId, userId.Hex())

	if(err != nil){
		log.Println(err.Error())
		log.Fatal("error jwt")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Registration complete",
		"token": tokenString,
	})

}



