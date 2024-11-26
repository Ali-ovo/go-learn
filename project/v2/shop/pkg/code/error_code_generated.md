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
| ErrEsDatabase | 100601 | 404 | EsDatabase error |
| ErrEsUnmarshal | 100602 | 500 | Es unmarshal error |
| ErrConnectGRPC | 100701 | 500 | Connect to grpc error |
| ErrRedisDatabase | 100801 | 500 | Redis data base error |
| ErrGoodsNotFound | 101101 | 404 | Goods not found |
| ErrCategoryNotFound | 101102 | 404 | Category not found |
| ErrBrandsNotFound | 101103 | 404 | Brand not found |
| ErrBannerNotFound | 101104 | 404 | Banner not found |
| ErrCategoryBrandsNotFound | 101105 | 404 | CategoryBrands not found |
| ErrInventoryNotFound | 101201 | 404 | Goods not found |
| ErrInvSellDetailNotFound | 101202 | 400 | Inventory sell detail not found |
| ErrInvNotEnough | 101203 | 404 | Inventory not enough |
| ErrstockNotFound | 101204 | 404 | Order not found |
| ErrOrderNotFound | 101301 | 404 | Order not found |
| ErrShopCartNotFound | 101302 | 404 | ShopCart not found |
| ErrOrderDtm | 101303 | 404 | Dtm unknonwn error |
| ErrNotGoodsSelect | 101304 | 404 | No Goods selected |
| ErrUserNotFound | 101001 | 404 | User not found |
| ErrUserAlreadyExists | 101002 | 400 | User already exists |
| ErrUserPasswordIncorrect | 101003 | 400 | User password is incorrect |
| ErrSmsSend | 101004 | 400 | Send sms error |
| ErrJWTDeploy | 101005 | 500 | JWT deploy error |
| ErrJWTReadFiled | 101006 | 500 | JWT read field error |
| ErrInvalidPrivKey | 101007 | 500 | Internal server error |
| ErrFailedTokenCreation | 101008 | 500 | Internal server error |
| ErrCodeNotExist | 101009 | 400 | Sms code incorrect or expired |
| ErrCodeExpired | 101010 | 400 | Sms code expired |

