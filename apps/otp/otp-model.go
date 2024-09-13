package otp

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/thrillee/triq/internals/schemas"
)

type OTP struct {
	schemas.BaseModel

	TargetRef string

	Target     string
	Code       string
	Expiration time.Time
}

func (u OTP) GetID() interface{} {
	return u.ID
}

func (u OTP) GetModelName() interface{} {
	return "OTP"
}

func (u OTP) QueryRepo() schemas.Repository[OTP] {
	return schemas.NewBaseRepository[OTP](&OTP{})
}

func (u *OTP) FindByEventUserId(ctx context.Context, targetRef OTPTargetType, target string) *OTP {
	var otp OTP
	u.QueryRepo().Query(ctx).Model(OTP{TargetRef: string(targetRef), Target: target}).First(&otp)
	return &otp
}

func (u *OTP) FindByTargetRef(ctx context.Context, target string) *OTP {
	var otp OTP
	u.QueryRepo().Query(ctx).Model(OTP{Target: target}).First(&otp)
	return &otp
}

func (u *OTP) Save(ctx context.Context) {
	if u.ID == uuid.Nil {
		u.QueryRepo().Save(ctx, u)
	} else {
		u.QueryRepo().Edit(ctx, u.ID, u)
	}
}
