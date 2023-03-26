package main

import (
	"bytes"
	"io"
	"log"
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024, // 100MB
	})

	app.Get("/", healthCheck)
	app.Get("/health", healthCheck)
	app.Get("/photo/:name", getImage_Handler)
	app.Post("/images", uploadImage_Handler)

	log.Fatal(app.Listen(":8080"))
	PrintMemUsage("Server started")
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"alive": true, "v": "5"})
}

func getImage_Handler(c *fiber.Ctx) error {
	PrintMemUsage("Requested image")
	imageName := c.Params("name")

	fileBytes, err := getImageFromBucket(imageName)
	if err != nil {
		log.Println("Failed to get image from bucket: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
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
		log.Println("Failed to process image: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// w.WriteHeader(http.StatusOK)
	// w.Header().Set("Content-Type", "application/octet-stream")
	// w.Write(newImage)
	// c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Type", "image/png")
	c.Send(newImage)

	runtime.GC()
	PrintMemUsage("After gc")
	return nil
}

func uploadImage_Handler(c *fiber.Ctx) error {
	log.Println("File size: ", c.Request().Header.ContentLength())
	// read the file
	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Failed to read file: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// read the file into a byte array
	fileBytes, err := file.Open()
	if err != nil {
		log.Println("Failed to open file: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	defer fileBytes.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, fileBytes); err != nil {
		log.Println("Failed to copy file: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	// upload the file to the bucket
	err = uploadImageToBucket(file.Filename, buf.Bytes())
	if err != nil {
		log.Println("Failed to upload image to bucket: ", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"success": true})
}
