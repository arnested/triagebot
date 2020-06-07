package triagebot

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value(payloadKey{}).(Payload)
		if token, ok := os.LookupEnv("ZULIP_TOKEN"); !ok || payload.Token != token {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		re := regexp.MustCompile(`@reload\.dk$`)
		if !re.MatchString(payload.Message.SenderEmail) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			response := Response{
				Content: ExternalUserText,
			}
			json.NewEncoder(w).Encode(response)

			return
		}

		next.ServeHTTP(w, r)
	})
}
