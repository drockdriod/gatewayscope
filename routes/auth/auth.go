package auth

import (
	"github.com/gin-gonic/gin"
    authController "github.com/drockdriod/gatewayscope/controllers/auth"
)

func GetGroup(parentPath *gin.RouterGroup) *gin.RouterGroup {
	r := parentPath.Group("/auth")
	
	// Be sure to use struts here to define a schema in which the JSON would conform to
	
	r.POST("/register", authController.Register)

	r.GET("/accounts", authController.GetAccounts)


	return r

}