package common


import (
	"github.com/gin-gonic/gin"
    "github.com/drockdriod/gatewayscope/utils/crypto/jwt"
	"log"
)


func indexOf(element string, data []string ) (int) {
   for k, v := range data {
       if element == v {
           return k
       }
   }
   return -1    //not found.
}


func ClientAuthMiddleware() (gin.HandlerFunc){
	return func(c *gin.Context) {

		clientId := c.Param("clientId")

		authMiddleware, err := jwt.InitClientAuthMiddleware(clientId)

		c.Set("authMiddleware", authMiddleware)

		log.Println(authMiddleware)

		if( err != nil){
			log.Println(err.Error())
		}

		authMiddleware.MiddlewareFuncAlt(c)
	}
}

func CompareUserPermissions(permissions []string, userPermissions []string) bool {
	granted := false

	for n := 0; n < len(userPermissions); n++ {
		if indexOf(userPermissions[n],permissions) != -1 {
			granted = true
		}
	}

	return granted
}

func SetContextValue(property string, value interface{}) gin.HandlerFunc{
	return func(c *gin.Context){
		c.Set(property, value)	
	}
}



