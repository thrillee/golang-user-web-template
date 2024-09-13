package users

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thrillee/triq/apps/otp"
	"github.com/thrillee/triq/internals/common"
	"github.com/thrillee/triq/internals/schemas"
)

type eventHandler func(context.Context, *User) error

type userEventType string

const (
	NEW_USER              userEventType = "new_user"
	USER_VERIFIED         userEventType = "user_verified"
	USER_RESET_PASWORD    userEventType = "user_reset_password"
	USER_FORGOT_PASSWORD  userEventType = "user_forgot_password"
	USER_CHANGED_PASSWORD userEventType = "user_changed_password"
)

var userEventFactory = map[userEventType][]eventHandler{
	NEW_USER:              {handleRequestAccountVerify},
	USER_VERIFIED:         {},
	USER_RESET_PASWORD:    {},
	USER_CHANGED_PASSWORD: {},
	USER_FORGOT_PASSWORD:  {handleRequestForgotPasswordVerify},
}

var otpService = otp.NewOTPService()

func handleVerifyAccountResend(ctx context.Context, d *OTPResendPayload) error {
	if err := common.Validate(d); err != nil {
		return err
	}

	query := User{AccountRef: d.TargetRef}

	var user User
	result := query.Query(ctx).Model(query).First(&user)
	err := schemas.HandleDBError(result, "Account")
	if err != nil {
		return err
	}

	if user.IsVerified {
		return common.NewSessionError(401, fmt.Errorf("Account is already verified"))
	}

	return handleRequestAccountVerify(ctx, &user)
}

func handleRequestAccountVerify(ctx context.Context, u *User) error {
	return handleRequestVerify(ctx, otp.OTP_VERIFY_ACCOUNT, u)
}

func handleRequestForgotPasswordVerify(ctx context.Context, u *User) error {
	return handleRequestVerify(ctx, otp.OTP_RESET_PASSWORD, u)
}

func handleRequestVerify(ctx context.Context, targetType otp.OTPTargetType, u *User) error {
	return otpService.RequestOTP(ctx, otp.RequestOTPPayload{
		OTPDisburse: []otp.OTPDispatchType{otp.OTP_DISPATCH_EMAIL, otp.OTP_DISPATCH_SMS},
		OTPType:     otp.OTP_SHORT,
		TargetRef:   u.AccountRef,
		Expiration:  time.Minute * 10,
		Target:      targetType,
		DispatchInfo: otp.DispatchInfo{
			Fullname:   string(u.FullName),
			Email:      u.Email,
			Phone:      u.Phone,
			Expiration: "10 Minutes",
		},
	})
}

func publishEvent(ctx context.Context, et userEventType, u *User) error {
	handlers, ok := userEventFactory[et]
	if !ok {
		return fmt.Errorf("User Event Handler not found for event: %s\n", et)
	}

	for _, f := range handlers {
		err := f(ctx, u)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"Event-Type": et,
				"User":       u.toString(),
			}).Errorf("User Event Failed: %v\n", err)
		}
	}
	return nil
}
