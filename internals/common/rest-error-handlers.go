package common

import (
	"github.com/gofiber/fiber/v2"
)

func SendSuccess(message string, c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{"success": true, "message": message})
}

func SendSessionError(s *SessionError, c *fiber.Ctx) error {
	return c.Status(int(s.ResponseCode)).JSON(fiber.Map{"success": false, "message": s.Err.Error()})
}

func SendError(err error, c *fiber.Ctx) error {
	if err == nil {
		return nil
	} else if s, ok := err.(*FieldValidationError); ok {
		return c.Status(400).JSON(fiber.Map{"success": false, "error": s.Errs})
	} else if s, ok := err.(*SessionError); ok && s.Err != nil {
		return c.Status(int(s.ResponseCode)).JSON(fiber.Map{"success": false, "message": s.Err.Error()})
	} else {
		return c.Status(400).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
}
