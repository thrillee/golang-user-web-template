package oauth

import (
	"context"
	"fmt"
)

type HandleOAuth func(context.Context, string) (*OAuthUser, error)

const (
	FACEBOOK_PROVIDER string = "facebook"
	DISCORD_PROVIDER  string = "discord"
	GOOGLE_PROVIDER   string = "google"
)

var oAuthFactory = map[string]HandleOAuth{
	GOOGLE_PROVIDER:   GetGoogleUserData,
	DISCORD_PROVIDER:  GetDiscordUserData,
	FACEBOOK_PROVIDER: GetUserInfoFromFacebook,
}

func HandleOAuthLogin(ctx context.Context, provider, code string) (*OAuthUser, error) {
	fc, ok := oAuthFactory[provider]
	if !ok {
		return nil, fmt.Errorf("OAuth Provider %v not implemented", provider)
	}

	return fc(ctx, code)
}
