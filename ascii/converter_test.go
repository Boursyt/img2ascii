package ascii

import (
	"errors"
	"image"
	"image/color"
	"strings"
	"testing"
	"unicode/utf8"
)

func solidImage(w, h int, c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

func TestConvertToAsciiRejectsInvalidSize(t *testing.T) {
	img := solidImage(2, 2, color.White)

	_, err := ConvertToAscii(img, 0, 10)
	if !errors.Is(err, ErrInvalidSize) {
		t.Fatalf("expected ErrInvalidSize, got %v", err)
	}
}

func TestConvertToAsciiWithRampRejectsInvalidRamp(t *testing.T) {
	img := solidImage(2, 2, color.White)

	_, err := ConvertToAsciiWithRamp(img, 10, 10, "#")
	if !errors.Is(err, ErrInvalidRamp) {
		t.Fatalf("expected ErrInvalidRamp, got %v", err)
	}
}

func TestConvertToAsciiWithRampMapsBlackAndWhite(t *testing.T) {
	black := solidImage(1, 1, color.Black)
	blackASCII, err := ConvertToAsciiWithRamp(black, 1, 1, " #")
	if err != nil {
		t.Fatal(err)
	}
	if blackASCII != " \n" {
		t.Fatalf("expected black pixel to use first ramp character, got %q", blackASCII)
	}

	white := solidImage(1, 1, color.White)
	whiteASCII, err := ConvertToAsciiWithRamp(white, 1, 1, " #")
	if err != nil {
		t.Fatal(err)
	}
	if whiteASCII != "#\n" {
		t.Fatalf("expected white pixel to use last ramp character, got %q", whiteASCII)
	}
}

func TestConvertToAsciiFitsWithinRequestedSize(t *testing.T) {
	img := solidImage(100, 50, color.White)

	out, err := ConvertToAscii(img, 20, 10)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSuffix(out, "\n"), "\n")
	if len(lines) > 10 {
		t.Fatalf("expected at most 10 lines, got %d", len(lines))
	}
	for _, line := range lines {
		if utf8.RuneCountInString(line) > 20 {
			t.Fatalf("expected line width at most 20 cells, got %d in %q", utf8.RuneCountInString(line), line)
		}
	}
}

func TestConvertToAsciiWithUnicodeRamp(t *testing.T) {
	img := solidImage(1, 1, color.White)

	out, err := ConvertToAsciiWithRamp(img, 1, 1, RampBlocks)
	if err != nil {
		t.Fatal(err)
	}
	if !utf8.ValidString(out) {
		t.Fatalf("expected valid UTF-8 output, got %q", out)
	}
	if out != "█\n" {
		t.Fatalf("expected white pixel to use last unicode ramp character, got %q", out)
	}
}
