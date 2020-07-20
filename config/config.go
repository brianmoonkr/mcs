package config // import "github.com/teamgrit-lab/cojam/config"

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/globalsign/mgo"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/oauth"
	"github.com/teamgrit-lab/cojam/component/tglog"
)

// CF ...
var CF *Config

// Config ...
type Config struct {
	Prop          *Properties
	DBConn        *DBConnection
	ProjectRoot   string
	ExecutionMode string
	DomainURI     string
}

// DBConnection ...
type DBConnection struct {
	RDB          *gorm.DB
	Redis        *redis.Pool
	ConnRedis	 *redis.Pool
	SessionRedis *redis.Pool
	Mongo        *mgo.Session
}

// Properties ...
type Properties struct {
	Domain      string `json:"domain"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	MediaAdminPort	int	`json:"media_admin_port"`
	MediaAdminSecret	string	`json:"media_admin_secret"`
	ThumbOutStart	string	`json:"thumb_out_start"`
	ProjectName string `json:"project_name"`
	JanusRecCmd string	`json:"janus_rec_cmd"` 
	JanusRecPath string	`json:"janus_rec_path"` 
	RepositoryPath string `json:"repository_path"` 
	SharedAccessPath string `json:"shared_access_path"` 
	TextroomPost	string `json:"textroom_post"` 

	AdminInfo struct {
		Email string `json:"email"`
	} `json:"admin_info"`

	CdnInfo struct {
		CdnUrl string `json:"cdn_url"`
		UploadUrl string `json:"upload_url"`
		FtpId string `json:"ftp_id"`
		FtpPwd string `json:"ftp_pwd"`
	} `json:"cdn_info"`	

	Storage struct {
		RDB struct {
			URL  string `json:"url"`
			Open int    `json:"open"`
			Idle int    `json:"idle"`
		} `json:"rdb"`

		Mongo struct {
			URL    string   `json:"url"`
			Pool   int      `json:"pool"`
			DBName []string `json:"db_name"`
		} `json:"mongo"`

		ConnRedis struct {
			URL       string `json:"url"`
			MaxIdle   int    `json:"max_idle"`
			MaxActive int    `json:"max_active"`
			DBnum     int    `json:"dbnum"`
		} `json:"conn_redis"`

		SessionRedis struct {
			URL       string `json:"url"`
			MaxIdle   int    `json:"max_idle"`
			MaxActive int    `json:"max_active"`
			DBnum     int    `json:"dbnum"`
		} `json:"session_redis"`

		Redis struct {
			URL       string `json:"url"`
			MaxIdle   int    `json:"max_idle"`
			MaxActive int    `json:"max_active"`
			DBnum     int    `json:"dbnum"`
		} `json:"redis"`
	} `json:"storage"`

	API struct {
		OAuth struct {
			Facebook struct {
				OAuth
			}
			Google struct {
				OAuth
			}
			Naver struct {
				OAuth
			}
			Kakao struct {
				OAuth
			}
			Line struct {
				OAuth
			}
		} `json:"oauth"`
		Sendgrid struct {
			Key string `json:"key"`
		} `json:"sendgrid"`
		WebrtcSignal struct {
			URL string `json:"url"`
		} `json:"webrtc_signal"`
		OneSignal struct {
			Key string `json:"key"`
		} `json:"one_signal"`
		Teamgrit struct {
			VodUpload struct {
				Key string `json:"key"`
				URL string `json:"url"`
				CDN string `json:"cdn"`
			} `json:"vod_upload"`
			ImageUpload struct {
				Key string `json:"key"`
				URL string `json:"url"`
				CDN string `json:"cdn"`
			} `json:"image_upload"`
		} `json:"teamgrit"`
		Slack struct {
			Channel     string `json:"channel"`
			InComingURL string `json:"incoming_url"`
		} `json:"slack"`
	} `json:"api"`

	MediaServers []string `json:"media_servers"`
}

// OAuth ...
type OAuth struct {
	ID       string `json:"id"`
	Key      string `json:"key"`
	CallBack string `json:"callback"`
}

// InitConfig ...
func InitConfig(app *iris.Application) {
	CF = &Config{}
	initExecutionMode()
	initProperties()
	initDBConnection()
	setProjectRoot()
	setDomainURI()
	setOauthProviders()
}

func initExecutionMode() {
	CF.ExecutionMode = os.Getenv("COJAM_EXECMODE")
	if len(CF.ExecutionMode) == 0 {
		//CF.ExecutionMode = "local"
		CF.ExecutionMode = "dev"
	}

	tglog.Logger.Println("ExecutionMode : ", CF.ExecutionMode)
}

func initProperties() {
	configFile, err := os.Open("config/properties/config." + CF.ExecutionMode + ".json")
	defer configFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&CF.Prop)
	if err != nil {
		log.Fatal(err)
	}
}

func initDBConnection() {
	CF.DBConn = &DBConnection{}
	CF.DBConn.RDB = newRDB()
	CF.DBConn.ConnRedis = newConnRedis()
	//CF.DBConn.Redis = newRedis()
	//CF.DBConn.SessionRedis = newSessionRedis()
	//CF.DBConn.Mongo = newMongo()
}

func newRDB() *gorm.DB {
	var err error
	conn, err := gorm.Open("mysql", CF.Prop.Storage.RDB.URL)
	if err != nil {
		log.Fatal("ERROR RDB ", err)
	}
	conn.DB().SetMaxIdleConns(CF.Prop.Storage.RDB.Idle)
	conn.DB().SetMaxOpenConns(CF.Prop.Storage.RDB.Open)
	conn.DB().SetConnMaxLifetime(3 * time.Minute)
	conn.SetLogger(gorm.Logger{tglog.Logger})
	conn.LogMode(true)
	conn.SingularTable(true)

	fmt.Println("DB RDB Connection OK.")
	return conn
}

func newRedis() *redis.Pool {

	pool := &redis.Pool{
		MaxIdle:     CF.Prop.Storage.Redis.MaxIdle, // 항상 pool에 대기중인 conn 수
		MaxActive:   CF.Prop.Storage.Redis.MaxIdle, // 최대 연결 수
		Wait:        true,                          // maxActive 값을 넘을때 maxActive 값에 들어올때까지 conn blocking
		IdleTimeout: 3 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", CF.Prop.Storage.Redis.URL, redis.DialDatabase(CF.Prop.Storage.Redis.DBnum))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	fmt.Println("Redis Connection OK.")
	return pool
}

func newConnRedis() *redis.Pool {

	pool := &redis.Pool{
		MaxIdle:     CF.Prop.Storage.ConnRedis.MaxIdle, // 항상 pool에 대기중인 conn 수
		MaxActive:   CF.Prop.Storage.ConnRedis.MaxIdle, // 최대 연결 수
		Wait:        true,                          // maxActive 값을 넘을때 maxActive 값에 들어올때까지 conn blocking
		IdleTimeout: 3 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", CF.Prop.Storage.ConnRedis.URL, redis.DialDatabase(CF.Prop.Storage.ConnRedis.DBnum))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	fmt.Println("ConnRedis Connection OK.")
	return pool
}

func newSessionRedis() *redis.Pool {
	pool := &redis.Pool{
		MaxIdle:     CF.Prop.Storage.SessionRedis.MaxIdle, // 항상 pool에 대기중인 conn 수
		MaxActive:   CF.Prop.Storage.SessionRedis.MaxIdle, // 최대 연결 수
		Wait:        true,                                 // maxActive 값을 넘을때 maxActive 값에 들어올때까지 conn blocking
		IdleTimeout: 3 * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", CF.Prop.Storage.SessionRedis.URL, redis.DialDatabase(CF.Prop.Storage.SessionRedis.DBnum))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	fmt.Println("SessionRedis Connection OK.")
	return pool
}

func newMongo() *mgo.Session {
	dial := fmt.Sprintf("%s?maxPoolSize=%d", CF.Prop.Storage.Mongo.URL, CF.Prop.Storage.Mongo.Pool)
	conn, err := mgo.Dial(dial)
	if err != nil {
		log.Fatal("ERROR Mongo", err)
	}
	conn.SetMode(mgo.Monotonic, true)

	fmt.Println("DB Mongo Connection OK.")
	return conn
}

func setProjectRoot() {
	var err error
	CF.ProjectRoot, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
}

func setDomainURI() {
	fmt.Printf("mode: %s", CF.ExecutionMode)

	if CF.ExecutionMode == "local" {
		CF.DomainURI = fmt.Sprintf("http://%s:%s", CF.Prop.Domain, strconv.Itoa(CF.Prop.Port))
		return
	}

	if CF.ExecutionMode == "dev" {
		CF.DomainURI = fmt.Sprintf("https://%s", CF.Prop.Domain)
		return
	}

	if CF.ExecutionMode == "real" {
		CF.DomainURI = fmt.Sprintf("https://%s", CF.Prop.Domain)
		return
	}
}

func setOauthProviders() {
	oauth.UseProviders(
		oauth.NewKakao(CF.Prop.API.OAuth.Kakao.ID, CF.Prop.API.OAuth.Kakao.Key, CF.Prop.API.OAuth.Kakao.CallBack),
		oauth.NewGoogle(CF.Prop.API.OAuth.Google.ID, CF.Prop.API.OAuth.Google.Key, CF.Prop.API.OAuth.Google.CallBack),
		oauth.NewFacebook(CF.Prop.API.OAuth.Facebook.ID, CF.Prop.API.OAuth.Facebook.Key, CF.Prop.API.OAuth.Facebook.CallBack),
		oauth.NewNaver(CF.Prop.API.OAuth.Naver.ID, CF.Prop.API.OAuth.Naver.Key, CF.Prop.API.OAuth.Naver.CallBack),
		oauth.NewLine(CF.Prop.API.OAuth.Line.ID, CF.Prop.API.OAuth.Line.Key, CF.Prop.API.OAuth.Line.CallBack),
	)
}
