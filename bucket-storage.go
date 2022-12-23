package main

import (
	"fmt"
	"log"
	"io/ioutil"
		// Imports the Google Cloud Storage client package.
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"

)

var bucketName = "dodolandia-layouts-originals-public"
var prefix = "new-originals/"

type Image struct {
	Name string
	Format string
	Width  int
	Height int
	Quality int
}

func listAllFromBucket() {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	query := &storage.Query{Prefix: prefix}
	it := bucket.Objects(ctx, query)

	count := 0
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Errorf("listBucket: unable to list bucket %q: %v", bucket, err)
			return
		}
		count++
		fmt.Println("cur object:", obj)
	}
	fmt.Println("total-count", count)
}

func getImageFromBucket(imageName string) []byte {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bucketName)

	// Creates a ObjectHandle instance.
	object := bucket.Object(prefix + imageName)

	// Creates a Reader instance.
	reader, err := object.NewReader(ctx)
	if err != nil {
		log.Println("Failed to create reader: ", imageName, bucketName, err)
		
	}

	// Reads the contents of the object.
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}

	return data
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