package triagebot

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value(payloadKey{}).(ZulipPayload)

		re := regexp.MustCompile(`@reload\.dk$`)
		if !re.MatchString(payload.Message.SenderEmail) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			response := ZulipResponse{
				Content: ExternalUserText,
			}
			json.NewEncoder(w).Encode(response)

			return
		}

		next.ServeHTTP(w, r)
	})
}
