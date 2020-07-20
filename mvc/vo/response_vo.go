package vo

import (
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/tglog"
)

// ResponseVO ...
type ResponseVO struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// Send 는 오류가 없는 정상적인 response를 보낸다.
func (r *ResponseVO) Send(ctx iris.Context, msg string) {
	r.Status = iris.StatusOK
	r.Message = msg
	ctx.JSON(r)
}

// Send422 는 오류가 없는 정상적인 response를 보낸다.
// 422: UNPROCESSABLE ENTITY
// 요청을 처리할 수 없음.
// 올바른 요청일 수 있으나, 해당 작업에는 유효하지 않다.
func (r *ResponseVO) Send422(ctx iris.Context, msg string) {
	r.Status = iris.StatusUnprocessableEntity
	r.Message = msg
	ctx.JSON(r)
}

// Error400 은 요청을 처리할 수 없다.
func (r *ResponseVO) Error400(ctx iris.Context, err error, spot string) {
	ctx.Values().Set("RDBRollBack", true)
	errMsg := err.Error()
	tglog.PrintErr(ctx, spot, errMsg)
	r.Status = iris.StatusBadRequest
	r.Message = "잠시 서비스에 죄송한 상황이 발생했습니다."
	ctx.JSON(r)
}

func (r *ResponseVO) ErrorMessage(ctx iris.Context, status int, spot string) {
	ctx.Values().Set("RDBRollBack", true)
	//tglog.PrintErr(ctx, spot)
	r.Status = status
	r.Message = spot
	ctx.JSON(r)
}
