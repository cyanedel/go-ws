package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	annyeonghaseyo()

	app := fiber.New()
	app2 := fiber.New()

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

	app.Post("/team", func(c *fiber.Ctx) error {

		fmt.Println("team body:", c.Body())
		return c.SendString("200")
	})

	app.Post("/teamforwarder", func(c *fiber.Ctx) error {
		proxy.Do(c, "127.0.0.1:3000/team")
		fmt.Println("team forwarder body:", c.Body())
		return c.SendString("200")
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		fmt.Println("upload 3000 body:", c.Body())
		return c.SendString("200")
	})

	app2.Get("/halo/:lang", func(c *fiber.Ctx) error {
		lang := c.Params("lang")
		msg := "Hello World!"
		if lang == "jp" {
			msg = "こんにちは世界!"
		} else if lang == "kr" {
			msg = "안녕하세요 세상"
		}

		return c.SendString(msg)
	})

	app2.Get("/annyeong", func(c *fiber.Ctx) error {
		return proxy.Do(c, "127.0.0.1:3000/hello/kr")
	})

	app2.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		fmt.Println("file", file)

		internalServerUrl := "http://127.0.0.1:3000/upload" // Replace with your internal server URL

		client := &http.Client{}
		req, err := http.NewRequest("POST", internalServerUrl, nil) // We'll send file data later
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		// Prepare file data for sending
		fileData, err := file.Open()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		defer fileData.Close()

		// Attach file data to the request
		req.Body = io.NopCloser(fileData)

		// Send the request to the internal server
		resp, err := client.Do(req)
		if err != nil {
			return c.SendStatus(http.StatusBadGateway)
		}
		defer resp.Body.Close()

		// Handle response from internal server
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.SendString(string(body))
	})

	go func() {
		app.Listen(":3000")
	}()

	app2.Listen(":3001")
}

func annyeonghaseyo() {
	fmt.Println("안녕하세요 세상")
}
