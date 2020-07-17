package zulip

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const timeout time.Duration = 5 * time.Second

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

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL.String(), strings.NewReader(payload.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer response.Body.Close()
}
