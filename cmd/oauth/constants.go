package oauth

import "hirelin-auth/internal/logger"

var Providers = providerStruct{
	GOOGLE: "google",
	GITHUB: "github",
}

var Environment = environmentStruct{
	clientUrl:   logger.GetEnv("CLIENT_URL"),
	serverUrl:   logger.GetEnv("SERVER_URL"),
	signInRoute: "/auth/signin",
	logOutRoute: "/auth/logout",
	errorRoute:  "/auth/error",
}

var ProviderEndpoints = map[string]providerEndPointType{
	"google": {
		AuthUri:     "https://accounts.google.com/o/oauth2/auth",
		TokenUri:    "https://oauth2.googleapis.com/token",
		UserInfoUri: "https://www.googleapis.com/oauth2/v3/userinfo",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"openid"},
	},
	//	TODO: Add endpoints for others
}
