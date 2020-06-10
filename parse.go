package triagebot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type payloadKey struct{}

func parseZulipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		payload := ZulipPayload{}
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			response := fmt.Sprintf("could not parse request body: %s", err.Error())
			http.Error(w, response, http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), payloadKey{}, payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseScheduleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		payload := SchedulePayload{}
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			response := fmt.Sprintf("could not parse request body: %s", err.Error())
			http.Error(w, response, http.StatusBadRequest)

			return
		}

		ctx := context.WithValue(r.Context(), payloadKey{}, payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
