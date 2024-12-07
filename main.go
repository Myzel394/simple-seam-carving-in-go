package main

import (
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"os"
	"myzel394.app/image-stuff/imageutils"
)

func main() {
	// Read image
	reader, _ := os.Open("./assets/water2.png")
	rawImage, _, _ := image.Decode(reader)
	readImage := imageutils.ImageAnalyzer{Image: rawImage}

	bounds := readImage.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	writeImage := image.NewRGBA(bounds)

	// Main action
	for x := range width {
		for y := range height {
			energy := uint8(readImage.CalculateEnergyAt(x, y))

			color := color.RGBA{energy, energy, energy, 255}
			writeImage.Set(x, y, color)
		}
	}

	// Out image
	writer, _ := os.Create("image.png")
	png.Encode(writer, writeImage)
}
