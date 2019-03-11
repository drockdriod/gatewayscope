package clients

import (
	"github.com/gin-gonic/gin"
	// "log"
    // "github.com/drockdriod/gatewayscope/utils/crypto/jwt"
    "github.com/drockdriod/gatewayscope/utils/common"
    clientController "github.com/drockdriod/gatewayscope/controllers/clients"
)

func GetGroup(parentPath *gin.RouterGroup) *gin.RouterGroup {
	r := parentPath.Group("/clients/:clientId")

	r.Use(common.ClientAuthMiddleware())
	{
		r.POST("/permissions/:userId/check", clientController.ComparePermissionsByUserId)
	}



	// Permissions route:
	// send JWT, with user id in route as a param
	// 
	
	

	return r

}