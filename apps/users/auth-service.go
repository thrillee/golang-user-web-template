package users

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/thrillee/triq/internals/common"
	"github.com/thrillee/triq/internals/schemas"
	"github.com/thrillee/triq/internals/security"
	"github.com/thrillee/triq/oauth"
)

type AuthService interface {
	HandleLogin(context.Context, LoginPayload) (*User, error)
}

var (
	AUTH_EXP           time.Duration = time.Hour * 7
	AUTH_COOKIE_NAME   string        = os.Getenv("AUTH_COOKIE_NAME")
	AUTH_COOKIE_DOMAIN string        = os.Getenv("AUTH_COOKIE_DOMAIN")
)

type AUTH_PROVIDER string

const (
	EMAIL_AUTH   AUTH_PROVIDER = "email"
	PHONE_AUTH   AUTH_PROVIDER = "phone"
	GOOGLE_AUTH  AUTH_PROVIDER = AUTH_PROVIDER(oauth.GOOGLE_PROVIDER)
	DISCORD_AUTH AUTH_PROVIDER = AUTH_PROVIDER(oauth.DISCORD_PROVIDER)
	FACBOOK_AUTH AUTH_PROVIDER = AUTH_PROVIDER(oauth.FACEBOOK_PROVIDER)
)

type loginHandler func(context.Context, *LoginPayload) (*LoginResponse, error)

var authFactory = map[AUTH_PROVIDER]loginHandler{
	EMAIL_AUTH:   localLoginHandler,
	PHONE_AUTH:   localLoginHandler,
	GOOGLE_AUTH:  hanldeOAuthUser,
	DISCORD_AUTH: hanldeOAuthUser,
	FACBOOK_AUTH: hanldeOAuthUser,
}

func hanldeOAuthUser(ctx context.Context, d *LoginPayload) (*LoginResponse, error) {
	if d.AuthToken == "" {
		return nil, fmt.Errorf("Auth token is required for oauth")
	}
	var query User
	oAuthUserResult, err := oauth.HandleOAuthLogin(ctx, string(d.AuthProvider), d.AuthToken)
	if err != nil {
		return nil, err
	}
	var user User
	result := query.Query(ctx).Model(User{Email: oAuthUserResult.Email}).First(&user)
	err = schemas.HandleDBError(result, "Account")
	if err != nil {
		user, err = createOAuthUser(ctx, oAuthUserResult)
	}

	return &LoginResponse{
		MfaRequired: false,
		AuthUser:    &user,
	}, nil
}

func localLoginHandler(ctx context.Context, d *LoginPayload) (*LoginResponse, error) {
	var query User
	if d.AuthProvider == EMAIL_AUTH {
		query = User{Email: d.Username}
	} else {
		query = User{Phone: d.Username}
	}

	var user User
	result := query.Query(ctx).Model(query).First(&user)
	err := schemas.HandleDBError(result, "Account")
	if err != nil {
		return nil, err
	}

	if !security.CheckPasswordHash(d.Password, user.Password) {
		return nil, fmt.Errorf("Incorrect Email or Password")
	}

	return &LoginResponse{
		MfaRequired: false,
		AuthUser:    &user,
	}, nil
}

func (i UserServiceImp) Login(ctx context.Context, d *LoginPayload) (*LoginResponse, error) {
	if err := common.Validate(d); err != nil {
		return nil, err
	}

	handler, ok := authFactory[d.AuthProvider]
	if !ok {
		options := []string{}
		for k := range authFactory {
			options = append(options, string(k))
		}
		return nil, common.NewSessionError(
			common.BAD_REQUEST,
			fmt.Errorf("Authentication provider '%s' not found: Options are (%v)", d.AuthProvider, strings.Join(options, ", ")))
	}

	res, err := handler(ctx, d)
	if err != nil {
		log.Errorf("Login Error: %v", err)
		return nil, common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Invalid Credentials"))
	}

	token, err := createToken(tokenProps{
		email:      res.AuthUser.Email,
		accountRef: res.AuthUser.AccountRef,
		exp:        time.Now().Add(AUTH_EXP),
	})
	if err != nil {
		return nil, common.NewSessionError(common.INTERNVAL_SERVER_ERROR, err)
	}

	cookie := LoginCookie{
		Name:     AUTH_COOKIE_NAME,
		Value:    token,
		Path:     "/",
		Domain:   AUTH_COOKIE_DOMAIN,
		MaxAge:   int(time.Now().Add(AUTH_EXP).Unix()),
		Expires:  time.Now().Add(AUTH_EXP),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "lax",
	}

	res.Cookie = cookie

	user := res.AuthUser
	user.LastLogin = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	user.QueryRepo().Edit(ctx, user.ID, user)
	return res, nil
}

func (i UserServiceImp) GetAuthUser(ctx context.Context, authToken string) (*User, error) {
	claims, err := getTokenClaims(authToken)
	if err != nil {
		return nil, common.NewSessionError(common.INTERNVAL_SERVER_ERROR, err)
	}

	accountRef := claims["accountRef"].(string)
	var user User
	result := i.model.Query(ctx).Model(User{AccountRef: accountRef}).First(&user)
	err = schemas.HandleDBError(result, "Account")
	if err != nil {
		return nil, common.NewSessionError(common.INTERNVAL_SERVER_ERROR, err)
	}

	return &user, nil
}
