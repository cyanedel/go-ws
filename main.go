package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	annyeonghaseyo()

	app := fiber.New()

	app.Get("/hello/:lang", func(c *fiber.Ctx) error {
		lang := c.Params("lang")
		msg := "Hello World!"
		if lang == "jp" {
			msg = "こんにちは世界!"
		} else if lang == "kr" {
			msg = "안녕하세요 세상"
		}

		return c.SendString(msg)
	})

	app.Listen(":3000")
}

func annyeonghaseyo() {
	fmt.Println("안녕하세요 세상")
}
