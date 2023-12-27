package routes // model : go-be

import (
	"go-be/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func Routes(app *fiber.App) {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	requestAuthen := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	})

	auth := app.Group("/auth")
	auth.Post("/regis", controllers.Regis)
	auth.Post("/login", controllers.Login)
	auth.Get("/readed", requestAuthen, controllers.LoginReaded)

}
