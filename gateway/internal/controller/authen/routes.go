package authen

import (
	"gateway/pb/authpb"
	"gateway/pb/userpb"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine, authService authpb.AuthenticationServiceClient, userService userpb.UserServiceClient) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", func(c *gin.Context) {
			Login(c, authService)
		})
		authGroup.POST("/register", func(c *gin.Context) {
			Register(c, userService)
		})
	}
}
