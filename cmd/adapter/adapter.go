package adapter

import (
	"context"
	"net/http"
	"time"

	"hirelin-auth/cmd/oauth"
	"hirelin-auth/internal/logger"
	"hirelin-auth/internal/server"
	dbType "hirelin-auth/internal/server/db"
	"hirelin-auth/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

func OauthSqlcAdapter() func(oauth.AdapterParams) {
	return func(params oauth.AdapterParams) {
		db := server.GetDB()
		clientUrl := logger.GetEnv("CLIENT_URL")
		session, err := utils.GenerateSession()

		if err != nil {
			http.Redirect(params.Response, params.Request, clientUrl+"/auth/error?error="+err.Error(), http.StatusFound)
			return
		}

		// find existing user
		user, err := db.GetUserByEmail(context.Background(), params.UserData.Email)
		var errStr string
		if err != nil {
			errStr = err.Error()
			if errStr != "no rows in result set" {
				http.Redirect(params.Response, params.Request, clientUrl+"/auth/error?error="+errStr, http.StatusFound)
				return
			}
		}

		// user not exists = account creation
		if errStr == "no rows in result set" {
			dbReturns, err := db.NewOAuthUserTransaction(context.Background(), dbType.NewOAuthUserTransactionParams{
				Email:             params.UserData.Email,
				Name:              pgtype.Text{String: params.UserData.Name, Valid: true},
				Image:             params.UserData.Picture,
				EmailVerified:     pgtype.Timestamp{Time: time.Now(), Valid: true},
				Provider:          params.Provider,
				ProviderAccountID: params.UserData.Sub,
				AccessToken:       pgtype.Text{String: params.TokenData.AccessToken, Valid: true},
				RefreshToken:      pgtype.Text{String: params.TokenData.RefreshToken, Valid: true},
				ExpiresAt:         pgtype.Int4{Int32: int32(params.TokenData.ExpiresIn), Valid: true},
				TokenType:         pgtype.Text{String: params.TokenData.TokenType, Valid: true},
				Scope:             pgtype.Text{String: params.TokenData.Scope, Valid: true},
				IDToken:           pgtype.Text{String: params.TokenData.IDToken, Valid: true},
				SessionState:      pgtype.Text{String: params.Code, Valid: true},
			})
			if err != nil {
				http.Redirect(params.Response, params.Request, clientUrl+"/auth/error?error="+err.Error(), http.StatusFound)
				return
			}

			user = dbType.User{
				ID: dbReturns.UserID,
			}
		} else {
			// account not exists = account linking
			_, err := db.GetAccountByProviderId(context.Background(), dbType.GetAccountByProviderIdParams{
				Provider:          params.Provider,
				ProviderAccountID: params.UserData.Sub,
			})
			if err != nil {
				errStr = err.Error()
				if errStr != "no rows in result set" {
					http.Redirect(params.Response, params.Request, clientUrl+"/auth/error?error="+errStr, http.StatusFound)
					return
				}
			}

			err = nil
			if errStr == "no rows in result set" {
				_, err = db.CreateOAuthAccount(context.Background(), dbType.CreateOAuthAccountParams{
					UserID:            user.ID,
					Provider:          params.Provider,
					ProviderAccountID: params.UserData.Sub,
					AccessToken:       pgtype.Text{String: params.TokenData.AccessToken, Valid: true},
					RefreshToken:      pgtype.Text{String: params.TokenData.RefreshToken, Valid: true},
					ExpiresAt:         pgtype.Int4{Int32: int32(params.TokenData.ExpiresIn), Valid: true},
					TokenType:         pgtype.Text{String: params.TokenData.TokenType, Valid: true},
					Scope:             pgtype.Text{String: params.TokenData.Scope, Valid: true},
					IDToken:           pgtype.Text{String: params.TokenData.IDToken, Valid: true},
					SessionState:      pgtype.Text{String: params.Code, Valid: true},
				})
			} else {
				_, err = db.UpdateOAuthAccount(context.Background(), dbType.UpdateOAuthAccountParams{
					ProviderAccountID: params.UserData.Sub,
					AccessToken:       pgtype.Text{String: params.TokenData.AccessToken, Valid: true},
					RefreshToken:      pgtype.Text{String: params.TokenData.RefreshToken, Valid: true},
					ExpiresAt:         pgtype.Int4{Int32: int32(params.TokenData.ExpiresIn), Valid: true},
					TokenType:         pgtype.Text{String: params.TokenData.TokenType, Valid: true},
					Scope:             pgtype.Text{String: params.TokenData.Scope, Valid: true},
					IDToken:           pgtype.Text{String: params.TokenData.IDToken, Valid: true},
					SessionState:      pgtype.Text{String: params.Code, Valid: true},
					UserID:            user.ID,
					Provider:          params.Provider,
				})
			}
			if err != nil {
				http.Redirect(params.Response, params.Request, clientUrl+"/auth/error?error="+err.Error(), http.StatusFound)
				return
			}
		}
		// create session
		_, err = db.CreateSession(context.Background(), dbType.CreateSessionParams{
			SessionToken: session.SessionId.String(),
			RefreshToken: session.RefreshJwt,
			UserID:       user.ID,
			ExpiresAt:    pgtype.Timestamp{Time: session.Expiry, Valid: true},
		})

		if err != nil {
			http.Redirect(params.Response, params.Request, clientUrl+"/auth/error?error="+err.Error(), http.StatusFound)
			return
		}

		redirect := params.State.Redirect
		if redirect == "" {
			redirect = "/"
		}

		http.SetCookie(params.Response, &http.Cookie{
			Name:     utils.SESSION_TOKEN_NAME,
			Value:    session.SessionId.String(),
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(time.Until(session.Expiry).Seconds()),
		})
		http.SetCookie(params.Response, &http.Cookie{
			Name:     utils.REFRESH_JWT_NAME,
			Value:    session.RefreshJwt,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(utils.REFRESH_LIFETIME),
		})

		http.Redirect(params.Response, params.Request, clientUrl+redirect, http.StatusFound)
	}
}
