package triagebot

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		payload, ok := req.Context().Value(payloadKey{}).(ZulipPayload)
		if !ok {
			http.Error(resp, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		re := regexp.MustCompile(`@reload\.dk$`)
		if !re.MatchString(payload.Message.SenderEmail) {
			resp.Header().Set("Content-Type", "application/json; charset=utf-8")

			response := ZulipResponse{
				Content: ExternalUserText,
			}
			_ = json.NewEncoder(resp).Encode(response)

			return
		}

		next.ServeHTTP(resp, req)
	})
}
