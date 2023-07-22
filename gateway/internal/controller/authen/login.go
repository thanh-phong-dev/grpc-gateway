package authen

import (
	"fmt"
	"gateway/pb/authpb"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type loginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context, authService authpb.AuthenticationServiceClient) {
	var loginRequest loginRequest
	err := ctx.ShouldBindBodyWith(&loginRequest, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("%+v", err),
		})
		return
	}

	loginRequestPb := authpb.LoginRequest{
		Username: loginRequest.UserName,
		Password: loginRequest.Password,
	}
	loginResponse, err := authService.Login(ctx, &loginRequestPb)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err_message": fmt.Sprintf("%+v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}
