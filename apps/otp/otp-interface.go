package otp

import (
	"context"
	"time"
)

type (
	OTPCodeType     string
	OTPDispatchType string
	OTPTargetType   string
)

const (
	OTP_LONG  OTPCodeType = "long"
	OTP_SHORT OTPCodeType = "short"
)

const (
	OTP_DISPATCH_SMS   OTPDispatchType = "sms"
	OTP_DISPATCH_EMAIL OTPDispatchType = "email"
	OTP_DISPATCH_ALL   OTPDispatchType = "all"
)

const (
	OTP_VERIFY_ACCOUNT OTPTargetType = "verify_account"
	OTP_RESET_PASSWORD OTPTargetType = "reset_password"
)

type DispatchInfo struct {
	Fullname   string
	Email      string
	Phone      string
	Expiration string
}

type RequestOTPPayload struct {
	TargetRef    string
	Target       OTPTargetType
	Expiration   time.Duration
	OTPType      OTPCodeType
	OTPDisburse  []OTPDispatchType
	DispatchInfo DispatchInfo
}

type VerifyOTPPayload struct {
	Target    OTPTargetType
	TargetRef string
	OTPCode   string
}

type OTPResendPayload struct {
	TargetRef string        `json:"target_ref"`
	Target    OTPTargetType `json:"target"`
}

type OTPService interface {
	RequestOTP(context.Context, RequestOTPPayload) error
	VerifyOTP(context.Context, VerifyOTPPayload) (bool, error)
	ResendOTP(context.Context, OTPResendPayload) error
}
