package users

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/thrillee/triq/apps/otp"
	"github.com/thrillee/triq/internals/common"
	"github.com/thrillee/triq/internals/schemas"
	"github.com/thrillee/triq/oauth"
)

func (i UserServiceImp) ResetPassword(ctx context.Context, d *VerifyAccountPayload) error {
	if err := common.Validate(d); err != nil {
		return err
	}

	var user User
	result := i.model.Query(ctx).Model(User{AccountRef: d.AccountID}).First(&user)
	err := schemas.HandleDBError(result, "Account")
	if err != nil {
		return common.NewSessionError(common.BAD_REQUEST, err)
	}

	ok, err := otpService.VerifyOTP(ctx, otp.VerifyOTPPayload{
		OTPCode:   d.OTP,
		TargetRef: d.AccountID,
		Target:    d.Target,
	})
	if err != nil || !ok {
		return common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Verification Failed: %v", err))
	}

	user.IsVerified = true
	user.SetPassword(ctx, d.Password)

	go publishEvent(ctx, USER_RESET_PASWORD, &user)

	return nil
}

func (i UserServiceImp) ForgotPassword(ctx context.Context, d *VerifyAccountPayload) error {
	if err := common.Validate(d); err != nil {
		return err
	}

	var user User
	result := i.model.Query(ctx).Model(User{Email: d.AccountID}).First(&user)
	err := schemas.HandleDBError(result, "Account")
	if err != nil {
		return common.NewSessionError(common.BAD_REQUEST, err)
	}

	go publishEvent(ctx, USER_FORGOT_PASSWORD, &user)
	return nil
}

func (i UserServiceImp) VerifyAccount(ctx context.Context, d *VerifyAccountPayload) (*User, error) {
	if err := common.Validate(d); err != nil {
		return nil, err
	}

	var user User
	result := i.model.Query(ctx).Model(User{AccountRef: d.AccountID}).First(&user)
	err := schemas.HandleDBError(result, "Account")
	if err != nil {
		return nil, common.NewSessionError(common.BAD_REQUEST, err)
	}

	ok, err := otpService.VerifyOTP(ctx, otp.VerifyOTPPayload{
		OTPCode:   d.OTP,
		TargetRef: d.AccountID,
		Target:    d.Target,
	})
	if err != nil || !ok {
		return nil, common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Verification Failed: %v", err))
	}

	user.IsVerified = true
	user.QueryRepo().Edit(ctx, user.ID, &user)

	go publishEvent(ctx, USER_VERIFIED, &user)

	return &user, nil
}

func (i UserServiceImp) CreateUser(ctx context.Context, u *NewUserPayload) (*User, error) {
	if err := common.Validate(u); err != nil {
		return nil, err
	}

	user := User{
		FullName:   u.Fullname,
		Username:   u.Username,
		Email:      u.Email,
		Phone:      u.Phone,
		Active:     true,
		IsVerified: false,
		AccountRef: uuid.NewString(),
	}

	if user.EmailExists(ctx, u.Email) {
		return nil, common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Email already exits"))
	}

	if user.PhoneExists(ctx, u.Phone) {
		return nil, common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Phone number already exits"))
	}

	_, err := i.repo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	user.SetPassword(ctx, u.Password)

	go publishEvent(ctx, NEW_USER, &user)

	return &user, nil
}

func createOAuthUser(ctx context.Context, u *oauth.OAuthUser) (User, error) {
	user := User{
		FullName:   schemas.MyString(u.Fullname),
		Username:   u.Email,
		Email:      u.Email,
		Phone:      u.Phone,
		Active:     true,
		IsVerified: false,
		AccountRef: uuid.NewString(),
	}

	_, err := user.QueryRepo().Create(ctx, &user)

	user.SetPassword(ctx, uuid.NewString())

	go publishEvent(ctx, USER_VERIFIED, &user)

	return user, common.NewSessionError(common.BAD_REQUEST, err)
}
