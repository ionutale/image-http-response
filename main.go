package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	// Imports the Google Cloud Storage client package.

	"github.com/discord/lilliput"
	"google.golang.org/api/iterator"
	// Imports the Google Cloud Storage client package.
	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

func listAllFromBucket() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	// projectID := "beta-dodolandia"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := "dodolandia-layouts-originals"

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

	// Sets the name for the new bucket.
	bucketName := "dodolandia-layouts-originals"

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

func main() {
	// listAllFromBucket()
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/photo", handler)
	http.ListenAndServe(":8080", nil)
}

type ImageFormat struct {
	Format string
	Width  int
	Height int
	Quality int
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	imageName := "test.png"
	if r.URL.Query().Get("name") != "" {
		// dodolandia-layouts-originals/original/5fc3816ade78a2000b9364c6
		imageName = "original/" + r.URL.Query().Get("name")
	}

	// fileBytes, err := ioutil.ReadFile(imageName)
	// if err != nil {
	// 	panic(err)
	// }

	fileBytes := getImageFromBucket(imageName)
	
	// read query params into struct
	format := r.URL.Query().Get("fm")
	width := r.URL.Query().Get("w")
	height := r.URL.Query().Get("h")
	quality := r.URL.Query().Get("q")	

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(newImage)
}

func getEncodeOptions (format string, quality int) map[int]int {
	if format == "jpeg" || format == "jpg" {
		return map[int]int{lilliput.JpegQuality: quality}
	} else if format == "png" {
		return map[int]int{lilliput.PngCompression: quality}
	} else if format == "webp" {
		return map[int]int{lilliput.WebpQuality: quality}
	} else {
		return map[int]int{}
	}
}

func ProcessImage(image []byte, format string, width int, height int, quality int) ([]byte, error) {
	decoder, err := lilliput.NewDecoder(image)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		os.Exit(1)
	}
	defer decoder.Close()

	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	fmt.Println(format, width, height, quality)

	opts := &lilliput.ImageOptions{
		FileType:             format,
		Width:                width,
		Height:               height,
		ResizeMethod:         lilliput.ImageOpsFit,
		NormalizeOrientation: true,
		EncodeOptions:        getEncodeOptions(format, quality),
	}

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)
	
	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error transforming image, %s\n", err)
		panic(err)
	}

	return outputImg, err
}

