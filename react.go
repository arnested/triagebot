package triagebot

import (
	"net/http"

	"arnested.dk/go/triagebot/internal/zulip"
)

func reactMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		payload, ok := req.Context().Value(payloadKey{}).(ZulipPayload)
		if !ok {
			http.Error(resp, http.StatusText(http.StatusNoContent), http.StatusNoContent)

			return
		}

		go zulip.ThumbsUp(payload.Message.ID)

		next.ServeHTTP(resp, req)
	})
}
