package controllerUser

import (
	"fmt"
	"shop/gmicro/pkg/common/core"
	"shop/gmicro/pkg/common/time"
	"shop/gmicro/server/restserver/middlewares"
	gin2 "shop/pkg/translator/gin"

	"github.com/gin-gonic/gin"
)

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
}

// UpdateUser 更新用户信息
//
//	@Description:
//	@receiver us
//	@param ctx
func (us *userController) UpdateUser(ctx *gin.Context) {
	updateForm := UpdateUserForm{}
	if err := ctx.ShouldBind(&updateForm); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	userID, _ := ctx.Get(middlewares.KeyUserID)
	userIDInt := uint64(userID.(float64))
	userDTO, err := us.srv.Get(ctx, userIDInt)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	userDTO.NickName = updateForm.Name
	userDTO.Gender = updateForm.Gender
	userDTO.Birthday, err = time.ToTime(fmt.Sprint(updateForm.Birthday + " 00:00:00"))

	err = us.srv.Update(ctx, userDTO)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)
}
