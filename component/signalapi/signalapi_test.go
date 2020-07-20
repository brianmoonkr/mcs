package signalapi_test

import (
	"fmt"
	"testing"

	"github.com/teamgrit-lab/shop-beta/component/signalapi"
)

func TestCreateLiveRoom(t *testing.T) {
	liveRoom := &signalapi.LiveRoom{
		SignalServerURL: "https://dev-sig.cojam.tv",
		ServiceGroup:    "MEDIA",
		ServiceType:     "LIVE",
		LiveSeq:         "12",
		UserSessionID:   "ZTc3NzRiM2QtNWFjNC00OGI5LThjZDMtMzgzMmNlZjExMjg5c2Vuc2Vva2lAdGVhb=",
	}

	resp, err := liveRoom.CreateLiveRoom()
	if err != nil {
		fmt.Printf("err : liveRoom.CreateLiveRoom() : %+v", err)
	}

	fmt.Printf("CreateLiveRoom : %+v\n", resp)
}

func TestDecodeGlobalSessionID(t *testing.T) {
	sid := "TUVESUEkTElWRSRiajAxc2pxM3E1NjV2Z2k1YzQ2MCQxNSRaVGMzTnpSaU0yUXROV0ZqTkMwME9HSTVMVGhqWkRNdE16Z3pNbU5sWmpFeE1qZzVjMlZ1YzJWdmEybEFkR1ZoYj0kc2lnMS5jb2phbS50diQyMDE5LTA0LTIzIDA1OjAzOjMxJDIwMTktMDQtMjMgMDU6MDM6MzE="
	_, _ := DecodeGlobalSessionID(sid)
}
