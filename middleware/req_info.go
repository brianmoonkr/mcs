package middleware

import (
	"encoding/base64"

	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/constdf"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
	"github.com/teamgrit-lab/cojam/component/tglog"
	"github.com/teamgrit-lab/cojam/component/util"
)

// ReqInfo ...
func ReqInfo(ctx iris.Context) {
	var requestID string

	if ctx.Values().Get("IsLogin").(bool) {
		sessionKeyBase64 := ctx.Values().Get(ctxkey.SessionKey).(string)
		sessionKey, err := base64.StdEncoding.DecodeString(sessionKeyBase64)
		if err != nil {
			tglog.Logger.Errorf("ReqInfo - sessionKey : %s\n", string(sessionKey))
		}

		requestID = string(sessionKey)
	} else {
		requestID = util.MakeUniqueID()
	}

	tglog.Logger.Infof("User-Agent : %+v\n", ctx.GetHeader("User-Agent"))

	isMobile := ctx.IsMobile()
	userAgent := constdf.GLOBAL_SESSION_ID_USERAGENT_WEB

	tglog.Logger.Infof("RequestID: %s, IsMobile : %t, method: %s, path: %s\n",
		requestID,
		isMobile,
		ctx.Method(),
		ctx.Path())

	if isMobile {
		userAgent = constdf.GLOBAL_SESSION_ID_USERAGENT_APP
	}

	ctx.Values().Set(ctxkey.RequestID, requestID)
	ctx.Values().Set(ctxkey.UserAgent, userAgent)

	ctx.Next()
}
