package ascii

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/nfnt/resize"
)

const (
	RampShort    = " .:-=+*#%@"
	RampShortAlt = " .,:;i1tfLCG08@"
	RampBourke70 = " .'`^\",:;Il!i><~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	RampBlocks   = " ░▒▓█"
	RampInverted = "@%#*+=-:. "
	RampBinary   = " #"
)

// charAspect compensates for terminal cells being taller than wide (~2:1).
const charAspect = 2.0

// imgResize scales the image to fit within maxW x maxH cells while preserving
func imgResize(img image.Image, maxW, maxH int) image.Image {
	bounds := img.Bounds()
	imgW := float64(bounds.Dx())
	imgH := float64(bounds.Dy())

	scaleW := float64(maxW) / (imgW * charAspect)
	scaleH := float64(maxH) / imgH
	scale := scaleW
	if scaleH < scale {
		scale = scaleH
	}

	newW := uint(imgW * charAspect * scale)
	newH := uint(imgH * scale)
	if newW == 0 {
		newW = 1
	}
	if newH == 0 {
		newH = 1
	}

	return resize.Resize(newW, newH, img, resize.Lanczos3)
}

// rgbaAt: convert the Red Green Blue Alpha values of the image into a single uint32
func rgbaAt(img image.Image, x, y int) (r, g, b, a uint32) {
	return img.At(x, y).RGBA()
}

// luminance: calculate the luminance of the image
func luminance(r, g, b uint32) float32 {
	return 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)
}

func ConvertToAscii(img image.Image, maxW, maxH int) string {
	img = imgResize(img, maxW, maxH)
	bounds := img.Bounds()
	var sb strings.Builder

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := rgbaAt(img, x, y)
			lum := luminance(r, g, b)
			shade := int(lum * float32(len(RampShortAlt)-1) / 0xFFFF)
			sb.WriteByte(RampShortAlt[shade])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
