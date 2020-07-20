package session //import "github.com/teamgrit-lab/cojam/component/session"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/util"
	"github.com/teamgrit-lab/cojam/config"
	"github.com/teamgrit-lab/cojam/mvc/domain"
)

const LOGIN_COOKIE_NAME = "cojam_sessionid"
const LOGIN_COOKIE_EXPIRE_TIME = 86400 //1일
const COOKIE_SEPERATOR = "#TEAMGRIT#"

const FLASH_MESSAGE_NAME = "flashMsg"

type UserSession struct {
	domain.User
	Roles     []string
	AvatarURL string `json:"avatar_url"`
}

// SetSession ...
func SetSession(key string, user *UserSession) error {
	conn := config.CF.DBConn.SessionRedis.Get()
	defer conn.Close()

	b, _ := json.Marshal(user)
	conn.Send("MULTI")
	conn.Send("SET", key, b)
	conn.Send("EXPIRE", key, LOGIN_COOKIE_EXPIRE_TIME)
	_, err := conn.Do("EXEC")
	if err != nil {
		return fmt.Errorf("session.SetSession : conn.Do('EXEC') : %+v", err)
	}

	return err
}

// GetSession ...
func GetSession(key string) (*UserSession, error) {
	conn := config.CF.DBConn.SessionRedis.Get()
	defer conn.Close()

	var userSession *UserSession

	b, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &userSession)
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

// GetUserIDBySessionID ...
func GetUserIDBySessionID(sessionID string) (userID string) {
	decodingID, _ := base64.StdEncoding.DecodeString(sessionID)
	deSessionID := strings.Split(string(decodingID), COOKIE_SEPERATOR)
	if len(deSessionID) == 2 {
		// uuid-userID
		userID = deSessionID[1]
	}
	return
}

// MakeCookieValue ...
func MakeCookieValue(userID string) string {
	return base64.StdEncoding.EncodeToString([]byte(util.MakeUniqueID() + COOKIE_SEPERATOR + userID))
}

// MakeExpireTime ...
func MakeExpireTime() time.Time {
	return time.Now().Add(LOGIN_COOKIE_EXPIRE_TIME * time.Second)
}

// DeleteSession ...
func DeleteSession(key string) error {
	conn := config.CF.DBConn.SessionRedis.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

// SetExpireTime ...
func SetExpireTime(key string) error {
	conn := config.CF.DBConn.SessionRedis.Get()
	defer conn.Close()

	_, err := redis.Bool(conn.Do("EXPIRE", key, LOGIN_COOKIE_EXPIRE_TIME))
	return err
}

// IsExistence ...
func IsExistence(key string) bool {
	conn := config.CF.DBConn.SessionRedis.Get()
	defer conn.Close()

	ok, _ := redis.Bool(conn.Do("EXISTS", key))
	return ok
}

// SetFlash ...
func SetFlash(ctx iris.Context, value []byte) {
	ctx.SetCookie(&http.Cookie{Name: FLASH_MESSAGE_NAME, Value: encode(value), Path: "/"})
}

// GetFlash ...
func GetFlash(ctx iris.Context) ([]byte, error) {
	c := ctx.GetCookie(FLASH_MESSAGE_NAME)
	flashMsg, err := decode(c)
	if err != nil {
		return nil, err
	}
	ctx.SetCookie(&http.Cookie{Name: FLASH_MESSAGE_NAME, MaxAge: -1, Expires: time.Unix(1, 0)})
	return flashMsg, nil
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}
