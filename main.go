package main // import "github.com/teamgrit-lab/cojam"

import (
	stdCtx "context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net"
	"net/url"
	"strings"
	"encoding/json"
	"strconv"
	"log"
	"errors"

	"github.com/go-resty/resty/v2"
//	"github.com/garyburd/redigo/redis"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	
	"github.com/teamgrit-lab/cojam/component/tglog"
	"github.com/teamgrit-lab/cojam/config"
	"github.com/teamgrit-lab/cojam/middleware"
	"github.com/teamgrit-lab/cojam/mvc/controller"
	"github.com/teamgrit-lab/cojam/mvc/dao"

	"github.com/teamgrit-lab/cojam/mvc/domain"
	"github.com/teamgrit-lab/cojam/mvc/service"
	. "github.com/teamgrit-lab/cojam/transcode"
	"github.com/robfig/cron"
)

var cronInstace *cron.Cron
var taskFunc = make(map[string]func())

func GetCrontabInstance() *cron.Cron {
	if cronInstace != nil {
		return cronInstace
	}
	cronInstace = cron.New()
	cronInstace.Start()

	iris.RegisterOnInterrupt(func() {
		cronInstace.Stop()
	})
	return cronInstace
}

func AddTaskFunc(name string, schedule string, f func()) {
	if _, ok := taskFunc[name]; !ok {
		fmt.Println("Add a new task:", name)

		cInstance := GetCrontabInstance()
		cInstance.AddFunc(schedule, f)

		taskFunc[name] = f
	} else {
		fmt.Println("Don't add same task `" + name + "` repeatedly!")
	}
}

func getHostName(strUrl string) (string, error) {
	u, err := url.Parse(strUrl)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return u.Hostname(), nil
}

func getHostPort(strUrl string) (string, error) {
	u, err := url.Parse(strUrl)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	_, port, _ := net.SplitHostPort(u.Host)

	return port, nil
}

func main() {
	app := iris.New()

	initApp(app)
	tglog.NewLogFile()
	config.InitConfig(app)

	if len(os.Args) > 1 && os.Args[1] == "init" {		
		// redisConn
		redisPool := config.CF.DBConn.ConnRedis
		redisConn := redisPool.Get()
		redisConn.Do("FLUSHDB")
		err := redisConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	/*	blocked in 2020-06-15
		// redis
		redisPool = config.CF.DBConn.Redis
		redisConn = redisPool.Get()
		redisConn.Do("FLUSHDB")
		err = redisConn.Close()
		if err != nil {
			log.Fatal(err)
		}
		// redisSession
		redisPool := config.CF.DBConn.SessionRedis
		redisConn := redisPool.Get()
		redisConn.Do("FLUSHDB")
		err := redisConn.Close()
		if err != nil {
			log.Fatal(err)
		}
		*/
		tglog.Logger.Printf("===== REDIS DB Flushed. =====\n")
	}
	
	initMediaConnection()
	
	go AddTaskFunc("mcs scheduled task", "0 10 19 * * ?", ChannelManager)
	
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.HTML("<h1>Oops!!!!!!!!!!!!!!!</h1><h1>500</h1>")
	})
	
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.HTML("<h1>Oops !!!!!!!!!!!!!!!</h1><h1>404</h1>")
	})
	
	mvc.Configure(app.Party("/api",
		middleware.APIURLAuth,
		middleware.BasicAuth,
		middleware.Cors,
		middleware.RDBTx,
		middleware.Recover,
	), apiParty)

	tglog.Logger.Printf("\n\n\n\n\n")
	tglog.Logger.Println("===== ****************************** =====")
	tglog.Logger.Println("===== ****************************** =====")
	tglog.Logger.Println("===== START : Media Control Service  =====")
	tglog.Logger.Println("===== **********(with iris)********* =====")
	tglog.Logger.Printf("===== ****************************** =====\n")
	tglog.Logger.Println("IRIS Version : ", iris.Version)

	app.Run(
		iris.Addr(fmt.Sprintf("%s:%d", config.CF.Prop.Address, config.CF.Prop.Port)),
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed), // skip err server closed when CTRL/CMD+C pressed:
		iris.WithOptimizations,                        // enables faster json serialization and more:
		iris.WithoutPathCorrectionRedirection,
	)

}

