package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"hirelin-auth/internal/server"
	"hirelin-auth/internal/utils"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgtype"

	types "hirelin-auth/cmd/types"
	database "hirelin-auth/internal/server/db"
)

// GlobalMiddleWare helps in refreshing token upon expiry
func GlobalMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get session ID
		sessionToken, err := r.Cookie(utils.SESSION_TOKEN_NAME)
		var sessionData types.SessionData

		if err != nil {
			// Get refresh token when session id is not present
			refreshToken, err := r.Cookie(utils.REFRESH_JWT_NAME)
			var claims jwt.MapClaims

			if err == nil {
				claims, err = utils.ParseJWT(refreshToken.Value)
				if err != nil {
					fmt.Println("Error parsing refresh token:", err)
					next.ServeHTTP(w, r)
					return
				}
			}

			//  Refresh session if refresh token is present
			if err == nil {
				session, err := utils.GenerateSession()

				if err != nil {
					fmt.Println("Error generating session:", err)
					next.ServeHTTP(w, r)
					return
				}

				db := server.GetDB()

				_, err = db.UpdateSession(context.Background(), database.UpdateSessionParams{
					SessionToken:   session.SessionId.String(),
					RefreshToken:   session.RefreshJwt,
					ExpiresAt:      pgtype.Timestamp{Time: session.Expiry, Valid: true},
					SessionToken_2: claims["session_id"].(string),
				})

				if err == nil {
					http.SetCookie(w, &http.Cookie{
						Name:     utils.SESSION_TOKEN_NAME,
						Value:    session.SessionId.String(),
						Path:     "/",
						HttpOnly: true,
						Secure:   true,
						MaxAge:   int(time.Until(session.Expiry).Seconds()),
						SameSite: http.SameSiteStrictMode,
					})
					http.SetCookie(w, &http.Cookie{
						Name:     utils.REFRESH_JWT_NAME,
						Value:    session.RefreshJwt,
						Path:     "/",
						HttpOnly: true,
						Secure:   true,
						MaxAge:   int(utils.REFRESH_LIFETIME),
						SameSite: http.SameSiteStrictMode,
					})
					session_pgtype, _ := utils.StringToPgtypeUUID(session.SessionId.String())
					sessionData.Token = pgtype.UUID{Bytes: session_pgtype.Bytes, Valid: true}
					sessionData.Valid = true
				} else {
					sessionData.Token = pgtype.UUID{Valid: false}
					sessionData.Valid = false
				}
			}
		} else {
			session_pgtype, _ := utils.StringToPgtypeUUID(sessionToken.Value)
			sessionData.Token = session_pgtype
			sessionData.Valid = true
		}

		r = r.WithContext(context.WithValue(r.Context(), types.SessionKey, sessionData))
		r = r.WithContext(context.WithValue(r.Context(), types.UserKey, types.UserData{Valid: false}))

		next.ServeHTTP(w, r)
	})
}

// ProtectedMiddleware is used to protect routes that require authentication and adds user data to the request Context
func ProtectedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return GlobalMiddleWare(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionData := r.Context().Value(types.SessionKey).(types.SessionData)

		if !sessionData.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		db := server.GetDB()
		user, err := db.GetSessionUser(context.Background(), sessionData.Token.String())

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userId_pgtype, _ := utils.StringToPgtypeUUID(user.ID.String())

		r = r.WithContext(context.WithValue(r.Context(), types.UserKey, types.UserData{
			Valid:    true,
			ID:       userId_pgtype,
			Name:     user.Name,
			Email:    user.Email,
			UserName: user.Username,
		}))

		next.ServeHTTP(w, r)
	}))
}
