package clients

import (
	"github.com/gin-gonic/gin"
	// "log"
    // "github.com/drockdriod/gatewayscope/utils/crypto/jwt"
    commonUtil "github.com/drockdriod/gatewayscope/utils/common"
    clientController "github.com/drockdriod/gatewayscope/controllers/clients"
)

func GetGroup(parentPath *gin.RouterGroup) *gin.RouterGroup {
	r := parentPath.Group("/clients/:clientId")

	r.Use(commonUtil.ClientAuthMiddleware())
	{
		r.POST("/permissions/check", clientController.ComparePermissionsByUser)
	}

	r.Use(commonUtil.SetContextValue("AUTHORIZER_TYPE", "CLIENT"), commonUtil.ClientAuthMiddleware())
	{
		r.GET("/users", clientController.GetUsers)
	}

	// Permissions route:
	// send JWT, with user id in route as a param
	// 
	
	

	return r

}

