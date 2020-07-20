package tglog // import "github.com/teamgrit-lab/cojam/component/tglog"

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
)

var fileLog os.File

// Logger ..
var Logger *golog.Logger

// NewLogFile ...
func NewLogFile() {
	Logger = golog.New()

	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(err)
	}

	go func() {
		var hour int
		for {
			t := time.Now().In(loc)
			nowHour := t.Hour()
			if nowHour != hour {
				hour = nowHour
				fileLog.Close()
				createFile(t)
			}
			time.Sleep(time.Millisecond * 1)
		}
	}()
}

func createFile(t time.Time) {
	sYear := strconv.Itoa(t.Year())
	sMonth := strconv.Itoa(int(t.Month()))
	sDay := strconv.Itoa(t.Day())
	sHour := strconv.Itoa(t.Hour())

	if len(sMonth) == 1 {
		sMonth = fmt.Sprintf("0%s", sMonth)
	}

	if len(sDay) == 1 {
		sDay = fmt.Sprintf("0%s", sDay)
	}

	if len(sHour) == 1 {
		sHour = fmt.Sprintf("0%s", sHour)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath := fmt.Sprintf("%s/logs/%s%s/%s/", pwd, sYear, sMonth, sDay)
	err = os.MkdirAll(filePath, 0777)
	if err != nil {
		panic(err)
	}

	fileName := fmt.Sprintf("%s.log", t.Format("2006010215"))
	fileLog, err := os.OpenFile(filePath+fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	Logger.SetOutput(io.MultiWriter(fileLog, os.Stdout))
}

// PrintErr ...
func PrintErr(ctx iris.Context, spot, errMsg string) {
	Logger.Errorf("RequestID: %s, %s : %+v", ctx.Values().GetString(ctxkey.RequestID), spot, errMsg)
}

// Trace ...
func Trace(ctx iris.Context, spot, errMsg string) {
	Logger.Infof("RequestID: %s, %s : %+v", ctx.Values().GetString(ctxkey.RequestID), spot, errMsg)
}

// PrintErrRedirect ...
func PrintErrRedirect(ctx iris.Context, spot, errMsg string) {
	PrintErr(ctx, spot, errMsg)
	ctx.Redirect("/", iris.StatusTemporaryRedirect)
}
