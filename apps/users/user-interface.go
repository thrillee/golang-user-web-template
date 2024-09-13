package users

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/thrillee/triq/apps/otp"
	"github.com/thrillee/triq/internals/schemas"
)

type LoginNextAction string

const (
	GOOGLE_AUTH_COMPLETE LoginNextAction = "GOOGLE_AUTH_COMPLETE"
)

type LoginCookie struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Path        string    `json:"path"`
	Domain      string    `json:"domain"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	SameSite    string    `json:"same_site"`
	SessionOnly bool      `json:"session_only"`
}

type LoginResponse struct {
	MfaRequired bool
	AuthUser    *User
	NextAction  LoginNextAction
	Cookie      LoginCookie
}

type VerifyAccountPayload struct {
	AccountID string            `json:"account_id" validate:"required"`
	OTP       string            `json:"otp"`
	Password  string            `json:"password"`
	Target    otp.OTPTargetType `json:"target"`
}

type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type NewUserPayload struct {
	Email    string           `json:"email" validate:"required,email"`
	Phone    string           `json:"phone" validate:"required"`
	Password string           `json:"password"`
	Fullname schemas.MyString `json:"fullname"`
	Username string           `json:"username" validate:"required"`
}

type LoginPayload struct {
	AuthProvider AUTH_PROVIDER `json:"auth_provider" validate:"required"`
	Username     string        `json:"username" validate:"required"`
	Password     string        `json:"password" validate:"required"`
	AuthToken    string        `json:"auth_token"`
}

type OTPResendPayload struct {
	TargetRef string `json:"target_ref"`
}

type UserService interface {
	CreateUser(context.Context, *NewUserPayload) (*User, error)
	VerifyAccount(context.Context, *VerifyAccountPayload) (*User, error)
	EditUser(context.Context, string, *NewUserPayload) (*User, error)
	ForgotPassword(context.Context, *VerifyAccountPayload) error
	ResetPassword(context.Context, *VerifyAccountPayload) error
	Login(context.Context, *LoginPayload) (*LoginResponse, error)
	ChangePassword(context.Context, *User, *ChangePasswordPayload) error
	UploadDisplayPicutre(context.Context, *User, *multipart.FileHeader) error
	GetAuthUser(context.Context, string) (*User, error)
}
