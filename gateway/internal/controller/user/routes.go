package user

import (
	"gateway/pb/userpb"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.Engine, userService userpb.UserServiceClient) {
	userGroup := r.Group("/user")
	{
		userGroup.GET("/:id", func(c *gin.Context) {
			GetUser(c, userService)
		})
	}

}
