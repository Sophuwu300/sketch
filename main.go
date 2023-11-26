package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/term"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

var w, h int

func getTermSize() {
	W, H, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not get terminal size")
		os.Exit(1)
	}
	w, h = W, H
}

func getSize(x, y int) (uint, uint) {
	if y > h {
		x = x * h / y
		y = h
	}
	if x > w {
		y = y * w / x
		x = w
	}
	return uint(x), uint(y) / 2
}

func getImg(path string) (image.Image, error) {
	var img image.Image
	file, err := os.Open(path)
	if err != nil {
		return img, fmt.Errorf("error opening file: %s", path)
	}
	defer file.Close()
	img, _, err = image.Decode(file)
	if err != nil {
		return img, fmt.Errorf("error decoding file: %s", path)
	}
	return img, nil
}

func printImg(path string) error {

	img, err := getImg(path)
	if err != nil {
		return err
	}

	W, H := getSize(img.Bounds().Dx(), img.Bounds().Dy())

	img = resize.Resize(W, H, img, resize.MitchellNetravali)

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			fmt.Printf("\033[48;2;%d;%d;%dm \033[0m", r>>8, g>>8, b>>8)
		}
		fmt.Println()
	}
	fmt.Println()
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "fatal: no image specified")
		os.Exit(1)
	}
	getTermSize()
	var errs []error
	for _, path := range os.Args[1:] {
		err := printImg(path)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}