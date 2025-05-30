package main

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

func LoadImage(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, err
}

func RenderImage(filepath string, row int, col int, width_ int, height_ int) {
	img, err := LoadImage(filepath)
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
		return
	}

	rgba := convertToRGBA(img)
	rgba, _ = Resize(*rgba, width_, height_)
	encoded := encodeImageToBase64RGBA(rgba)

	width := rgba.Rect.Dx()
	height := rgba.Rect.Dy()

	fmt.Printf("\033[s\033[%d;%dH", row, col)

	chunk_size := 4096
	pos := 0
	first := true
	for pos < len(encoded) {
		fmt.Print("\033_G")
		if first {
			fmt.Printf("q=2,a=T,f=32,s=%d,v=%d,", width, height)
			first = false
		}
		chunk_len := len(encoded) - pos
		if chunk_len > chunk_size {
			chunk_len = chunk_size
		}
		if pos+chunk_len < len(encoded) {
			fmt.Print("m=1")
		}
		fmt.Printf(";%s\033\\", encoded[pos:pos+chunk_len])
		pos += chunk_len
	}
	fmt.Print("\033[u")
}

func Resize(img image.RGBA, width int, height int) (*image.RGBA, error) {
	if width <= 0 || height <= 0 {
		return nil, nil
	}

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))
	scaleX := float64(img.Bounds().Dx()) / float64(width)
	scaleY := float64(img.Bounds().Dy()) / float64(height)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			srcX := int(float64(x) * scaleX)
			srcY := int(float64(y) * scaleY)

			r, g, b, a := getPixel(img, srcX, srcY)
			newImg.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return newImg, nil
}

func encodeImageToBase64RGBA(rgba *image.RGBA) string {
	return base64.StdEncoding.EncodeToString(rgba.Pix)
}

func getPixel(img image.RGBA, x int, y int) (uint8, uint8, uint8, uint8) {
	index := img.PixOffset(x, y)
	return img.Pix[index], img.Pix[index+1], img.Pix[index+2], img.Pix[index+3]
}

func convertToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba
}
