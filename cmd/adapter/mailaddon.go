package adapter

import (
	"bytes"
	"encoding/json"
	"hirelin-auth/cmd/oauth"
	"hirelin-auth/cmd/types"
	"hirelin-auth/internal/logger"
	"net/http"
)

func AuthMailAddon() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			state := r.Context().Value(types.StateKey).(oauth.StateType)
			user := r.Context().Value(types.OAuthUserKey).(oauth.UserData)

			payload := map[string]interface{}{
				"client_info": state.ClientInfo,
				"user":        user,
			}

			payloadBytes, err := json.Marshal(payload)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			res, err := http.Post(
				logger.GetEnv("SERVER_URL")+"/api/notification/authmail",
				"application/json",
				bytes.NewBuffer(payloadBytes),
			)

			if err != nil {
				logger.Logger.Println("Error sending auth mail notification:", err)
			} else {
				defer res.Body.Close()
			}

			next.ServeHTTP(w, r)
		})
	}
}
