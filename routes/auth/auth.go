package auth

import (
	"github.com/gin-gonic/gin"
    authController "github.com/drockdriod/gatewayscope/controllers/auth"
)

func GetGroup(parentPath *gin.RouterGroup) *gin.RouterGroup {
	r := parentPath.Group("/auth")

	// r.use()
		
	r.POST("/register", authController.Register)

	r.GET("/accounts", authController.GetAccounts)


	return r

}