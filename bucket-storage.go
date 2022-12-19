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

	query := &storage.Query{Prefix: "ori"}
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
	object := bucket.Object(imageName)

	// Creates a Reader instance.
	reader, err := object.NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	// Reads the contents of the object.
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to read data: %v", err)
	}

	return data
}