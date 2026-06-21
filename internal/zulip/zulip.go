package zulip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const timeout time.Duration = 5 * time.Second

var (
	errInvalidZulipUser = errors.New("zulip response did not include a valid user")
	errNilZulipResponse = errors.New("execute zulip request: nil response")
)

type User struct {
	UserID   int    `json:"user_id"`
	FullName string `json:"full_name"`
}

type Payload struct {
	User User `json:"user"`
}

func (user User) Tag() string {
	return "@**" + user.FullName + "|" + strconv.Itoa(user.UserID) + "**"
}

// UserByEmail get a Zulip user by their email.
func UserByEmail(ctx context.Context, email string) (string, error) {
	//nolint:exhaustruct
	apiURL := url.URL{
		User:   url.UserPassword(os.Getenv("ZULIP_BOT_MAIL"), os.Getenv("ZULIP_BOT_APIKEY")),
		Host:   "reload.zulipchat.com",
		Path:   "/api/v1/users/" + email,
		Scheme: "https",
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create zulip request: %w", err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("execute zulip request: %w", err)
	}

	if response == nil {
		return "", fmt.Errorf("%w", errNilZulipResponse)
	}

	defer response.Body.Close()

	var payload Payload

	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return "", fmt.Errorf("decode zulip response: %w", err)
	}

	if payload.User.UserID == 0 {
		return "", fmt.Errorf("%w", errInvalidZulipUser)
	}

	return payload.User.Tag(), nil
}

// ThumbsUp on a message.
func ThumbsUp(ctx context.Context, messageID int) {
	//nolint:exhaustruct
	apiURL := url.URL{
		User:   url.UserPassword(os.Getenv("ZULIP_BOT_MAIL"), os.Getenv("ZULIP_BOT_APIKEY")),
		Host:   "reload.zulipchat.com",
		Path:   fmt.Sprintf("/api/v1/messages/%d/reactions", messageID),
		Scheme: "https",
	}

	payload := url.Values{}
	payload.Set("emoji_name", "+1")

	ctx, cancel := context.WithTimeout(ctx, timeout)
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
