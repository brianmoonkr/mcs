package middleware

import (
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
	"github.com/teamgrit-lab/cojam/component/session"
)

// CheckSession ...
func CheckSession(ctx iris.Context) {
	//tglog.Logger.Info("[middleware] CheckSession !!!")

	// client 브라우져에서 쿠키값 확인
	cookieStr := ctx.GetCookie(session.LOGIN_COOKIE_NAME)
	if len(cookieStr) == 0 {
		ctx.Values().Set("IsLogin", false)
		ctx.Next()
		return
	}

	// cookie 가 있으면 redisSession 에서 세션 정보를 정확히 확인.
	if !session.IsExistence(cookieStr) {
		ctx.Values().Set("IsLogin", false)
		ctx.Next()
		return
	}

	// 세션정보를 레디스세션디비에서 가져온다.
	sessionInfo, err := session.GetSession(cookieStr)
	if err != nil {
		ctx.Values().Set("IsLogin", false)
		ctx.Next()
		return
	}

	// 레디스세션디비의 세션시간을 초기화
	session.SetExpireTime(cookieStr)

	// 로그인 되어 있으면 context data 에 로그인 상태를 셋팅한다.
	ctx.Values().Set("IsLogin", true)
	// 세션 정보 셋팅

	ctx.Values().Set(ctxkey.SessionKey, cookieStr)
	ctx.Values().Set(ctxkey.SessionInfo, sessionInfo)

	ctx.Next()
}
