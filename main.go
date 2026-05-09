package main

import (
	"fmt"
	"image"
	"os"

	"img2ascii/ascii"

	"golang.org/x/term"
)

// terminalSize returns the usable terminal width and height in cells.
// Falls back to 80x24 when stdout is not a terminal.
func terminalSize() (int, int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 || h <= 0 {
		return 80, 24
	}
	// keep one row free so the prompt does not push the top off-screen
	return w, h - 1
}

func main() {
	file, err := os.Open("./ascii/test_image/mlp2.jpeg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	w, h := terminalSize()
	asciiArt := ascii.ConvertToAscii(img, w, h)
	fmt.Print(asciiArt)
}
