package middleware

import (
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
	"github.com/teamgrit-lab/cojam/component/tglog"
	"github.com/teamgrit-lab/cojam/config"
)

// RDBTx ...
func RDBTx(ctx iris.Context) {

	var isTx bool
	tx := config.CF.DBConn.RDB

	switch ctx.Method() {
	case "POST", "PUT", "DELETE":
		ctx.Values().Set("RDBRollBack", false)
		isTx = true
		tx = tx.Begin()
		//tglog.Logger.Infof("[RDB Transaction] RequestID: %s, RDB_TX Begin()", ctx.Values().GetString("RequestID"))
		ctx.Values().Set(ctxkey.RDB_CONN, tx)
	default:
		isTx = false
		ctx.Values().Set(ctxkey.RDB_CONN, tx)
		//tglog.Logger.Infof("[RDB Transaction] RequestID: %s, Not TX", ctx.Values().GetString("RequestID"))
	}

	defer func(isTx bool) {
		if isTx {
			if ctx.GetStatusCode() >= 500 || ctx.Values().Get("RDBRollBack").(bool) {
				tx.Rollback()
				tglog.Logger.Errorf("[RDB Transaction] RequestID: %s, RDB_TX Rollback()", ctx.Values().GetString(ctxkey.RequestID))
			} else {
				tx.Commit()
				//tglog.Logger.Infof("[RDB Transaction] RequestID: %s, RDB_TX Commit()", ctx.Values().GetString("RequestID"))
			}
		}
	}(isTx)

	ctx.Next()
}
