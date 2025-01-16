package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	GoogleUserInfoEndpoint = "https://www.googleapis.com/oauth2/v2/userinfo"
)

type OauthClient struct {
	Google *oauth2.Config
}

func NewOauthClient(cnf *Config) *OauthClient {
	return &OauthClient{
		Google: initGoogleOauth(cnf),
	}
}

func initGoogleOauth(cnf *Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cnf.Env.GetString("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: cnf.Env.GetString("GOOGLE_OAUTH_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:3000/api/auth/google/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}
}
