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
    "github.com/drockdriod/gatewayscope/utils/crypto"
)

var client, dbContext, err = db.Client()


func GetAccounts(c *gin.Context){
	accounts := db.GetItems("accounts", bson.D{})

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,	
	})
}


func Register(c *gin.Context) {

	var jsonBody commonModels.Register

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }

	jsonBody.HashPassword = crypto.HashAndSalt(jsonBody.Password)
	log.Println(jsonBody)



	res, err := db.InsertObj("accounts", jsonBody)

	log.Println(res)

	if err != nil { 
		c.AbortWithError(400, err)
		return 
	}




	c.JSON(http.StatusOK, gin.H{
		"message": "Account registered",
	})

}