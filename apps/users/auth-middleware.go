package users

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/thrillee/triq/internals/common"
)

func IsAuthenticated(c *fiber.Ctx) error {
	authToken := c.Cookies(AUTH_COOKIE_NAME)
	claims, err := getTokenClaims(authToken)
	if err != nil {
		return common.SendError(common.NewSessionError(
			common.UNAUTHORIZED, fmt.Errorf("Session Expired. Kindly Login again")), c)
	}

	accountRef := claims["accountRef"].(string)
	query := User{}
	user := query.FindByAccountRef(c.Context(), accountRef)
	if user.IsNew() {
		return common.SendError(common.NewSessionError(
			common.UNAUTHORIZED, fmt.Errorf("Invalid User Session")), c)
	}

	c.Locals("user", user)

	return c.Next()
}

func GetAuthUser(c *fiber.Ctx) (*User, error) {
	userGenenric := c.Locals("user")
	user := userGenenric.(*User)
	if user.IsNew() {
		return nil, common.NewSessionError(common.BAD_REQUEST, fmt.Errorf("Invalid user session"))
	}

	return user, nil
}
