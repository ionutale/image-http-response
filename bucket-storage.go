package main

import (
	"fmt"
	"io/ioutil"
	"log"

	// Imports the Google Cloud Storage client package.
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

var bucketName = "dodolandia-layouts"
var prefix = "original/"

type Image struct {
	Name string
	Format string
	Width  int
	Height int
	Quality int
}

func getImageFromBucket(imageName string) []byte, error {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates a ObjectHandle instance.
	object := bucket.Object(prefix + imageName)
	log.Println("image name", prefix + imageName)

	// Creates a Reader instance.
	reader, err := object.NewReader(ctx)
	if err != nil {
		return nil, err
	}

	// Reads the contents of the object.
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func uploadImageToBucket(imageName string, imageBytes []byte) error {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return err
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates a ObjectHandle instance.
	object := bucket.Object(prefix + imageName)

	// Creates a Writer instance.
	writer := object.NewWriter(ctx)

	// Writes data to the object.
	if _, err := writer.Write(imageBytes); err != nil {
		log.Fatalf("Failed to write data: %v", err)
		return err
	}

	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
		return err
	}

	fmt.Printf("File %v uploaded.\n", imageName)
	return nil
}