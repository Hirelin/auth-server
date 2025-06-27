package oauth

import "net/http"

type providerStruct struct {
	GOOGLE string
	GITHUB string
}

type environmentStruct struct {
	clientUrl   string
	serverUrl   string
	signInRoute string
	logOutRoute string
	errorRoute  string
}

type ProviderConfig struct {
	clientId     string
	clientSecret string
}

type providerWithConfig struct {
	name   string
	config ProviderConfig
}

type providerEndPointType struct {
	AuthUri     string
	TokenUri    string
	UserInfoUri string
	Scopes      []string
}

type StateType struct {
	Redirect   string    `json:"redirect"`
	Provider   string    `json:"provider"`
	ExpiresAt  int64     `json:"expires_at"`  // duration = 10 minutes
	ClientInfo ClintInfo `json:"client_info"` // e.g., "web", "mobile"
}

type ClintInfo struct {
	Browser Software   `json:"browser"`
	Device  DeviceType `json:"device"`
	OS      Software   `json:"os"`
	CPU     string     `json:"cpu"`
}
type Software struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type DeviceType struct {
	Type  string `json:"type"`
	Model string `json:"model"`
}

type TokenData struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

type UserData struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
	EmailVerified bool   `json:"email_verified"`
}

type AdapterParams struct {
	TokenData *TokenData
	UserData  *UserData
	Response  http.ResponseWriter
	Request   *http.Request
	Provider  string
	State     StateType
	Code      string
}
