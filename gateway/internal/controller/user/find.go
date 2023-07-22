package user

import (
	"fmt"
	"gateway/pb/userpb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUser(ctx *gin.Context, userService userpb.UserServiceClient) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid id",
		})
		return
	}

	userInfoRequest := userpb.UserInfoRequest{
		Id: int64(id),
	}
	userInfoResponse, err := userService.GetUserInfo(ctx, &userInfoRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"err_message": fmt.Sprintf("%+v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, userInfoResponse)
}
