package middleware

import (
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/constdf"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
	"github.com/teamgrit-lab/cojam/component/session"
)

// AdminAuth ...
func AdminAuth(ctx iris.Context) {
	// *step - Request 쿠키값 확인
	cookieStr := ctx.GetCookie(session.LOGIN_COOKIE_NAME)
	if len(cookieStr) == 0 {
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
		return
	}

	// *step - 세션정보를 레디스세션에서 가져온다.
	sessionInfo, err := session.GetSession(cookieStr)
	if err != nil {
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
		return
	}
	// Update : 세션 ExpireTime
	session.SetExpireTime(cookieStr)

	// *step - admin 권한확인
	if !CheckRoles(sessionInfo) {
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
		return
	}

	// 세션 정보 셋팅
	ctx.Values().Set(ctxkey.SessionKey, cookieStr)
	ctx.Values().Set(ctxkey.SessionInfo, sessionInfo)
	ctx.Values().Set("IsLogin", true)
	ctx.Next()
}

// CheckRoles ...
func CheckRoles(sessionInfo *session.UserSession) bool {
	for _, v := range sessionInfo.Roles {
		if v == constdf.USER_AUTH_CODE_ADMIN {
			return true
		}
	}
	return false
}
