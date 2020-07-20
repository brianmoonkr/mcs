package middleware

import (
	"fmt"
	"strings"

	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
	"github.com/teamgrit-lab/cojam/component/session"
	"github.com/teamgrit-lab/cojam/mvc/vo"
)

// APIURLAuth 는 API 전용 미들웨어이다.
// 로그인된 사용자만 허용해야하는 API를 필터한다.
func APIURLAuth(ctx iris.Context) {
	path := ctx.Path()

	var err error
	var cookieStr string
	var sessionInfo *session.UserSession

	// Request 쿠키값 확인
	cookieStr = ctx.GetCookie(session.LOGIN_COOKIE_NAME)

	if len(cookieStr) != 0 {
		if isAPIURLAuth(path) {
			// 세션정보를 레디스세션에서 가져온다.
			sessionInfo, err = session.GetSession(cookieStr)
			if err != nil {
				sendError(ctx)
				return
			}

			// Update : 세션 ExpireTime
			session.SetExpireTime(cookieStr)
		}

		// 세션 정보 셋팅
		ctx.Values().Set(ctxkey.SessionKey, cookieStr)
		ctx.Values().Set(ctxkey.SessionInfo, sessionInfo)
		ctx.Values().Set("IsLogin", true)
		ctx.Next()
		return
	}

	ctx.Values().Set("IsLogin", false)
	ctx.Next()
}

// isAPIURLAuth 는 로그인된 사용자만 진입가능한
// URL 체크 함수이다.
func isAPIURLAuth(path string) bool {
	authURLs := []string{
		"/api/v3/create",
		"/api/v2/create",
		"/api/v1/live/registration",
		"/api/v1/live/start",
		"/api/v1/live/end",
		"/api/v1/live/join",
		"/api/v1/live/room/make",
		"/api/v1/custcenter/faq",
		"/api/v1/channel/subscription",
		"/api/v1/vod/comment",
		"/api/v1/vod/like",
		"/api/v1/vod/like/check",
		"/api/v1/user/notice/agree",
		"/api/v1/user/nicknm",
		"/api/v1/user/withdrawal",
	}

	for _, v := range authURLs {
		if strings.Contains(path, v) {
			return true
		}
	}
	return false
}

func sendError(ctx iris.Context) {
	res := new(vo.ResponseVO)
	res.Error400(ctx, fmt.Errorf("No Session"), "middleware - APIURLAuth()")
}
