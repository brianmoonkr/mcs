package middleware

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/tglog"
	"github.com/teamgrit-lab/cojam/config"
)

// Recover ..
func Recover(ctx iris.Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.StatusCode(500)
			CommonRecover(ctx, err)
		}
	}()
	ctx.Next()
}

// CommonRecover ...
func CommonRecover(ctx iris.Context, err interface{}) {
	if ctx.IsStopped() {
		return
	}

	var stacktrace string
	pjName := config.CF.Prop.ProjectName
	for i := 1; ; i++ {
		_, f, l, got := runtime.Caller(i)
		if !got {
			break
		}

		if strings.Contains(f, pjName) {
			stacktrace += fmt.Sprintf("\t%s:%d\n", strings.Split(f, pjName)[1], l)
		}
	}

	// when stack finishes
	logMessage := fmt.Sprintf("\n\n============================== ERROR LOG START ================================")
	logMessage += fmt.Sprintf("\nAt Request: %s\n", fmt.Sprintf("%d %s %s %s", ctx.GetStatusCode(), ctx.Path(), ctx.Method(), ctx.RemoteAddr()))
	logMessage += fmt.Sprintf("ERROR Message : %s\n", err)
	logMessage += fmt.Sprintf("\n%s\n", stacktrace)
	logMessage += fmt.Sprintf("\n================================ ERROR LOG END ================================\n\n")
	tglog.Logger.Error(logMessage)
	ctx.StopExecution()
}
