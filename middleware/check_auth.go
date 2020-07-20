package middleware

import (
	"strings"

	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
)

// CheckAuth ...
func CheckAuth(ctx iris.Context) {
	path := ctx.Path()

	// 로그인 없이 이용가능한 페이지 체크
	if !isLoginAuthURL(path) {
		ctx.Next()
		return
	}

	if !ctx.Values().Get("IsLogin").(bool) {
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
		return
	}

	userSession := ctx.Values().Get(ctxkey.SessionInfo)
	if userSession == nil {
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
		return
	}

	ctx.Next()
	return
}

// isLoginAuthURL 는 로그인해야 진입할 수 있는 페이지체크.
// URL 로 체크한다.
func isLoginAuthURL(path string) bool {
	authURLs := []string{
		"/live/core",
	}

	for _, v := range authURLs {
		if strings.Contains(path, v) {
			return true
		}
	}
	return false
}
