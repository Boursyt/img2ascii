package main

import (
	"flag"
	"fmt"
	"image"
	"os"

	"img2ascii/ascii"

	"golang.org/x/term"
)

var ramps = map[string]string{
	"short":     ascii.RampShort,
	"short-alt": ascii.RampShortAlt,
	"bourke70":  ascii.RampBourke70,
	"blocks":    ascii.RampBlocks,
	"inverted":  ascii.RampInverted,
	"binary":    ascii.RampBinary,
}

var rampNames = []string{"short", "short-alt", "bourke70", "blocks", "inverted", "binary"}

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

func displayRamp() {
	for _, name := range rampNames {
		fmt.Printf("%-10s [%s]\n", name, ramps[name])
	}
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: img2ascii [flags] <image-path>")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	width := flag.Int("width", 0, "output width in terminal cells (0 uses terminal width)")
	height := flag.Int("height", 0, "output height in terminal cells (0 uses terminal height minus one row)")
	rampName := flag.String("ramp", "short", "character ramp: short, short-alt, bourke70, blocks, inverted, binary")
	showRamps := flag.Bool("display-ramp", false, "display available character ramps")
	flag.Parse()

	if *showRamps {
		displayRamp()
		return
	}

	if flag.NArg() != 1 {
		fail("usage: img2ascii [flags] <image-path>")
	}

	ramp, ok := ramps[*rampName]
	if !ok {
		fail("unknown ramp %q", *rampName)
	}

	termW, termH := terminalSize()
	if *width == 0 {
		*width = termW
	}
	if *height == 0 {
		*height = termH
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fail("open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		fail("decode image: %v", err)
	}

	asciiArt, err := ascii.ConvertToAsciiWithRamp(img, *width, *height, ramp)
	if err != nil {
		fail("convert image: %v", err)
	}
	fmt.Print(asciiArt)
}
