package signalapi // import "github.com/teamgrit-lab/cojam/component/signalapi"

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
	"github.com/teamgrit-lab/cojam/component/ctxkey"
	"github.com/teamgrit-lab/cojam/component/session"
	"github.com/teamgrit-lab/cojam/component/util"
	"github.com/teamgrit-lab/cojam/config"
)

const SERVIOCE_GROUP_MEDIA = "MEDIA"
const SERVICE_TYPE_LIVE = "LIVE"
const SERVICE_TYPE_P2P = "P2P"
const RTCMODEL_SFU = "sfu"
const RTCMODEL_MESH = "mesh"

// LiveRoom ...
type LiveRoom struct {
	SignalServerURL string
	UserAgent       string
	ServiceGroup    string
	ServiceType     string
	LiveSeq         string
	UserSessionID   string
}

// RespLiveRoom ...
type RespLiveRoom struct {
	RoomToken             string `json:"response"`
	Transaction           string `json:"transaction"`
	CojamServiceSessionID string
}

func NewMakeCojamServiceID(ctx iris.Context) (string, error) {
	serviceGroup := ctx.Params().Get("service_group")
	serviceType := ctx.Params().Get("service_type")
	liveSeq := ctx.Params().Get("live_seq")
	sessionID := ctx.Params().Get("session_id")

	if len(serviceGroup) == 0 {
		return "", fmt.Errorf("Fail Verify Parameter")
	}

	if len(serviceType) == 0 {
		return "", fmt.Errorf("Fail Verify Parameter")
	}

	if len(sessionID) == 0 {
		return "", fmt.Errorf("Fail Verify Parameter")
	}

	liveRoom := &LiveRoom{
		UserAgent:     ctx.Values().Get(ctxkey.UserAgent).(string),
		ServiceGroup:  serviceGroup,
		ServiceType:   serviceType,
		LiveSeq:       liveSeq,
		UserSessionID: sessionID,
	}

	id := liveRoom.MakeGlobalSessionID()

	return id, nil

}

// NewCreateLiveRoomByMakeRoom ...
func NewCreateLiveRoomByMakeRoom(ctx iris.Context) (*RespLiveRoom, error) {
	serviceType := ctx.Params().Get("service_type")

	if len(serviceType) == 0 {
		return nil, fmt.Errorf("Fail Verify Parameter")
	}

	liveRoom := &LiveRoom{
		SignalServerURL: config.CF.Prop.API.WebrtcSignal.URL,
		UserAgent:       ctx.Values().Get(ctxkey.UserAgent).(string),
		ServiceGroup:    SERVIOCE_GROUP_MEDIA,
		ServiceType:     serviceType,
		LiveSeq:         "0",
		UserSessionID:   ctx.Values().Get(ctxkey.SessionKey).(string),
	}

	respLiveRoom, err := liveRoom.CreateLiveRoom()
	if err != nil {
		return nil, err
	}

	return respLiveRoom, nil
}

// CreateLiveRoom ...
func (liveRoom *LiveRoom) CreateLiveRoom() (respLiveRoom *RespLiveRoom, err error) {
	//fmt.Printf("LiveRoom.CreateLiveRoom : %+v\n", liveRoom)

	transactionCode := util.MakeUniqueID()
	cojamServiceSessionID := liveRoom.MakeGlobalSessionID()

	rtcModel := getRTCModel(liveRoom.ServiceType)

	client := &http.Client{
		Timeout: time.Duration(3 * time.Second),
	}

	var requestData = []byte(`{
		"type": "create",
		"title": "",
		"id": "` + liveRoom.LiveSeq + `",
		"owner": "` + cojamServiceSessionID + `",
		"user_id": ["` + cojamServiceSessionID + `", "` + cojamServiceSessionID + `"],
		"current_id": "` + cojamServiceSessionID + `",
		"description": "",
		"expiration_time": ` + strconv.Itoa(3600*24) + `,
		"max_users": ` + strconv.Itoa(100) + `,
		"rtc_model": "` + rtcModel + `",
		"enable_push_msg": ` + strconv.FormatBool(true) + `,
		"close_callback": ` + strconv.FormatBool(isCloseCallback(rtcModel)) + `,
		"transaction": "` + transactionCode + `"
		}`)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/room", liveRoom.SignalServerURL), bytes.NewBuffer(requestData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %+v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Do: %+v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll: %+v", err)
	}
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("close: %+v", err)
	}

	respLiveRoom = &RespLiveRoom{}
	err = json.Unmarshal(body, respLiveRoom)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(body): %+v", err)
	}

	if respLiveRoom.Transaction != transactionCode {
		return nil, fmt.Errorf("transaction 변조")
	}

	respLiveRoom.CojamServiceSessionID = cojamServiceSessionID
	return
}

// getRTCModel 는 sfu, mesh 선택.
//   - LIVE : sfu
//   - P2P : mesh
func getRTCModel(serviceType string) string {
	if serviceType == SERVICE_TYPE_LIVE {
		return RTCMODEL_SFU
	}
	return RTCMODEL_MESH
}

// isCloseCallback 은 serviceType 에 따른 콜백 유무
func isCloseCallback(rtcModel string) bool {
	if rtcModel == RTCMODEL_SFU {
		return true
	}
	return false
}

// DecodeGlobalSessionID ...
func DecodeGlobalSessionID(sid string) ([]string, error) {
	decodingSID, _ := base64.StdEncoding.DecodeString(sid)
	splitSID := strings.Split(string(decodingSID), "$")

	return splitSID, nil
}

// MakeGlobalSessionID ...
func (liveRoom *LiveRoom) MakeGlobalSessionID() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s$%s$%s$%s$%s$%s", liveRoom.UserAgent, liveRoom.ServiceGroup, liveRoom.ServiceType, util.MakeUniqueID(), liveRoom.LiveSeq, liveRoom.UserSessionID)))
}

// GetLiveSeqByDecodeGlobalSessionID ...
func GetLiveSeqByDecodeGlobalSessionID(sid string) (liveSeq uint64, err error) {
	decodeSID, err := DecodeGlobalSessionID(sid)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(decodeSID[4], 10, 64)
}

// GetUserSessionIDByDecodeGlobalSessionID ...
func GetUserSessionIDByDecodeGlobalSessionID(sid string) (userSessionID string, err error) {
	decodeSID, _ := DecodeGlobalSessionID(sid)
	return decodeSID[5], nil
}

// GetUserIDByDecodeGlobalSessionID ...
func GetUserIDByDecodeGlobalSessionID(sid string) (string, error) {
	userSessionID, err := GetUserSessionIDByDecodeGlobalSessionID(sid)
	if err != nil {
		return "", err
	}
	userID := session.GetUserIDBySessionID(userSessionID)
	return userID, nil
}
