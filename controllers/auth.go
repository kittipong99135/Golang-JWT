package controllers //model : go-be

import (
	"go-be/database"
	"go-be/models"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Regis(c *fiber.Ctx) error {
	db := database.DBConn
	var regisBody models.User

	err := c.BodyParser(&regisBody)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Register body has invalid.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	var userExists models.User
	result := db.Find(&userExists, "email = ?", strings.TrimSpace(regisBody.Email))
	if result.RowsAffected != 0 {
		return c.Status(503).JSON(fiber.Map{
			"massage": "Error : Email exists.",
			"status":  "error",
		})
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(regisBody.Password), 10)
	if err != nil {
		c.Status(503).JSON(fiber.Map{
			"message": "Error : Password has invalid.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	userRegisted := models.User{
		Email:    regisBody.Email,
		Password: "secretpass:" + string(passHash),
		Status:   "active",
		Role:     "admin",
	}

	db.Create(&userRegisted)
	return c.Status(503).JSON(fiber.Map{
		"massage": "Seccess : Register success.",
		"status":  "seccess",
		"newuser": userRegisted,
	})
}

func Login(c *fiber.Ctx) error {
	db := database.DBConn

	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginBody loginRequest
	err := c.BodyParser(&loginBody)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Logging body has some in invalid.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	var userSelect models.User
	selectResult := db.Find(&userSelect, "email = ?", strings.TrimSpace(loginBody.Email))
	if selectResult.RowsAffected == 0 {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Email has invalid.",
			"status":  "error",
		})
	}

	dwarfsPass := strings.Split(userSelect.Password, ":")[1:][0]
	err = bcrypt.CompareHashAndPassword([]byte(dwarfsPass), []byte(loginBody.Password))
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"message": "Error : Compare password error.",
			"status":  "error",
			"error":   err.Error(),
		})
	}

	claims := jwt.MapClaims{
		"id":     userSelect.ID,
		"email":  userSelect.Email,
		"status": userSelect.Status,
		"role":   userSelect.Role,
		"exp":    time.Now().Add(time.Minute * 5).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return c.Status(200).JSON(fiber.Map{
		"message": "Success : Login success.",
		"status":  "success",
		"token":   userSelect.Email + ":" + t,
	})
}