func apiParty(mvc *mvc.Application) {
	registerDI(mvc)
	mvc.Handle(new(controller.JanusController))

}

// registerDI 는 service, dao Dependency Injection.
func registerDI(mvc *mvc.Application) {
	//DAO
	janusDAO := dao.NewJanusDAO()
	//Service
	janusService := service.NewJanusService()
	//
	janusServiceDI := map[string]interface{}{
		"janusDAO": janusDAO,
	}

	// Service Dependency injection map
	janusService.SetService(janusServiceDI)
	
	// Service Register
	mvc.Register(
		janusService,
	)
}

func initApp(app *iris.Application) {
	app.Use(iris.Gzip)
	go shutDown(app)
}

func shutDown(app *iris.Application) {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		// kill -SIGINT XXXX or Ctrl+c
		os.Interrupt,
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill  is equivalent with the syscall.Kill
		os.Kill,
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	select {
	case <-ch:

		tglog.Logger.Printf("\n\n")
		tglog.Logger.Println("===== ShutDown START : CoJam Server With IRIS =====")
		config.CF.DBConn.RDB.Close()
		tglog.Logger.Println("Close RDB")
		//config.CF.DBConn.Mongo.Close()
		config.CF.DBConn.ConnRedis.Close()
		tglog.Logger.Println("Close ConnRedis")
		config.CF.DBConn.Redis.Close()
		tglog.Logger.Println("Close Redis")
		config.CF.DBConn.SessionRedis.Close()
		tglog.Logger.Println("Close SessionRedis")
		timeout := 10 * time.Second
		ctx, cancel := stdCtx.WithTimeout(stdCtx.Background(), timeout)
		defer cancel()
		tglog.Logger.Println("===== ShutDown END : CoJam Server With IRIS =====")

		app.Shutdown(ctx)
		//cronJob.Stop()

	}
}

func initMediaConnection() {
	// conn redis
	redisPool := config.CF.DBConn.ConnRedis
	redisConn := redisPool.Get()
	defer redisConn.Close()
	
	janusReq := struct {
		Janus      string   `json:"janus"`
		Transaction    string `json:"transaction"`
		AdminSecret string `json:"admin_secret"`
	}{}

	janusRes := struct {
		Janus      string   `json:"janus"`
		Transaction    string `json:"transaction"`
		Sessions 	[]uint64	`json:"sessions"`
	}{}
	
	client := resty.New()

	mediaServers := config.CF.Prop.MediaServers
	
	for i := 0; i < len(mediaServers); i++ {
		server := mediaServers[i]
		cons := strings.Split(server, "$$")
		httpConn := cons[0]

		tr := strconv.FormatInt(time.Now().Unix(), 10)

		janusReq.Janus = "ping"
		janusReq.Transaction = tr
		jsonMessage, _ := json.Marshal(janusReq)

		resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonMessage)).
		Post(fmt.Sprintf("%s/janus", httpConn))

		if err != nil {
			log.Fatal(err)
		}

		json.Unmarshal(resp.Body(), &janusRes)
		
		if janusRes.Janus == "pong" {	
			lastIndex := strings.LastIndex(httpConn, ":")
			httpAdminConn := fmt.Sprintf("%s:%d", httpConn[0:lastIndex], config.CF.Prop.MediaAdminPort)

			janusReq.Janus ="list_sessions"
			janusReq.Transaction = strconv.FormatInt(time.Now().Unix(), 10)
			janusReq.AdminSecret = config.CF.Prop.MediaAdminSecret

			jsonMessage, _ = json.Marshal(janusReq)
			
			resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(string(jsonMessage)).
			Post(fmt.Sprintf("%s/admin", httpAdminConn))

			if err != nil {
				log.Fatal(err)
			}

			json.Unmarshal(resp.Body(), &janusRes)
						
			sessionResp, _ := createJanusSession(httpConn, tr)
			sessionId := sessionResp.Data.Id
			handleResp, _ := attachJanusVideoPlugin(httpConn, tr, sessionId)  
			handleId := handleResp.Data.Id

			redisConn.Do("HSET", fmt.Sprintf("conn_key#%s", server), "count", len(janusRes.Sessions))
			redisConn.Do("HSET", fmt.Sprintf("conn_key#%s", server), "session", sessionId)
			redisConn.Do("HSET", fmt.Sprintf("conn_key#%s", server), "handle", handleId)
			redisConn.Do("SADD", "conn_keys", server)

		}

	}

}

