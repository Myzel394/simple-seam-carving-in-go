package main

import (
	"image"
	"image/color"
	"math"
)

type ImageAnalyzer struct {
	image.Image
}

// Calculate the energy using the dual-gradient energy function
func getEnergy(
	firstPixel color.Color,
	secondPixel color.Color,
) float64 {
	firstRed, firstGreen, firstBlue, _ := firstPixel.RGBA()
	firstRed, firstGreen, firstBlue = firstRed >> 8, firstGreen >> 8, firstBlue >> 8
	secondRed, secondGreen, secondBlue, _ := secondPixel.RGBA()
	secondRed, secondGreen, secondBlue = secondRed >> 8, secondGreen >> 8, secondBlue >> 8

	return float64(
		math.Abs(float64(firstRed - secondRed)) +
		math.Abs(float64(firstGreen - secondGreen)) +
		math.Abs(float64(firstBlue - secondBlue)),
	)
}

// Calculate a 0-255 grayscale value from a color using the NTSC formula
func colorToGray(color color.Color) float64 {
	red, green, blue, _ := color.RGBA()
	red, green, blue = red >> 8, green >> 8, blue >> 8

	return 0.299 * float64(red) +
		0.587 * float64(green) +
		0.114 * float64(blue)
}

func sumSlice(slice []float64) float64 {
	var sum float64
	for _, value := range slice {
		sum += value
	}
	return sum
}

func (image *ImageAnalyzer) CalculateEnergyAt(x int, y int) float64 {
	northPixel := image.At(
		x,
		max(0, y - 1),
	)
	northEastPixel := image.At(
		min(image.Bounds().Max.X, x + 1),
		max(0, y - 1),
	)
	eastPixel := image.At(
		min(image.Bounds().Max.X, x + 1),
		y,
	)
	southEastPixel := image.At(
		min(image.Bounds().Max.X, x + 1),
		min(image.Bounds().Max.Y, y + 1),
	)
	southPixel := image.At(
		x,
		min(image.Bounds().Max.Y, y + 1),
	)
	southWestPixel := image.At(
		max(0, x - 1),
		min(image.Bounds().Max.Y, y + 1),
	)
	westPixel := image.At(
		max(0, x - 1),
		y,
	)
	northWestPixel := image.At(
		max(0, x - 1),
		max(0, y - 1),
	)

	thisPixel := image.At(
		x,
		y,
	)

	_ = thisPixel
	_ = northPixel
	_ = northEastPixel
	_ = southPixel

	horizontalMatrix := []float64{
		colorToGray(northWestPixel), 0, -colorToGray(northEastPixel),
		2 * colorToGray(westPixel), 0, -2 * colorToGray(eastPixel),
		colorToGray(southWestPixel), 0, -colorToGray(southEastPixel),
	}
	verticalMatrix := []float64{
		colorToGray(northWestPixel), 2 * colorToGray(northPixel), colorToGray(northEastPixel),
		0, 0, 0,
		-colorToGray(southWestPixel), -2 * colorToGray(southPixel), -colorToGray(southEastPixel),
	}

	g_x := sumSlice(horizontalMatrix)
	g_y := sumSlice(verticalMatrix)

	return math.Sqrt(
		math.Pow(g_x, 2) +
		math.Pow(g_y, 2),
	)
}
