package auth

import (
	"github.com/gin-gonic/gin"
    ginJWT "github.com/appleboy/gin-jwt"
	// "log"
    commonUtils "github.com/drockdriod/gatewayscope/utils/common"
    authController "github.com/drockdriod/gatewayscope/controllers/auth"
)

func GetGroup(parentPath *gin.RouterGroup) *gin.RouterGroup {
	r := parentPath.Group("/auth/:clientId")
		
	parentPath.POST("/register", authController.ClientRegister)

	parentPath.POST("/login", authController.ClientLogin)

	r.Use(commonUtils.SetContextValue("AUTHORIZER_TYPE", "CLIENT"),commonUtils.ClientAuthMiddleware())
	{
		r.GET("/accounts", authController.GetAccounts)

		usersRoute := r.Group("/users")

		usersRoute.POST("/register", authController.UserRegister)
		usersRoute.POST("/login", func(c *gin.Context) {
			authMiddleware := c.MustGet("authMiddleware").(*ginJWT.GinJWTMiddleware)
			authMiddleware.LoginHandler(c)
		})
	}



	// Permissions route:
	// send JWT, with user id in route as a param
	// 
	
	

	return r

}