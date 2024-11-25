package controllerUser

import (
	"shop/gmicro/pkg/common/core"
	"shop/gmicro/server/restserver/middlewares"

	"github.com/gin-gonic/gin"
)

func (us *userController) GetUserDetail(ctx *gin.Context) {
	userID, _ := ctx.Get(middlewares.KeyUserID)
	userDTO, err := us.srv.Get(ctx, uint64(userID.(float64)))
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, gin.H{
		"name":     userDTO.NickName,
		"birthday": userDTO.Birthday.Format("2006-01-02"),
		"gender":   userDTO.Gender,
		"moble":    userDTO.Mobile,
	})
}