func ChannelManager() {

	var reqBody domain.ReqBody
	var reqMessage domain.ReqMessage

	client := resty.New()
	
	mediaServers := config.CF.Prop.MediaServers
	
	for i := 0; i < len(mediaServers); i++ {
		server := mediaServers[i]
		cons := strings.Split(server, "$$")
		httpConn := cons[0]
		wssConn := cons[1]

		tr := strconv.FormatInt(time.Now().Unix(), 10)

		sessionResp, _ := createJanusSession(httpConn, tr)
		sessionId := sessionResp.Data.Id
		handleResp, _ := attachJanusVideoPlugin(httpConn, tr, sessionId)  
		videoHandleId := strconv.FormatUint(handleResp.Data.Id, 10)
		//videoHandleId := handleResp.Data.Id

		handleResp, _ = attachJanusTextPlugin(httpConn, tr, sessionId)  
		textHandleId := strconv.FormatUint(handleResp.Data.Id, 10)
		//textHandleId := handleResp.Data.Id

		strSessionId := strconv.FormatUint(sessionId, 10)
		
		reqBody.Request = "list"
		reqMessage.Janus = "message"
		reqMessage.Transaction = tr
		reqMessage.Body = reqBody
		
		jsonMessage, err := json.Marshal(reqMessage)
		if err != nil {		
			return
		}
		
		resp, err := client.R().
			SetPathParams(map[string]string{
				"sessionId":    strSessionId,
				"pluginHandle": videoHandleId,
			}).
			SetHeader("Content-Type", "application/json").
			SetBody(string(jsonMessage)).
			Post(fmt.Sprintf("%s/janus/{sessionId}/{pluginHandle}", httpConn))

		if err != nil {
			log.Fatal(err)
			return
		}
		
		//fmt.Println("=== list room ===")
		//fmt.Println(resp)

		var videoroomCreateResp domain.JanusPluginCreate

		json.Unmarshal(resp.Body(), &videoroomCreateResp)

		roomList := videoroomCreateResp.PluginData.Data.RoomList

		rdbConn := config.CF.DBConn.RDB
		for _, roomItem := range roomList {
			b2bLive :=  &domain.B2bLive{}
			rdbConn.Where("videoroom = ? AND access_point = ?", strconv.FormatUint(roomItem.Room, 10), wssConn).Find(&b2bLive)

			if b2bLive.Status == "1000" {
				rdbConn.Model(&b2bLive).Where("videoroom = ? AND access_point = ?", strconv.FormatUint(roomItem.Room, 10), wssConn).Update("Status", "0000")

				client := resty.New()

				// Destroy Janus Videoroom
				var reqMessage domain.ReqMessage
				var reqBody domain.ReqBody
				reqBody.Request = "destroy"
				reqBody.Room = roomItem.Room
				reqBody.Secret = "teamgrit"

				reqMessage.Janus = "message"
				reqMessage.Transaction = tr
				reqMessage.Body = reqBody

				jsonMessage, err := json.Marshal(reqMessage)
				if err != nil {
					log.Fatal(err)
					return
				}

				_, err = client.R().
					SetPathParams(map[string]string{
						"sessionId":    strSessionId,
						"pluginHandle": videoHandleId,
					}).
					SetHeader("Content-Type", "application/json").
					SetBody(string(jsonMessage)).
					Post(fmt.Sprintf("%s/janus/{sessionId}/{pluginHandle}", httpConn))

				if err != nil {
					log.Fatal(err)
					return
				}

				fmt.Println("=== destroy videoroom room ===")
				//fmt.Println(resp)

				// destroy text room
				reqBody.Textroom = "destroy"
				reqBody.Room = roomItem.Room
				reqMessage.Body = reqBody

				jsonMessage, err = json.Marshal(reqMessage)
				if err != nil {
					log.Fatal(err)
					return
				}

				_, err = client.R().
					SetPathParams(map[string]string{
						"sessionId":    strSessionId,
						"pluginHandle": textHandleId,
					}).
					SetHeader("Content-Type", "application/json").
					SetBody(string(jsonMessage)).
					Post(fmt.Sprintf("%s/janus/{sessionId}/{pluginHandle}", httpConn))

				if err != nil {
					log.Fatal(err)
					return
				}


			}
		} 

	}
	fmt.Println("############### Daily Job Process is Done at 4:00 ##################")
}

