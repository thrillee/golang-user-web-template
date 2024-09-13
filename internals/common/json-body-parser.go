package common

import "github.com/gofiber/fiber/v2"

func JSONParseValidator(c *fiber.Ctx, payload interface{}) error {
	err := c.BodyParser(payload)
	if err != nil {
		return err
	}

	err = Validate(payload)
	if err != nil {
		return err
	}
	return nil
}
