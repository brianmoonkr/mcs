package transcode

import (
	"context"
	
	"go.uber.org/zap"
	"github.com/teamgrit-lab/cojam/config"
)

func TranscodeAll() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	use := Process{Logger: logger.Named("process")}
	message := `{"sid", "","path","` + config.CF.Prop.JanusRecPath + `"}`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Callback funtion

	use.Process(message, ctx)
}

//transcode mjr to mp4
func TranscodeBySession(seq string, sid string, subPath string,  cb string) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	fullPath := config.CF.Prop.JanusRecPath + "/" + subPath + "/records"

	use := Process{Logger: logger.Named("process")}
	message := `{"seq", "` + seq + `","sid", "` + sid + `","path","` + fullPath + `", "cb", "` + cb + `"}`
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	use.Process(message, ctx)	
	
}
