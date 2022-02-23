package triagebot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type payloadKey struct{}

func parseZulipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(resp, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		//nolint:exhaustivestruct
		payload := ZulipPayload{}

		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			response := fmt.Sprintf("could not parse request body: %s", err.Error())
			http.Error(resp, response, http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(req.Context(), payloadKey{}, payload)

		next.ServeHTTP(resp, req.WithContext(ctx))
	})
}

func parseScheduleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(resp, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		//nolint:exhaustivestruct
		payload := SchedulePayload{}
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			response := fmt.Sprintf("could not parse request body: %s", err.Error())
			http.Error(resp, response, http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(req.Context(), payloadKey{}, payload)

		next.ServeHTTP(resp, req.WithContext(ctx))
	})
}
