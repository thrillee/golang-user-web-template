package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/thrillee/triq/internals/common"
)

type UserREST struct {
	service UserService
}

func NewUserREST(app *fiber.App) {
	rest := UserREST{
		service: NewUserService(),
	}

	authURLs := app.Group("/api/v1/auth")
	authURLs.Post("login", rest.Login)
	authURLs.Post("resend-verify-otp", rest.ResendVerifyAccountOTP)
	authURLs.Post("forgot-password", rest.ForgotPassword)
	authURLs.Post("reset-password", rest.ResetPassword)
	authURLs.Post("register", rest.RegisterUser)
	authURLs.Post("verify", rest.VerifyAccount)

	userURL := app.Group("/api/v1/accounts")
	userURL.Use(IsAuthenticated)
	userURL.Put("change-dp", rest.UploadDisplayPicutre)
	userURL.Post("change-password", rest.ChangePassword)
	userURL.Get("auth-me", rest.GetUserInSession)
	userURL.Put("edit", rest.EditUser)
}

func (r UserREST) GetUserInSession(c *fiber.Ctx) error {
	user, err := GetAuthUser(c)
	if err != nil {
		return common.SendError(err, c)
	}
	return c.JSON(user)
}

func (r UserREST) ResendVerifyAccountOTP(c *fiber.Ctx) error {
	var data OTPResendPayload
	err := c.BodyParser(&data)
	if err != nil {
		return common.SendError(fmt.Errorf("Resend OTP Failed: %v", err), c)
	}

	err = handleVerifyAccountResend(c.Context(), &data)
	if err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success": true,
		"message": "A new OTP have been generated for you.",
	}

	return c.JSON(res)
}

func (r UserREST) GetAuthUser(c *fiber.Ctx) error {
	authUser, err := GetAuthUser(c)
	if err != nil {
		return common.SendError(fmt.Errorf("Invalid user session: %v", err), c)
	}

	return c.JSON(authUser)
}

func (r UserREST) UploadDisplayPicutre(c *fiber.Ctx) error {
	authUser, err := GetAuthUser(c)
	if err != nil {
		return common.SendError(fmt.Errorf("Upload Display Picture Failed: %v", err), c)
	}

	file, err := c.FormFile("dp")
	if err != nil {
		return common.SendError(fmt.Errorf("Upload Display Picture Failed: %v", err), c)
	}

	err = r.service.UploadDisplayPicutre(c.Context(), authUser, file)
	if err != nil {
		return common.SendError(fmt.Errorf("Upload Display Picture Failed: %v", err), c)
	}

	return c.JSON(authUser)
}

func (r UserREST) ChangePassword(c *fiber.Ctx) error {
	user, err := GetAuthUser(c)
	if err != nil {
		return common.SendError(fmt.Errorf("Change password Failed: %v", err), c)
	}

	var data ChangePasswordPayload
	err = c.BodyParser(&data)
	if err != nil {
		return common.SendError(fmt.Errorf("Change password Failed: %v", err), c)
	}

	if err = r.service.ChangePassword(c.Context(), user, &data); err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success": true,
		"message": "You have successfully changed your password",
	}

	return c.JSON(res)
}

func (r UserREST) Login(c *fiber.Ctx) error {
	var data LoginPayload
	err := c.BodyParser(&data)
	if err != nil {
		return common.SendError(err, c)
	}

	loginResponse, err := r.service.Login(c.Context(), &data)
	if err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success":    true,
		"user":       loginResponse.AuthUser,
		"mfa":        loginResponse.MfaRequired,
		"nextAction": loginResponse.NextAction,
	}

	cookie := new(fiber.Cookie)
	cookie.Name = loginResponse.Cookie.Name
	cookie.Path = loginResponse.Cookie.Path
	cookie.Domain = loginResponse.Cookie.Domain
	cookie.Value = loginResponse.Cookie.Value
	cookie.Expires = loginResponse.Cookie.Expires
	cookie.Secure = loginResponse.Cookie.Secure
	cookie.SameSite = loginResponse.Cookie.SameSite
	cookie.HTTPOnly = loginResponse.Cookie.HTTPOnly
	cookie.MaxAge = loginResponse.Cookie.MaxAge

	c.Cookie(cookie)

	return c.JSON(res)
}

func (r UserREST) ResetPassword(c *fiber.Ctx) error {
	var data VerifyAccountPayload
	err := c.BodyParser(&data)
	if err != nil {
		return common.SendError(err, c)
	}

	if err = r.service.ResetPassword(c.Context(), &data); err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success": true,
		"message": "You have successfully changed your password",
	}

	return c.JSON(res)
}

func (r UserREST) ForgotPassword(c *fiber.Ctx) error {
	var data VerifyAccountPayload
	err := c.BodyParser(&data)
	if err != nil {
		return common.SendError(err, c)
	}

	if err = r.service.ForgotPassword(c.Context(), &data); err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success": true,
		"message": "An OTP have been sent to you, use the OTP to complete your password reset.",
	}

	return c.JSON(res)
}

func (r UserREST) EditUser(c *fiber.Ctx) error {
	accountRef := c.Params("accountRef")

	var data NewUserPayload
	err := c.BodyParser(&data)
	if err != nil {
		return common.SendError(err, c)
	}

	user, err := r.service.EditUser(c.Context(), accountRef, &data)
	if err != nil {
		return common.SendError(err, c)
	}

	return c.JSON(user)
}

func (r UserREST) VerifyAccount(c *fiber.Ctx) error {
	var data VerifyAccountPayload
	err := c.BodyParser(&data)
	if err != nil {
		return common.SendError(err, c)
	}

	user, err := r.service.VerifyAccount(c.Context(), &data)
	if err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success":    true,
		"message":    "Account Verified Successfully",
		"userEmail":  user.Email,
		"accountRef": user.AccountRef,
	}

	return c.JSON(res)
}

func (r UserREST) RegisterUser(c *fiber.Ctx) error {
	var newUser NewUserPayload
	err := c.BodyParser(&newUser)
	if err != nil {
		return common.SendError(err, c)
	}

	user, err := r.service.CreateUser(c.Context(), &newUser)
	if err != nil {
		return common.SendError(err, c)
	}

	res := map[string]interface{}{
		"success":    true,
		"message":    "Registration successful",
		"userEmail":  user.Email,
		"accountRef": user.AccountRef,
	}

	return c.JSON(res)
}
