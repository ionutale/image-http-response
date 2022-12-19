package main

import (
	"fmt"
	"log"
	"github.com/gofiber/fiber/v2"
	"strconv"
)


func main() {
	app := fiber.New()

	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)
	app.Get("/photo", getImage_Handler)
	app.Get("/photo/:name", handleWithParams)

	log.Fatal(app.Listen(":8080"))
}

func handleWithParams(c *fiber.Ctx) error {
	fmt.Println("handleWithParams", c.Params("name"))
	return c.SendString("handleWithParams: " + c.Params("name"))
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"alive": true, "v": "3"})
}

func getImage_Handler(c *fiber.Ctx) error {

	imageName := ""
	if c.Params("name") != "" {
		imageName = "original/" + c.Params("name")
	}

	fileBytes := getImageFromBucket(imageName)
	
	// read query params into struct
	format := c.Query("fm")
	width := c.Query("w")
	height := c.Query("h")
	quality := c.Query("q")

	// convert query params to int
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		panic(err)
	}

	heightInt, err := strconv.Atoi(height)
	if err != nil {
		panic(err)
	}

	qualityInt, err := strconv.Atoi(quality)
	if err != nil {
		panic(err)
	}

	// process image
	newImage, err := ProcessImage(fileBytes, format, widthInt, heightInt, qualityInt)
	if err != nil {
		panic(err)
	}

	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/octet-stream")
	// w.Write(newImage)
	c.Set("Content-Type", "application/octet-stream")
	return c.Send(newImage)
}
