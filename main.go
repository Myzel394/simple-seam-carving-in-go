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
	reader, _ := os.Open("./assets/surfer.png")
	rawImage, _, _ := image.Decode(reader)
	readImage := imageutils.ImageAnalyzer{Image: rawImage}

	bounds := readImage.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y
	writeImage := image.NewRGBA(bounds)

	seams := imageutils.NewImageSeams()
	seams.CreateSeamsFromRectangle(bounds)

	// Main action
	for x := range width {
		for y := range height {
			energy := uint(readImage.CalculateEnergyAt(x, y))
			hexColor := uint8(energy)

			color := color.RGBA{hexColor, hexColor, hexColor, 255}
			writeImage.Set(x, y, color)

			seams.SetCostForNode(x, y, energy)
		}
	}

	seams.CreateOptimizedRoutes()
	lowestSeam := seams.GetLowestSeam()
	lowestSeam.WriteSeamChainToImage(writeImage)

	// Out image
	writer, _ := os.Create("image.png")
	png.Encode(writer, writeImage)
}
