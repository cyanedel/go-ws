package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

type FileDetails struct {
	Message  string `json:"message"`
	Status   int    `json:"status"`
	FileName string `json:"filename"`
	FileMime string `json:"mime"`
	FileExt  string `json:"ext"`
	FileSize int    `json:"size"`
}

func main() {
	annyeonghaseyo()

	app := fiber.New(fiber.Config{StreamRequestBody: true})
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

	app.Get("/download", func(c *fiber.Ctx) error {
		targetURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/2/24/Samsung_Logo.svg/2560px-Samsung_Logo.svg.png"

		resp, err := http.Get(targetURL)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		defer resp.Body.Close()

		// Change filename and extension here
		newFilename := "Samsung_Logo.jpg"

		// Set headers for download
		c.Set("Content-Disposition", "attachment; filename=\""+newFilename+"\"")
		c.Set("Content-Type", "image/jpg") // Adjust content-type as needed

		// Copy response body to the client
		_, err = io.Copy(c.Response().BodyWriter(), resp.Body)
		if err != nil {
			return err
		}

		return c.SendStatus(http.StatusOK)
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Error reading uploaded file")
		}

		fileHeader, _ := file.Open()
		defer fileHeader.Close()

		_, err = fileHeader.Seek(0, io.SeekStart)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error seeking in uploaded file")
		}

		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		ext = strings.TrimPrefix(ext, ".")

		m, _ := mimetype.DetectReader(fileHeader)

		wsResponse := FileDetails{
			Message:  "Complete",
			Status:   200,
			FileName: file.Filename,
			FileExt:  ext,
			FileMime: m.String(),
			FileSize: int(file.Size),
		}

		// response, _ := json.Marshal(wsResponse)

		fmt.Println(wsResponse)

		return c.JSON(wsResponse)
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
		file, _ := c.FormFile("file")

		log.Println("file", file)

		client := &http.Client{}
		req, _ := http.NewRequest("POST", "http://127.0.0.1:3000/upload", nil)
		fileData, _ := file.Open()
		defer fileData.Close()

		// Attach file data to the request
		req.Body = io.NopCloser(fileData)

		// Send the request to the internal server
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		// Handle response from internal server
		body, _ := io.ReadAll(resp.Body)

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
