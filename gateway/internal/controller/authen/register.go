package authen

import (
	"fmt"
	"gateway/pb/userpb"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Email    string `json:"email"`
}

func Register(ctx *gin.Context, userService userpb.UserServiceClient) {
	var registerRequest RegisterRequest
	err := ctx.ShouldBindBodyWith(&registerRequest, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("%+v", err),
		})
		return
	}

	registerRequestPb := userpb.RegisterRequest{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Name:     registerRequest.Name,
		Gender:   registerRequest.Gender,
		Email:    registerRequest.Email,
	}

	registerResponse, err := userService.Register(ctx, &registerRequestPb)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err_message": fmt.Sprintf("%+v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": registerResponse.Status,
	})
}
