package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func (hndler *SqsHandler) processImage(ctx context.Context, body []byte) ([]string, error) {
	// Decode the PNG image
	img, _, err := image.Decode(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	// Define dimensions for the images
	dimensions := []struct {
		Width    int
		Height   int
		Filename string
	}{
		{720, 480, "output_720x480.jpg"},
		{1200, 800, "output_1200x800.jpg"},
	}

	// Set the target file size in kilobytes
	targetFileSizeKB := 500

	reqId := uuid.New().String()
	var imageUrls []string
	for _, dim := range dimensions {
		resizedImage := imaging.Resize(img, dim.Width, dim.Height, imaging.Lanczos)

		var resizedImg image.Image = resizedImage

		// Loop to adjust the JPEG quality until the file size meets the target
		quality := 90
		for {
			// Create a buffer to store the JPEG image
			var buf bytes.Buffer
			// Encode the image to JPEG format with the current quality
			err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: quality})
			if err != nil {
				log.Fatal(err)
			}

			// Check the file size
			fileSizeKB := float64(buf.Len()) / 1024

			// If the file size is within the target, write the image to a file and break the loop
			if fileSizeKB <= float64(targetFileSizeKB) {
				id := uuid.New().String()
				fileName := fmt.Sprintf("%s-%d-%d.jpg", id, dim.Width, dim.Height)

				var fileUrl string
				fileUrl, err = hndler.Bucket.Put(ctx, reqId, fileName, "images/", buf.Bytes())
				if err != nil {
					return nil, err
				}
				imageUrls = append(imageUrls, fileUrl)
			}

			// Decrease the quality for the next iteration
			quality -= 10

			// Break the loop if the quality becomes too low
			if quality <= 0 {
				log.Fatal("Unable to compress to the desired size.")
			}
		}
	}

	return imageUrls, nil
}
