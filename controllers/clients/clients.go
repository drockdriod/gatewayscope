package clients

import (
	"github.com/drockdriod/gatewayscope/utils/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"log"
)

type RouteContext struct{
	Permissions []string `"json": "permissions"`
}

func ComparePermissionsByUserId(c *gin.Context){
	userId := c.Query("userId")

	var jsonBody RouteContext

	err := c.BindJSON(&jsonBody)

	if err != nil {
        c.AbortWithError(400, err)
        return
    }

    log.Println(userId)

    /**
     * get this from DB
     */
	userPermissions := []string{"write:all"}

	granted := common.CompareUserPermissions(jsonBody.Permissions, userPermissions)

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