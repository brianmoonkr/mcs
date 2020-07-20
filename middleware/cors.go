package middleware

import (
	"net/http"

	"github.com/kataras/iris"
)

const (
	options          string = "OPTIONS"
	allowOrigin      string = "Access-Control-Allow-Origin"
	allowMethods     string = "Access-Control-Allow-Methods"
	allowHeaders     string = "Access-Control-Allow-Headers"
	allowCredentials string = "Access-Control-Allow-Credentials"
	exposeHeaders    string = "Access-Control-Expose-Headers"
	credentials      string = "true"
	origin           string = "Origin"
	methods          string = "POST, GET, OPTIONS, PUT, DELETE, HEAD, PATCH"

	headers string = "Access-Control-Allow-Origin, Accept, Accept-Encoding, Authorization, Content-Length, Content-Type, X-CSRF-Token, X-Requested-With, KEY"
)

// Cors ...
func Cors(ctx iris.Context) {
	ctx.Header(allowOrigin, "*")
	ctx.Header(allowHeaders, headers)
	ctx.Header(allowMethods, methods)
//	ctx.Header(allowCredentials, credentials)
	ctx.Header(exposeHeaders, headers)

	if ctx.Method() == options {
		ctx.ResponseWriter().WriteHeader(http.StatusOK)
		return
	}
	ctx.Next()
}
