package domain

import "time"

// auth
type Auth struct {
	Seq        int    `json:"seq"`
	Key        string `json:"key"`
	Password   string `json:"pwd"`
	RemoteAddr string `json:"ip"`
}

type B2bLive struct {
	LiveSeq     string    `json:"live_seq"`
	ServiceSeq  int       `json:"service_seq"`
	RoomName    string    `json:"room_name"`
	AccessPoint string    `json:"access_point"`
	Videoroom   string    `json:"videoroom"`
	Textroom    string    `json:"textroom"`
	Record      bool      `json:"record"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Callback	string    `json:"callback"`
}

type MediaChannel struct {
	LiveSeq     string `json:"live_seq"`
	RoomName    string `json:"room_name"`
	AccessPoint string `json:"access_point"`
	Videoroom   string `json:"videoroom"`
	Textroom    string `json:"textroom"`
	Record      bool   `json:"record"`
	Description string `json:"description"`
}

// request format
type ReqMessage struct {
	Janus       string  `json:"janus"`
	Plugin      string  `json:"plugin"`
	Transaction string  `json:"transaction"`
	SessionId	uint64	`json:"session_id"`
	Body        ReqBody `json:"body"`
}

type ReqBody struct {
	Request    string `json:"request"`
	Textroom   string `json:"textroom"`
	Ptype      string `json:"ptype"`
	Publishers int	`json:"publishers"`
	Record     bool   `json:"record"`
	RecDir     string `json:"rec_dir"`
	FileName   string `json:"filename"`
	Videocodec string `json:"videocodec"`
	Audiocodec string `json:"audiocodec"`
	VideoorientExt	bool `json:"videoorient_ext"`
	Bitrate    uint64 `json:"bitrate"`
	AdminKey   string `json:"admin_key"`
	Secret     string `json:"secret"`
	Room       uint64 `json:"room"`
	Display    string `json:"display"`
	Id         uint64 `json:"id"`
	Post	   string	`json:"post"`
}

// response format
type JanusData struct {
	Id uint64 `json:"id"`
}

type JanusCreate struct {
	Janus       string    `json:"janus"`
	SessionId   uint64    `json:"session_id"`
	Transaction string    `json:"transaction"`
	Data        JanusData `json:"data"`
}

type JanusPluginVideoroomParticipant struct {
	Id      uint64 `json:"id"`
	Display string `json:"display"`
}

type JanusPluginVideoroomDetail struct {
	Room            uint64 `json:"room"`
	Description     string `json:"description"`
	PinRequired     bool   `json:"pin_required"`
	MaxPublishers   int    `json:"max_publishers"`
	Bitrate         uint64 `json:"bitrate"`
	FirFreq         int    `json:"fir_freq"`
	RequirePvtid    bool   `json:"require_pvtid"`
	NotifyJoining   bool   `json:"notify_joining"`
	Audiocodec      string `json:"audiocodec"`
	Videocodec      string `json:"videocodec"`
	Record          bool   `json:"record"`
	NumParticipants int    `json:"num_participants"`
}

type JanusPluginVideoroomData struct {
	Videoroom    string                            `json:"videoroom"`
	Textroom     string                            `json:"textroom"`
	Room         uint64                            `json:"room"`
	Exists		 bool							   `json:"exists"` 	
	Permanent    bool                              `json:"permanent"`
	RoomList     []JanusPluginVideoroomDetail      `json:"list"`
	Participants []JanusPluginVideoroomParticipant `json:"participants"`
}

type JanusPluginData struct {
	Plugin string                   `json:"plugin"`
	Data   JanusPluginVideoroomData `json:"data"`
}

type JanusPluginCreate struct {
	Janus       string          `json:"janus"`
	SessionId   uint64          `json:"session_id"`
	Transaction string          `json:"transaction"`
	Sender      uint64          `json:"sender"`
	PluginData  JanusPluginData `json:"plugindata"`
}

type PluginHandle struct {
	VideoroomHandleId uint64
	ChatHandleId      uint64
}
