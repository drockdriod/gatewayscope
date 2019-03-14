package clients

import (
	commonUtils "github.com/drockdriod/gatewayscope/utils/common"
	"github.com/drockdriod/gatewayscope/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type RouteContext struct{
	Permissions []string `"json": "permissions"`
}


func GetUsers(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{
		"users": "Access granted",
	})
}

func ComparePermissionsByUser(c *gin.Context){
	user := c.MustGet("USER").(models.User)

	log.Println("user")
	log.Println(user.Account)

	var jsonBody RouteContext

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }

    /**
     * get this from DB
     */
	userPermissions := []string{"write:all"}

	granted := commonUtils.CompareUserPermissions(jsonBody.Permissions, userPermissions)

	if granted == false {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Access denied",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Access granted",
	})
}