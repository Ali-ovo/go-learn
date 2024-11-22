# 错误码

！！系统错误码列表，由 `codegen -type=int -doc` 命令生成，不要对此文件做任何更改。

## 功能说明

如果返回结果中存在 `code` 字段，则表示调用 API 接口失败。例如：

```json
{
  "code": 100101,
  "message": "Database error"
}
```

上述返回中 `code` 表示错误码，`message` 表示该错误的具体信息。每个错误同时也对应一个 HTTP 状态码，比如上述错误码对应了 HTTP 状态码 500(Internal Server Error)。

## 错误码列表

系统支持的错误码列表如下：

| Identifier | Code | HTTP Code | Description |
| ---------- | ---- | --------- | ----------- |
| ErrConnectDB | 100501 | 500 | Internal server error |
| ErrGoodsNotFound | 100701 | 500 | Internal server error |
| ErrUserNotFound | 100601 | 404 | User not found |
| ErrUserAlreadyExists | 100602 | 400 | User already exists |
| ErrUserPasswordIncorrect | 100603 | 400 | User password is incorrect |
| ErrSmsSend | 100604 | 400 | Send sms error |
| ErrJWTDeploy | 100605 | 500 | JWT deploy error |
| ErrJWTReadFiled | 100606 | 500 | JWT read field error |
| ErrInvalidPrivKey | 100607 | 500 | Internal server error |
| ErrFailedTokenCreation | 100608 | 500 | Internal server error |
| ErrCodeNotExist | 100609 | 400 | Sms code incorrect or expired |
| ErrCodeExpired | 100610 | 400 | Sms code expired |

