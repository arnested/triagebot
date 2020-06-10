package triagebot

import (
	"net/http"
	"os"
)

func authenticationOutgoingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value(payloadKey{}).(ZulipPayload)
		if token, ok := os.LookupEnv("ZULIP_TOKEN"); !ok || payload.Token != token {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}

func authenticationScheduleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := r.Context().Value(payloadKey{}).(SchedulePayload)
		if token, ok := os.LookupEnv("ZULIP_TOKEN"); !ok || payload.Token != token {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		next.ServeHTTP(w, r)
	})
}
