package oauth

import (
	"context"
	"os"

	disgoauth "github.com/realTristan/disgoauth"
)

var dc *disgoauth.Client

func GetDiscardClient() *disgoauth.Client {
	if dc == nil {
		dc = disgoauth.Init(&disgoauth.Client{
			ClientID:     "", // os.Getenv("DISCORD_CLIENT_ID"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
			RedirectURI:  os.Getenv("DISCORD_REDIRECT_SECRET"),
			Scopes:       []string{disgoauth.ScopeIdentify},
		})
	}

	return dc
}

func GetDiscordUserData(ctx context.Context, codeFromURLParamaters string) (*OAuthUser, error) {
	accessToken, err := GetDiscardClient().GetOnlyAccessToken(codeFromURLParamaters)
	if err != nil {
		return nil, err
	}
	userData, err := disgoauth.GetUserData(accessToken)
	if err != nil {
		return nil, err
	}

	email, _ := userData["email"]
	phone, _ := userData["phone"]
	fullName, _ := userData["name"]
	imageURL, _ := userData["image"]

	return &OAuthUser{
		Fullname: fullName.(string),
		Email:    email.(string),
		Phone:    phone.(string),
		ImageURL: imageURL.(string),
	}, nil
}
