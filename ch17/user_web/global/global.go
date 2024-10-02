package global

import (
	"go-learn/ch17/user_web/config"
	"go-learn/ch17/user_web/proto"

	ut "github.com/go-playground/universal-translator"
)

var (
	Trans         ut.Translator
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	UserSrvClient proto.UserClient
)
