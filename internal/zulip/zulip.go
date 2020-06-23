package zulip

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// ThumbsUp on a message.
func ThumbsUp(messageID int) {
	apiURL := url.URL{
		User:   url.UserPassword(os.Getenv("ZULIP_BOT_MAIL"), os.Getenv("ZULIP_BOT_APIKEY")),
		Host:   "reload.zulipchat.com",
		Path:   fmt.Sprintf("/api/v1/messages/%d/reactions", messageID),
		Scheme: "https",
	}

	payload := url.Values{}
	payload.Set("emoji_name", "+1")

	response, err := http.Post(apiURL.String(), "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))

	if err != nil {
		return
	}

	defer response.Body.Close()
}
