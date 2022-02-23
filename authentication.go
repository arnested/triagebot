package triagebot

import (
	"net/http"
	"os"
)

func authenticationOutgoingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		payload, ok := req.Context().Value(payloadKey{}).(ZulipPayload)
		if !ok {
			http.Error(resp, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		if token, ok := os.LookupEnv("ZULIP_TOKEN"); !ok || payload.Token != token {
			http.Error(resp, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		next.ServeHTTP(resp, req)
	})
}

func authenticationScheduleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		payload, ok := req.Context().Value(payloadKey{}).(SchedulePayload)
		if !ok {
			http.Error(resp, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		if token, ok := os.LookupEnv("ZULIP_TOKEN"); !ok || payload.Token != token {
			http.Error(resp, http.StatusText(http.StatusForbidden), http.StatusForbidden)

			return
		}

		next.ServeHTTP(resp, req)
	})
}
