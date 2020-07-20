package constdf //import "github.com/teamgrit-lab/cojam/component/constdf"

import "strconv"

const OG_TAG_TITLE = "함께하는 참여 방송, CoJam.TV"
const OG_TAG_DESC = "함께해서 재미난 실시간 참여형 방송 서비스, 코잼TV, CoJamTV, 코잼티비, 코잼, CoJam, CoJam.TV"
const OG_TAG_IMAGE = "https://cojam.tv/static/cojam/img/img-og-main.jpg"

const LIVE_STATUS_REG = "1001"
const LIVE_STATUS_WAIT = "1002"
const LIVE_STATUS_ING = "1003"
const LIVE_STATUS_END = "1004"
const LIVE_STATUS_END_UNEXPECTEDLY = "1005"

const user_JOIN_INFO_PROVIDER_USER_KEY = "user_join_info:"
const user_JOIN_INFO_REDIS_EXPIRE_TIME = 1800 // 30분

const email_JOIN_AUTH_CODE_KEY = "email_join_auth_code_key:"
const email_JOIN_AUTH_CODE_KEY_EXPIRE_TIME = 180 // 3분

var VOD_FILE_TYPE = [2]string{"HLS", "MP4"}

// UserAuthCode --------------->
const USER_AUTH_CODE_ADMIN = "01"
const USER_AUTH_CODE_CORE = "02"
const USER_AUTH_CODE_JAMMER = "03"
const USER_AUTH_CODE_DEV = "04"

var user_AUTH_CODE = map[string]string{
	"01": "admin",
	"02": "core",
	"03": "jammer",
	"04": "dev",
}

// End Of UserAuthCode <---------

const GLOBAL_SESSION_ID_SERVICE_GROUP_MEDIA = "MEDIA"
const GLOBAL_SESSION_ID_SERVICE_TYPE_LIVE = "LIVE"
const GLOBAL_SESSION_ID_SERVICE_TYPE_P2P = "P2P"
const GLOBAL_SESSION_ID_USERAGENT_WEB = "WEB"
const GLOBAL_SESSION_ID_USERAGENT_APP = "APP"

// 고객센터 - 고객문의 Status
const CUSTOMER_CENTER_FAQ_STATUS_UNTREATED = "1001"
const CUSTOMER_CENTER_FAQ_STATUS_SUCCESS = "1002"

// Admin 고객센터 문의 검색
const CUSTOMER_CENTER_FAQ_SEARCH_FAQNUM = "1001"
const CUSTOMER_CENTER_FAQ_SEARCH_USERSEQ = "1002"
const CUSTOMER_CENTER_FAQ_SEARCH_USERNAME = "1003"
const CUSTOMER_CENTER_FAQ_SEARCH_QUESTION = "1004"
const CUSTOMER_CENTER_FAQ_SEARCH_ANSWER = "1005"

// Admin 사용자 검색
const USER_DATE_TYPE_ACCESS = "1001"
const USER_DATE_TYPE_CREATED = "1002"
const USER_SEARCH_USERSEQ = "1001"
const USER_SEARCH_USERNAME = "1002"
const USER_SEARCH_USEREMAIL = "1003"
const USER_SEARCH_CHANNELTITLE = "1004"
const USER_SEARCH_NICKNAME = "1005"

// Admin 사용자 블랙리스트 검색
const USER_BLACKLIST_SEARCH_USERSEQ = "1001"
const USER_BLACKLIST_SEARCH_USERNAME = "1002"
const USER_BLACKLIST_SEARCH_NICKNAME = "1003"
const USER_BLACKLIST_SEARCH_DESC = "1004"

// Admin 사용자 구독 검색
const USER_SUBSCRIPTION_SEARCH_USERSEQ = "1001"
const USER_SUBSCRIPTION_SEARCH_USERNAME = "1002"
const USER_SUBSCRIPTION_SEARCH_NICKNAME = "1003"
const USER_SUBSCRIPTION_SEARCH_CHANNELNUM = "1004"
const USER_SUBSCRIPTION_SEARCH_CHANNELNAME = "1005"