func contains(s []uint64, e uint64) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

// function :: create janus session
func createJanusSession(conn string, tr string) (domain.JanusCreate, error) {

	client := resty.New()

	// Janus Session Create
	var sessionResp domain.JanusCreate

	var reqMessage domain.ReqMessage
	reqMessage.Janus = "create"
	reqMessage.Transaction = tr

	jsonCreate, err := json.Marshal(reqMessage)
	if err != nil {
		log.Fatal(err)
		return sessionResp, errors.New("request message parsing error")
	}

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonCreate)).
		Post(fmt.Sprintf("%s/janus", conn))

	if err != nil {
		log.Fatal(err)
		return sessionResp, errors.New("errors in getting response from janus")
	}

	//fmt.Println(resp)

	json.Unmarshal(resp.Body(), &sessionResp)

	return sessionResp, nil
}

// function :: attach video room plugin
func attachJanusVideoPlugin(conn string, tr string, session uint64) (domain.JanusCreate, error) {

	var pluginResp domain.JanusCreate

	if session == 0 {
		return pluginResp, errors.New("request message parsing error")
	}

	client := resty.New()
	// Janus Videoroom Plugin Attach
	var reqMessage domain.ReqMessage

	reqMessage.Janus = "attach"
	reqMessage.Transaction = tr
	reqMessage.Plugin = "janus.plugin.videoroom"

	jsonAttach, err := json.Marshal(reqMessage)
	if err != nil {
		log.Fatal(err)
		return pluginResp, errors.New("request message parsing error")
	}

	resp, err := client.R().
		SetPathParams(map[string]string{
			"sessionId": strconv.FormatUint(session, 10),
		}).
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonAttach)).
		Post(fmt.Sprintf("%s/janus/{sessionId}", conn))

	if err != nil {
		log.Fatal(err)
		return pluginResp, errors.New("attach videoroom plugin")
	}

	//fmt.Println(resp)

	json.Unmarshal(resp.Body(), &pluginResp)

	return pluginResp, nil
}

// function :: attach video room plugin
func attachJanusTextPlugin(conn string, tr string, session uint64) (domain.JanusCreate, error) {

	var pluginResp domain.JanusCreate

	if session == 0 {
		return pluginResp, errors.New("request message parsing error")
	}

	client := resty.New()
	// Janus Textroom Plugin Attach
	var reqMessage domain.ReqMessage

	reqMessage.Janus = "attach"
	reqMessage.Transaction = tr
	reqMessage.Plugin = "janus.plugin.textroom"

	jsonAttach, err := json.Marshal(reqMessage)
	if err != nil {
		log.Fatal(err)
		return pluginResp, errors.New("request message parsing error")
	}

	resp, err := client.R().
		SetPathParams(map[string]string{
			"sessionId": strconv.FormatUint(session, 10),
		}).
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonAttach)).
		Post(fmt.Sprintf("%s/janus/{sessionId}", conn))

	if err != nil {
		log.Fatal(err)
		return pluginResp, errors.New("attach textroom plugin")
	}

	//fmt.Println(resp)

	json.Unmarshal(resp.Body(), &pluginResp)

	return pluginResp, nil
}

