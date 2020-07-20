
# MCS(Media Control Service)
COSHOP-WAS 미디어 연동 서비스 REST-API 서버
### API Specification
REST-API-MCS-v.1.2.pdf
( https://github.com/teamgrit-lab/mcs/blob/master/rest-api-spec/REST-API-MCS-v.1.2.pdf )

### config.[local/dev/real].json 설정 파일 구조
<pre>
<code>
{
  "domain": "192.168.0.45",             // 도메인 정보
  "address": "192.168.0.45",            // IP주소
  "port": 9200,                         // MCS Servcie Port
  "media_admin_port": 17088,            // janus admin port
  "media_admin_secret": "janusoverlord",// janus admin secret
  "project_name": "mcs",                // project name 
  "admin_info": {                       // admin info  
    "email": "admin@teamgrit.kr"
  },
  "janus_rec_cmd": "/opt/janus/bin/janus-pp-rec",   // janus mjr to mp4, opus format conversion tool path 
  "janus_rec_path": "/tmp/records",                 // storage path to janus videoroom recording   
  "repository_path": "/home/teamgrit/volumes/repo", //  after janus-pp-rec, ffmpeg, final mp4 file saved in this repository   
  "textroom_post": "http://1.214.216.250:19200/textroom",   // janus textroom post url config to save chatting message

  "storage": {
    "rdb": {                                                                        // maria db config
      "url": "grit:grit@tcp(192.168.0.45:3306)/cojam?charset=utf8&parseTime=true",
      "open": 100,
      "idle": 10
    },
    "mongo": {                                      // Not Available in mcs                                                               
      "url": "localhost:27017",
      "pool": 50,
      "db_name": [
        "webrtc-mvp-log",
        "mvp"
      ]
    },
    "conn_redis": {                                 //temporary store to manage janus connections   
      "url": "192.168.0.45:6379",
      "max_idle": 50,
      "max_active": 10000,
      "dbnum": 3
    },
    "redis": {                                      
      "url": "192.168.0.45:6379",
      "max_idle": 50,
      "max_active": 10000,
      "dbnum": 2
    },
    "session_redis": {                // temporary store to manange janus session                 
      "url": "192.168.0.45:6379",
      "max_idle": 50,
      "max_active": 10000,
      "dbnum": 1
    }
  },                                  // janus access config for more than 2 janus servers
                                      // format : {http_access_point}$${websocket_access_point}
  "media_servers": [
    "https://cojam.iptime.org:18089$$wss://cojam.iptime.org:8199"
  ]
}
</code>
</pre>

### Deployment
- cojam$ git pull

- [Dev]
- docker-compose.yml > services:dev:image (버전수정) "cojam/mcs-dev:version"
- Makefile > services:dev:image (버전수정) "cojam/mcs-dev:version"
- make build;make push

- [Pro]


- docker-compose.yml > services:dev:image (버전수정) "cojam/mcs:0.1.1"
- Makefile > services:dev:image (버전수정) "cojam/mcs:0.2.91"
- make build;make push
