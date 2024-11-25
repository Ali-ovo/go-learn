package JWTAuth

import (
	"os"
	doUser "shop/app/shop/api/internal_api/domain/do/user"
	"shop/gmicro/pkg/errors"
	"shop/gmicro/server/restserver/middlewares"
	"shop/gmicro/server/restserver/middlewares/auth"
	"shop/pkg/code"
	"shop/pkg/options"
	"time"

	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var j *middlewares.JWT

type CustomClaims struct {
	UserID      uint `json:"userid"`
	NickName    string
	AuthorityId uint
	jwt.RegisteredClaims
}

// generateJwt 生成 JWT 结构体 (自定义的JWT结构体)
//
//	@Description:
//	@param opts
//	@return error
func generateJwt(opts *options.JwtOptions) error {
	switch opts.Method {
	case "HS256", "HS512", "HS384":
		j = middlewares.NewJWT(middlewares.WithHsSymmetricEncrypt(opts.Method, []byte(opts.Key)))
	//case "ES256", "ES512", "ES384": // 目前 go-jwt 不支持 ES 加密 所以这里就不继续写了
	//	block, _ := pem.Decode([]byte(opts.Key))
	//	privateKey, _ := x509.ParseECPrivateKey(block.Bytes)
	//	j = middlewares.NewJWT(middlewares.WithEsAsymmetricEncrypt(opts.Method, privateKey))
	case "RS256", "RS512", "RS384":
		var keyData []byte
		if opts.PrivKeyFile == "" {
			return errors.WithCode(code.ErrJWTDeploy, "invalid jwt method")
		} else {
			filecontent, err := os.ReadFile(opts.PrivKeyFile)
			if err != nil {
				return errors.WithCode(code.ErrJWTReadFiled, "读取私钥文件失败")
			}
			keyData = filecontent
		}
		// TODO privatekey 文件 也是可以加密的
		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
		if err != nil {
			return errors.WithCode(code.ErrInvalidPrivKey, "给定的私钥无效")
		}
		j = middlewares.NewJWT(middlewares.WithRsAsymmetricEncrypt(opts.Method, privateKey))
	default:
		return errors.WithCode(code.ErrJWTDeploy, "invalid jwt method")
	}
	return nil
}

// CreateJWT 生成 自定义 token
//
//	@Description:
//	@param user
//	@param opts	JWT相关配置
//	@return string
//	@return error
func CreateJWT(user *doUser.UserDO, opts *options.JwtOptions) (string, error) {
	// 生成 token
	if j == nil {
		err := generateJwt(opts)
		if err != nil {
			return "", err
		}
	}

	claims := CustomClaims{
		UserID:      uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opts.Realm,
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(opts.Timeout)), // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Local()),                   // 生效时间
			IssuedAt:  jwt.NewNumericDate(time.Now().Local()),                   // 签发时间
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {
		return "", errors.WithCode(code.ErrFailedTokenCreation, "创建 JWT 令牌失败")
	}
	return token, nil
}

// NewJWTAuth 生成 gin-jwt
//
//	@Description:
//	@param opts
//	@return middlewares.AuthStrategy
func NewJWTAuth(opts *options.JwtOptions) auth.JWTStrategy {
	/*
		gin-jwt 内部也实现了 login 和 logout 的路由的实现
		但是 本项目是 自己去实现相关接口
		这里就不用 gin-jwt 的 登录 和 注销 接口了
		但是可以借鉴它的写法
		本项目实现的是 登入 注册 注销 三个接口
	*/

	gjwt, _ := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            opts.Realm,       // 唯一性 用来区分 众多JWT 的作用
		SigningAlgorithm: opts.Method,      // JWT的 加密解密算法
		Key:              []byte(opts.Key), // 密钥
		KeyFunc:          nil,              // 设置 Key 的 回调函数 里面做一些逻辑处理
		Timeout:          opts.Timeout,     // JWT 的过期时间
		MaxRefresh:       opts.MaxRefresh,  // JWT 的最大刷新时间
		Authenticator:    nil,              // JWT 的认证器 这里实现 相关认证方法
		Authorizator:     nil,              // JWT 认证成功后 执行的回调函数
		PayloadFunc:      nil,              // 创建JWT时, 在claims中添加数据的 回调函数
		Unauthorized:     nil,              // 认证失败的处理函数
		LoginResponse:    nil,              // 自定义登录逻辑 需要在 gin router 中 添加 LoginHandler 方法
		LogoutResponse: func(c *gin.Context, code int) { // 自定义注销逻辑 需要在 gin router 中 添加 LogoutHandler 方法
			c.JSON(code, nil)
		},
		RefreshResponse:       nil,                                           // 自定义刷新逻辑 需要在 gin router 中 添加 RefreshHandler 方法
		IdentityHandler:       nil,                                           // 自定义存储到 context 中的数据的 回调函数  如 identity := mw.IdentityHandler(c); c.Set(mw.IdentityKey, identity)
		IdentityKey:           middlewares.KeyUserID,                         // 用户身份信息的键名
		TokenLookup:           "header:Authorization,query:token,cookie:jwt", // JWT 的 Token 查找方式 header:<name>/cookie:<name>/query:<name> 默认: header:Authorization
		TokenHeadName:         "Bearer",                                      // JWT 的 Token 名称 默认: Bearer
		TimeFunc:              nil,                                           // 获取当前时间的函数 对测试有用
		HTTPStatusMessageFunc: nil,                                           // JWT 中间件出现问题时候 使用
		PrivKeyFile:           opts.PrivKeyFile,                              // 非对称加密的 私钥文件
		PrivKeyBytes:          nil,                                           // 非对称加密的 私钥 字节  优先私钥文件
		PubKeyFile:            opts.PubKeyFile,                               // 非对称加密的 公钥文件
		PrivateKeyPassphrase:  "",                                            // 用于 解密 私钥文件(PS jwt目前已弃用此函数)
		PubKeyBytes:           nil,                                           // 非对称加密的 私钥 字节  优先公钥文件
	})
	return auth.NewJWTStrategy(*gjwt)
}
