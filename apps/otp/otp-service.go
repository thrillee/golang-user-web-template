package otp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thrillee/triq/internals/schemas"
	"github.com/thrillee/triq/internals/security"
)

func NewOTPService() *OTPImpl {
	otp := OTP{}
	return &OTPImpl{
		repo:  otp.QueryRepo(),
		model: otp,
	}
}

type OTPImpl struct {
	repo  schemas.Repository[OTP]
	model OTP
}

func (i *OTPImpl) VerifyOTP(ctx context.Context, d VerifyOTPPayload) (bool, error) {
	targetOTP := i.model.FindByEventUserId(ctx, d.Target, d.TargetRef)
	if targetOTP == nil || !security.CheckPasswordHash(d.OTPCode, targetOTP.Code) {
		return false, fmt.Errorf("OTP is not valid")
	}

	if time.Now().After(targetOTP.Expiration) {
		return false, fmt.Errorf("OTP Expired")
	}

	targetOTP.Expiration = time.Now()
	targetOTP.Save(ctx)

	return true, nil
}

func (i *OTPImpl) RequestOTP(ctx context.Context, req RequestOTPPayload) error {
	r := OTP{}

	otpCode := codeGenerator(req.OTPType)
	log.Println("OTP: ", otpCode)

	codeHash, err := security.HashPassword(otpCode)
	if err != nil {
		return fmt.Errorf("Generate OTP Failed: %v", err)
	}

	targetOTP := r.FindByEventUserId(ctx, req.Target, req.TargetRef)
	if targetOTP == nil {
		targetOTP = &OTP{
			TargetRef: req.TargetRef,
			Target:    string(req.Target),
		}
	}

	targetOTP.Target = string(req.Target)
	targetOTP.TargetRef = req.TargetRef
	targetOTP.Code = codeHash
	targetOTP.Expiration = time.Now().Add(req.Expiration)
	targetOTP.Save(ctx)

	go handleDispatch(&dispatchPayload{
		expiration:   fmt.Sprintf("%s minutes", req.Expiration),
		dispatchType: req.OTPDisburse,
		dispatchInfo: req.DispatchInfo,
		codeType:     req.OTPType,
		code:         otpCode,
	})

	return nil
}

func codeGenerator(ct OTPCodeType) string {
	switch ct {
	case OTP_LONG:
		return uuid.NewString()
	default:
		return generateOTP(6)
	}
}
