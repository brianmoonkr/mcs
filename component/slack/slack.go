package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/teamgrit-lab/cojam/component/tglog"
	"github.com/teamgrit-lab/cojam/config"
)

// Slack ...
type Slack struct {
	UserName string `json:"username"`
	//IconEmoji string `json:"icon_emoji"`
	//IconURL   string `json:"icon_url"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

// Send ...
func (slack *Slack) Send() {
	execMode := config.CF.ExecutionMode
	if execMode == "local" {
		return
	}

	cfSlack := config.CF.Prop.API.Slack
	slack.UserName = "cojam"
	slack.Text = fmt.Sprintf("[%s] %s", execMode, slack.Text)
	params, _ := json.Marshal(slack)

	req, err := http.NewRequest(
		"POST",
		cfSlack.InComingURL,
		bytes.NewBuffer([]byte(params)),
	)
	if err != nil {
		tglog.Logger.Errorf("Slack Send http.NewRequest: %+v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		tglog.Logger.Errorf("Slack Send client.Do: %+v", err)
	}
	defer resp.Body.Close()
}
