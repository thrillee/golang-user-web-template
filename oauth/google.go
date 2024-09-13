package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
}

var AppConfig Config

func GoogleConfig() oauth2.Config {
	AppConfig.GoogleLoginConfig = oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return AppConfig.GoogleLoginConfig
}

func InitUser(rand string) string {
	return AppConfig.GoogleLoginConfig.AuthCodeURL(rand)
}

func GetGoogleUserData(ctx context.Context, code string) (*OAuthUser, error) {
	token, err := AppConfig.GoogleLoginConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, err
	}

	userDataByte, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(userDataByte, &jsonData)
	if err != nil {
		return nil, err
	}

	email, _ := jsonData["email"]
	name, _ := jsonData["name"]
	picture, _ := jsonData["picture"]
	return &OAuthUser{
		Fullname: name.(string),
		Email:    email.(string),
		ImageURL: picture.(string),
	}, nil
}
