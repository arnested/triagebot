package triagebot

import (
	"net/http"

	"arnested.dk/go/triagebot/internal/zulip"
)

func reactMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value(payloadKey{}).(ZulipPayload)
		go zulip.ThumbsUp(payload.Message.ID)

		next.ServeHTTP(w, r)
	})
}
