package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginReaded(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return c.Status(200).JSON(fiber.Map{
		"email":  claims["email"],
		"role":   claims["role"],
		"status": claims["status"],
	})
}
