package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	fn()

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	os.Stdout = oldStdout

	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}

	return string(out)
}

func withArgs(t *testing.T, args []string, fn func()) {
	t.Helper()

	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	t.Cleanup(func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	})

	os.Args = append([]string{"img2ascii"}, args...)
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flag.CommandLine.SetOutput(io.Discard)

	fn()
}

func writeTestPNG(t *testing.T, c color.Color) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, c)

	path := filepath.Join(t.TempDir(), "pixel.png")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		_ = file.Close()
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestMainDisplaysRamps(t *testing.T) {
	var out string
	withArgs(t, []string{"-display-ramp"}, func() {
		out = captureStdout(t, main)
	})

	for _, name := range rampNames {
		want := fmt.Sprintf("%-10s [%s]\n", name, ramps[name])
		if !strings.Contains(out, want) {
			t.Fatalf("expected ramp listing to contain %q, got %q", want, out)
		}
	}
}

func TestMainConvertsImageWithExplicitSizeAndRamp(t *testing.T) {
	path := writeTestPNG(t, color.White)

	var out string
	withArgs(t, []string{"-width", "1", "-height", "1", "-ramp", "binary", path}, func() {
		out = captureStdout(t, main)
	})

	if out != "#\n" {
		t.Fatalf("expected converted ASCII art, got %q", out)
	}
}
