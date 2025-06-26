package oauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"hirelin-auth/internal/logger"
)

var useProviders []providerWithConfig
var adapter func(AdapterParams)
var addons func(http.HandlerFunc) http.HandlerFunc

// SelectProviders sets the OAuth providers to be used for authentication.
//
// pass an adapter function to perform actions after successful authentication and fetching user data.
//
// Example usage:
//
//	oauth.SelectProviders(nil, oauth.Providers.GOOGLE, oauth.Providers.GITHUB)
func SelectProviders(a func(AdapterParams), providers ...string) {
	adapter = a
	for _, provider := range providers {
		newProvider := providerWithConfig{
			name: strings.ToLower(provider),
			config: ProviderConfig{
				clientId:     logger.GetEnv(strings.ToUpper(provider) + "_CLIENT_ID"),
				clientSecret: logger.GetEnv(strings.ToUpper(provider) + "_CLIENT_SECRET"),
			},
		}

		useProviders = append(useProviders, newProvider)
	}
}

// WithOAuth Adds the OAuth sign-in handler to the provided ServeMux.
//
// Example usage:
//
//	oauth.SelectProviders(oauth.Providers.GOOGLE)
//	oauth.WithOAuth(mux, middlewares...)
func WithOAuth(mux *http.ServeMux, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	// merge the middlewares
	signInHandler, callbackHandler := oAuthSignIn, handleOauthCallback
	for i := len(middlewares) - 1; i >= 0; i-- {
		signInHandler, callbackHandler = middlewares[i](signInHandler), middlewares[i](callbackHandler)
	}
	mux.HandleFunc("/api/auth/oauth/signin", signInHandler)

	// Add providers
	for _, provider := range useProviders {
		mux.HandleFunc("/api/auth/callback/"+strings.ToLower(provider.name), addons(callbackHandler))
	}
}

func oAuthSignIn(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	redirect := r.URL.Query().Get("redirect")

	if redirect == "" {
		redirect = "/"
	}

	if provider == "" {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=invalid_provider", http.StatusFound)
		return
	}

	var providerConfig providerWithConfig
	flag := true
	for _, p := range useProviders {
		if provider == p.name {
			providerConfig = p
			flag = false
			break
		}
	}
	if flag {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=invalid_provider", http.StatusFound)
	}

	state, err := GenerateStateHash(StateType{
		Redirect:  redirect,
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Provider:  providerConfig.name,
	})
	if err != nil {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=error_generating_state", http.StatusFound)
		return
	}

	redirectUrl := ProviderEndpoints[providerConfig.name].AuthUri + "?"
	redirectUrl += "client_id=" + providerConfig.config.clientId
	redirectUrl += "&redirect_uri=" + Environment.serverUrl + "/api/auth/callback/" + providerConfig.name
	redirectUrl += "&response_type=code&access_type=offline&prompt=consent"
	redirectUrl += "&state=" + state
	redirectUrl += "&scope=" + strings.Join(ProviderEndpoints[providerConfig.name].Scopes, " ")

	http.Redirect(w, r, redirectUrl, http.StatusFound)
}

func handleOauthCallback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	code := query.Get("code")
	state := query.Get("state")
	authError := query.Get("error")

	if authError != "" {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error="+authError, http.StatusFound)
		return
	}
	if code == "" || state == "" {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=missing_code_or_state", http.StatusFound)
		return
	}

	stateData, err := GetDataFromStateHash(state)
	if err != nil || stateData.ExpiresAt < time.Now().Unix() {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=invalid_state", http.StatusFound)
		return
	}

	tokenData, err := fetchTokenFromCode(code, stateData.Provider)
	if err != nil {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=failed_to_fetch_token", http.StatusFound)
		return
	}

	userData, err := fetchUserFromToken(tokenData.AccessToken, stateData.Provider)
	if err != nil {
		http.Redirect(w, r, Environment.clientUrl+Environment.errorRoute+"?error=failed_to_fetch_user", http.StatusFound)
		return
	}

	if adapter == nil {
		http.Redirect(w, r, Environment.clientUrl+stateData.Redirect, http.StatusFound)
		return
	}

	adapter(AdapterParams{
		TokenData: tokenData,
		UserData:  userData,
		Response:  w,
		Request:   r,
		Provider:  stateData.Provider,
		Code:      code,
		State:     stateData,
	})
}

func fetchTokenFromCode(code string, provider string) (*TokenData, error) {
	var providerConfig providerWithConfig
	flag := true
	for _, p := range useProviders {
		if provider == p.name {
			providerConfig = p
			flag = false
			break
		}
	}
	if flag {
		return nil, errors.New("invalid_provider")
	}

	url := ProviderEndpoints[provider].TokenUri + "?"
	url += "client_id=" + providerConfig.config.clientId
	url += "&client_secret=" + providerConfig.config.clientSecret
	url += "&code=" + code
	url += "&grant_type=authorization_code"
	url += "&redirect_uri=" + Environment.serverUrl + "/api/auth/callback/" + providerConfig.name

	res, err := http.Post(url, "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed_to_fetch_token")
	}

	var tokenData TokenData
	if err := json.NewDecoder(res.Body).Decode(&tokenData); err != nil {
		return nil, err
	}

	return &tokenData, nil
}

func fetchUserFromToken(token string, provider string) (*UserData, error) {
	res, err := http.Get(ProviderEndpoints[provider].UserInfoUri + "?access_token=" + token)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed_to_fetch_user")
	}

	var userData UserData
	if err := json.NewDecoder(res.Body).Decode(&userData); err != nil {
		return nil, err
	}

	return &userData, nil
}

// WithAddons adds extra middleware functions to callback routes.
//
// Note:
//
//  1. These addons will be executed at the end of executing the callback.
//
//  2. Execution order is the same as the order of the addons passed.
func WithAddons(adds ...func(http.HandlerFunc) http.HandlerFunc) {
	addons = func(h http.HandlerFunc) http.HandlerFunc {
		for _, addon := range adds {
			h = addon(h)
		}
		return h
	}
}