// Admin 사용자 좋아요 검색
const USER_LIKE_SEARCH_USERSEQ = "1001"
const USER_LIKE_SEARCH_USERNAME = "1002"
const USER_LIKE_SEARCH_NICKNAME = "1003"
const USER_LIKE_SEARCH_LIVESEQ = "1004"
const USER_LIKE_SEARCH_LIVENAME = "1005"

// Admin 사용자 회원탈퇴 검색
const USER_WITHDRAWAL_SEARCH_USERSEQ = "1001"
const USER_WITHDRAWAL_SEARCH_USERNAME = "1002"
const USER_WITHDRAWAL_SEARCH_NICKNAME = "1003"
const USER_WITHDRAWAL_SEARCH_DESC = "1004"

const INIT_CHANNEL_TITLE = "의 생방송입니다"

// Live ----->
const live_KEY = "live:"
const LIVE_KEY_EXPIRE_TIME = 86400 // 1일
const LIVE_KEY_LIKE_CNT = "like_cnt"
const LIVE_KEY_VIEW_CNT = "view_cnt"
const LIVE_KEY_LIVE_CNT = "live_cnt"
const LIVE_KEY_JOIN_LIST = "join_list"

// <----- Live

// File UPLOAD
const FILE_UPLOAD_TYPE_VOD = "vod"
const FILE_UPLOAD_TYPE_IMAGE = "image"
const FILE_UPLOAD_SERVICE_TYPE_LIVE_THUMBNAIL = "live_thumbnail"
const FILE_UPLOAD_SERVICE_TYPE_USER_PROFILE = "profile"
const FILE_UPLOAD_SERVICE_TYPE_CHANNEL_THUMBNAIL = "channel_thumbnail"

// VOD UPLOAD
const VOD_UPLOAD_STATUS_STANDBY = "1001"
const VOD_UPLOAD_STATUS_FINISH = "1002"
const VOD_UPLOAD_SERVICE_TYPE_LIVE = "live"
const VOD_UPLOAD_SERVICE_TYPE_VOD = "vod"

// 리스트 정렬기준
const LIST_ORDER_CODE_LATEST = "01"
const LIST_ORDER_CODE_POPULARITY = "02"
const LIST_ORDER_CODE_LIKE = "03"

// VOD Like Type
const VOD_LIKE_TYPE_LIKE = "01"
const VOD_LIKE_TYPE_DISLIKE = "02"

// VOD Comment Like Type
const VOD_COMMENT_LIKE_TYPE_LIKE = "01"
const VOD_COMMENT_LIKE_TYPE_DISLIKE = "02"

// 신고 상태
const REPORT_ILLEGALITY_STATUS_UNTREATED = "1001"
const REPORT_ILLEGALITY_STATUS_COMPLETE = "1002"
const REPORT_ILLEGALITY_STATUS_NOPROBLEM = "1003"

// 사용자
const USER_STATUS_NORMALITY = "01"
const USER_STATUS_WITHDRAWAL = "02"
const USER_STATUS_BLACKLIST_1 = "03"
const USER_STATUS_BLACKLIST_2 = "04"
const USER_STATUS_BLACKLIST_3 = "05"
const USER_STATUS_BLACKLIST_4 = "06"
const USER_STATUS_BLACKLIST_5 = "07"

// MakeLiveKey ...
func MakeLiveKey(liveSeq uint64) string {
	return live_KEY + strconv.FormatUint(liveSeq, 10)
}

// MakeEmailJoinAuthCodeKey ...
func MakeEmailJoinAuthCodeKey(userID string) string {
	return email_JOIN_AUTH_CODE_KEY + userID
}

// GetEmailJoinAuthCodeKeyExpire ...
func GetEmailJoinAuthCodeKeyExpire() int {
	return email_JOIN_AUTH_CODE_KEY_EXPIRE_TIME
}

// MakeUserJoinInfoProviderUserKey ...
func MakeUserJoinInfoProviderUserKey(key string) string {
	return user_JOIN_INFO_PROVIDER_USER_KEY + key
}

// GetUserJoinInfoProviderUserKeyExpire ...
func GetUserJoinInfoProviderUserKeyExpire() int {
	return user_JOIN_INFO_REDIS_EXPIRE_TIME
}

// GetUserAuthCodeName ...
func GetUserAuthCodeName(key string) string {
	return user_AUTH_CODE[key]
}
