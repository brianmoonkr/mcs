package middleware

import (
	"github.com/kataras/iris"
)

// APIURLAuth 는 API 전용 미들웨어이다.
// 로그인된 사용자만 허용해야하는 API를 필터한다.
func BasicAuth(ctx iris.Context) {
	username, password, _ := ctx.Request().BasicAuth()
	//key := ctx.GetHeader("key")
	ctx.Values().Set("AuthId", username)
	ctx.Values().Set("AuthPwd", password)
	//ctx.Values().Set("key", key)

	ctx.Next()
}