func removeRoom(server string, sessionId string, handleId string, tr string, roomId uint64, record bool, callback string) {

	cons := strings.Split(server, "$$")
	httpConn := cons[0]
	wssConn := cons[1]

	//sessionKey := fmt.Sprintf("%s_%s", wssConn, strconv.FormatUint(roomId, 10))

	hostPort, err := getHostPort(wssConn)
	if err != nil {
	        log.Fatal(err)	
		return
	}

	//removeRedisRoomInfo(sessionKey)	blocked in 2020-06-15

	client := resty.New()

	// Destroy Janus Videoroom
	var reqMessage domain.ReqMessage
	var reqBody domain.ReqBody
	reqBody.Request = "destroy"
	reqBody.Room = roomId
	reqBody.Secret = "teamgrit"

	reqMessage.Janus = "message"
	reqMessage.Transaction = tr
	reqMessage.Body = reqBody

	jsonMessage, err := json.Marshal(reqMessage)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = client.R().
		SetPathParams(map[string]string{
			"sessionId":    sessionId,
			"pluginHandle": handleId,
		}).
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonMessage)).
		Post(fmt.Sprintf("%s/janus/{sessionId}/{pluginHandle}", httpConn))

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("=== destroy videoroom room ===")
	//fmt.Println(resp)

	// destroy text room
	reqBody.Textroom = "destroy"
	reqBody.Room = roomId
	reqMessage.Body = reqBody

	jsonMessage, err = json.Marshal(reqMessage)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = client.R().
		SetPathParams(map[string]string{
			"sessionId":    sessionId,
			"pluginHandle": handleId,
		}).
		SetHeader("Content-Type", "application/json").
		SetBody(string(jsonMessage)).
		Post(fmt.Sprintf("%s/janus/{sessionId}/{pluginHandle}", httpConn))

	if err != nil {
		log.Fatal(err)
		return
	}

	//fmt.Println("=== destroy text room ===")
	//fmt.Println(resp)

	b2bLive :=  &domain.B2bLive{}

	rdbConn := config.CF.DBConn.RDB
	rdbConn.Where("videoroom = ?", strconv.FormatUint(roomId, 10)).Find(&b2bLive).Update("Status", "0000")

	if record {
		go TranscodeBySession(b2bLive.LiveSeq, strconv.FormatUint(roomId, 10), hostPort, callback)
	}

}
/* * blocked in 2020-06-15
func removeRedisRoomInfo(key string) {

	redisPool := config.CF.DBConn.SessionRedis
	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err := redisConn.Do("DEL", key)
	if err != nil {
		panic(err)
	}
}
*/
func setSessionCountByConn(key string) {

	server := strings.Split(key, "#")[1]
	cons := strings.Split(server, "$$")
	httpConn := cons[0]

	lastIndex := strings.LastIndex(httpConn, ":")
	httpAdminConn := fmt.Sprintf("%s:%d", httpConn[0:lastIndex], config.CF.Prop.MediaAdminPort)

	client := resty.New()

	janusAdminReq := struct {
		Janus      string   `json:"janus"`
		Transaction    string `json:"transaction"`
		AdminSecret string `json:"admin_secret"`
	}{}

	tr := strconv.FormatInt(time.Now().Unix(), 10)

	janusAdminReq.Janus ="list_sessions"
	janusAdminReq.Transaction = tr
	janusAdminReq.AdminSecret = config.CF.Prop.MediaAdminSecret

	jsonMessage, _ := json.Marshal(janusAdminReq)
	
	resp, err := client.R().
	SetHeader("Content-Type", "application/json").
	SetBody(string(jsonMessage)).
	Post(fmt.Sprintf("%s/admin", httpAdminConn))

	if err != nil {
		log.Fatal(err)
	}

	janusAdminRes := struct {
		Janus      string   `json:"janus"`
		Transaction    string `json:"transaction"`
		Sessions 	[]uint64	`json:"sessions"`
	}{}

	json.Unmarshal(resp.Body(), &janusAdminRes)

	redisPool := config.CF.DBConn.ConnRedis
	redisConn := redisPool.Get()
	defer redisConn.Close()

	_, err = redisConn.Do("HSET", key, "count", len(janusAdminRes.Sessions))
	if err != nil {
		panic(err)
	} 

}



