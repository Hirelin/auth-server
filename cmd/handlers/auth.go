package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"hirelin-auth/cmd/types"
	"hirelin-auth/internal/server"
	"hirelin-auth/internal/utils"
)

func GetSessionUser(w http.ResponseWriter, r *http.Request) {
	sessionData := r.Context().Value(types.SessionKey).(types.SessionData)

	if sessionData.Valid {
		db := server.GetDB()

		user, err := db.GetSessionData(context.Background(), sessionData.Token.String())

		if err != nil {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
			return
		}

		var recruiter any
		if user.RecruiterID.String() != "" {
			recruiter = map[string]interface{}{
				"id":           user.RecruiterID.String(),
				"name":         user.RecruiterName,
				"organization": user.Organization,
				"phone":        user.Phone,
				"address":      user.RecruiterAddress,
				"position":     user.Position,
			}
		} else {
			recruiter = nil
		}

		response := map[string]interface{}{
			"message": "User retrieved successfully",
			"status":  "success",
			"session": map[string]interface{}{
				"user": map[string]interface{}{
					"id":        user.ID.String(),
					"name":      user.Name,
					"email":     user.Email,
					"image":     user.Image,
					"recruiter": recruiter,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func LogoutAPI(w http.ResponseWriter, r *http.Request) {
	// session, _ := r.Cookie(utils.SESSION_TOKEN_NAME)
	session := r.Context().Value(types.SessionKey).(types.SessionData)

	http.SetCookie(w, &http.Cookie{
		Name:     utils.SESSION_TOKEN_NAME,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     utils.REFRESH_JWT_NAME,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	if session.Valid {
		db := server.GetDB()
		db.DeleteSession(context.Background(), session.Token.String())
	}

	redirect := r.URL.Query().Get("redirect")

	if redirect == "" {
		redirect = "/"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message":  "Logged out successfully",
		"status":   "success",
		"redirect": redirect,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}
