package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"golang.org/x/term"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"strings"
)

func getTermSize() (int, int) {
	W, H, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || !(W > 0 && H > 0) {
		fmt.Fprintln(os.Stderr, "fatal: could not get terminal size")
		os.Exit(1)
	}
	return W, H
}

func (img *Immg) FitSize(W, H int) {
	y := float64(img.Img.Bounds().Dy())
	x := float64(img.Img.Bounds().Dx())
	w := float64(W)
	h := float64(H)*2 - 1
	if x > w {
		y = y * w / x
		x = w
	}
	if y > h {
		x = x * h / y
		y = h
	}
	img.Img = resize.Resize(uint(math.Round(x)), uint(math.Round(y)), img.Img, resize.MitchellNetravali)
}

func (img *Immg) OpenImg(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %s", path)
	}
	defer file.Close()
	img.Img, _, err = image.Decode(file)
	if err != nil {
		return fmt.Errorf("error decoding file: %s", path)
	}
	return nil
}

type Immg struct {
	Img image.Image
	X   int
	Y   int
}

func (img *Immg) ForX() {
	img.X = 0
	for img.X < 2*(img.Img.Bounds().Dx()) {
		img.Doalp()
		img.X++
		img.Doalp()
		img.X++
		fmt.Print("▀\033[0m")
	}
}

func (img *Immg) Print() {
	img.Y = img.Img.Bounds().Min.Y
	for img.Y < img.Img.Bounds().Max.Y {
		img.ForX()
		fmt.Println()
		img.Y += 2
	}
}

func (img *Immg) Doalp() {
	r, g, b, _ := img.Img.At(img.Img.Bounds().Min.X+img.X/2, img.Y+img.X%2).RGBA()
	fmt.Printf("\033[%d8;2;%d;%d;%dm", 3+(img.X%2), r>>8, g>>8, b>>8)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "fatal: no image specified")
		os.Exit(1)
	}
	W, H := getTermSize()
	var img Immg
	var err error
	for _, path := range os.Args[1:] {
		if strings.HasPrefix(path, "-") && len(path) == 2 && '1' <= path[1] && path[1] <= '9' {
			W *= (int(path[1]) - '0')
			W /= 10
			H *= (int(path[1]) - '0')
			H /= 10
			continue
		}
		err = img.OpenImg(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening image: %s\n", path)
			continue
		}
		img.FitSize(W, H)
		img.Print()
	}
}
