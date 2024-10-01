package api

import (
	"context"
	"fmt"
	"go-learn/ch17/user_web/forms"
	"go-learn/ch17/user_web/global"
	"go-learn/ch17/user_web/global/response"
	"go-learn/ch17/user_web/middlewares"
	"go-learn/ch17/user_web/models"
	"go-learn/ch17/user_web/proto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}

	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err

	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	// 将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}

			return
		}
	}
}

func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)

	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}

	ctx.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errs.Translate(global.Trans))})
}

func GetUserList(ctx *gin.Context) {
	// ip := "127.0.0.1"
	// port := 50051

	// 链接用户 grpc
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Error("[GetUserList] 链接 [用户服务失败]", "msg", err.Error())
	}
	defer conn.Close()

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)

	zap.S().Infof("访问用户： %d", currentUser.ID)

	// 初始化客户端
	userSrvClient := proto.NewUserClient(conn)
	pn := ctx.DefaultQuery("pn", "0")
	pnInt, _ := strconv.Atoi(pn)

	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Error("[GetUserList] 获取用户列表失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			// BirthDay: time.Time(time.Unix(int64(value.BirthDay), 0)),
			// BirthDay: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02 15:04:05"),
			BirthDay: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)

	}

	ctx.JSON(http.StatusOK, result)
}

func PassWordLogin(ctx *gin.Context) {

	// 表单验证
	passwordLoginForm := forms.PassWordLoginForm{}

	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {

		ctx.JSON(http.StatusBadRequest, map[string]string{
			"captcha": "验证码错误",
		})
		return
	}

	// 链接用户 grpc
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Error("[GetUserList] 链接 [用户服务失败]", "msg", err.Error())
	}
	defer conn.Close()

	// 初始化客户端
	userSrvClient := proto.NewUserClient(conn)

	// 登录逻辑
	if rsp, err := userSrvClient.GetUserByMobile(ctx, &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	}); err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}

			return
		}

	} else {

		// 验证密码
		if checkRsp, pasErr := userSrvClient.CheckPassWord(ctx, &proto.PasswordCheckInfo{
			Password:        passwordLoginForm.PassWord,
			EncryptPassword: rsp.PassWord,
		}); pasErr != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]string{
				"mobile": "密码错误",
			})
		} else {
			if checkRsp.Success {
				// 生成 token
				j := middlewares.NewJWT()
				claims := models.CustomClaims{
					ID:          uint(rsp.Id),
					NickName:    rsp.NickName,
					AuthorityId: uint(rsp.Role),
					StandardClaims: jwt.StandardClaims{
						NotBefore: time.Now().Unix(),               // 生效时间
						ExpiresAt: time.Now().Unix() + 60*60*24*30, // 过期时间
						Issuer:    "Ali",
					},
				}

				token, err := j.CreateToken(claims)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{
						"msg": "生成token失败",
					})
					return
				}

				ctx.JSON(http.StatusOK, gin.H{
					"id":         rsp.Id,
					"nick_name":  rsp.NickName,
					"token":      token,
					"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
				})
			} else {
				ctx.JSON(http.StatusBadRequest, map[string]string{
					"mobile": "登录失败",
				})
			}

		}

	}

}

func Register(ctx *gin.Context) {
	// 用户注册

	// 表单验证
	registerForm := forms.RegisterForm{}

	if err := ctx.ShouldBind(&registerForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 验证码校验
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err == redis.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	} else {
		if value != registerForm.Code {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"captcha": "验证码错误",
			})
			return
		}
	}

	// 链接用户 grpc
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", global.ServerConfig.UserSrvInfo.Host, global.ServerConfig.UserSrvInfo.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		zap.S().Error("[GetUserList] 链接 [用户服务失败]", "msg", err.Error())
	}
	defer conn.Close()

	// 初始化客户端
	userSrvClient := proto.NewUserClient(conn)
	user, err := userSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   registerForm.Mobile,
		PassWord: registerForm.PassWord,
		NickName: registerForm.Mobile,
	})

	if err != nil {
		zap.S().Error("[Register] 创建用户失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 生成 token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               // 生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, // 过期时间
			Issuer:    "Ali",
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})

}
