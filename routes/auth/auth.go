package auth

import (
	"github.com/gin-gonic/gin"
	// "log"
    // "github.com/drockdriod/gatewayscope/utils/crypto/jwt"
    "github.com/drockdriod/gatewayscope/utils/common"
    authController "github.com/drockdriod/gatewayscope/controllers/auth"
)

func GetGroup(parentPath *gin.RouterGroup) *gin.RouterGroup {
	r := parentPath.Group("/auth/:clientId")
		
	parentPath.POST("/register", authController.Register)

	parentPath.POST("/login", authController.Login)

	r.Use(common.ClientAuthMiddleware())
	{
		r.GET("/accounts", authController.GetAccounts)
		r.POST("/users/register", authController.UserRegister)
	}



	// Permissions route:
	// send JWT, with user id in route as a param
	// 
	
	

	return r

}