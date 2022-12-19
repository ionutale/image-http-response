package main

import (
	"fmt"
	"net/http"
	"strconv"
)


func main() {
	// listAllFromBucket()
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/", http.HandlerFunc(healthCheck))
	fmt.Println("server is up and running")
	http.Handle("/photo", handler)
	http.ListenAndServe(":8080", nil)
}

type ImageFormat struct {
	Format string
	Width  int
	Height int
	Quality int
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"alive": true, "v": "3"}`))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	imageName := "test.png"
	if r.URL.Query().Get("name") != "" {
		// dodolandia-layouts-originals/original/5fc3816ade78a2000b9364c6
		imageName = "original/" + r.URL.Query().Get("name")
	}

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
