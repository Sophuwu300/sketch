package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"os/signal"
	"regexp"

	"github.com/nfnt/resize"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

func getTermSize() (float64, float64) {
	W, H, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || !(W > 0 && H > 0) {
		fmt.Fprintln(os.Stderr, "fatal: could not get terminal size")
		os.Exit(1)
	}
	return float64(W), float64(H)
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
	Img   image.Image
	X     int
	Y     int
	Scale float64
}

func (img *Immg) SetScale(s string) {
	for i := len(s) - 1; 0 <= i; i-- {
		if '9' >= s[i] && '0' <= s[i] {
			img.Scale = (img.Scale + float64(s[i]-'0')) / 10
		}
	}
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

func F(n int) float64 {
	return float64(n)
}

func fxf(a float64) func(b float64) uint {
	return func(b float64) uint {
		return uint(math.Round(a * b))
	}
}

func (img *Immg) Print() {
	w, h := getTermSize()
	b := img.Img.Bounds()
	x, y := F(b.Dx()), F(b.Dy())
	h = h*2 - 1
	if x > w {
		y, x = y*w/x, w
	}
	if y > h {
		y, x = h, x*h/y
	}
	u := fxf(img.Scale)
	img.Img = resize.Resize(u(x), u(y), img.Img, resize.MitchellNetravali)
	img.Y = img.Img.Bounds().Min.Y
	for img.Y < img.Img.Bounds().Max.Y {
		fmt.Print("\r")
		img.ForX()
		fmt.Println()
		img.Y += 2
	}
}

func (img *Immg) Doalp() {
	r, g, b, _ := img.Img.At(img.Img.Bounds().Min.X+img.X/2, img.Y+img.X%2).RGBA()
	fmt.Printf("\033[%d8;2;%d;%d;%dm", 3+(img.X%2), r>>8, g>>8, b>>8)
}

var scaleFlag = regexp.MustCompile(`-[0-9]+`)

func init() {
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "--help" {
			fmt.Println("Usage: imgcat [options] [image files]")
			fmt.Println("Options:")
			fmt.Println("  -h, --help       Show this help message")
			fmt.Println("  -sN              Set scale to N/10 (e.g., -s8 sets scale to 0.8)")
			fmt.Println()
			fmt.Println("To use interactive gallery mode, run without any arguments.\nOr with a directory as the only argument.")
			os.Exit(0)
		}
		if arg == "--version" {
			fmt.Printf("sketch: %s\nversion: %s\n", os.Args[0], SKETCH_VERSION)
			os.Exit(0)
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		gallery()
		return
	}
	if len(os.Args) == 2 {
		st, err := os.Stat(os.Args[1])
		if err == nil && st.IsDir() {
			err = os.Chdir(os.Args[1])
			if err != nil {
				fmt.Fprintln(os.Stderr, "fatal: could not change directory")
				os.Exit(1)
			}
			gallery()
			return
		}
	}
	var img Immg
	img.Scale = 1
	var err error
	for _, path := range os.Args[1:] {
		if scaleFlag.MatchString(path) {
			img.SetScale(path)
			continue
		}
		err = img.OpenImg(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening image: %s\n", path)
			continue
		}
		img.Print()
		if len(os.Args) > 2 {
			fmt.Println()
		}
	}
}

func gallery() {
	g := Gallary{}
	err := g.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer g.deferred()
	g.Print()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, unix.SIGWINCH)
	go func() {
		for range ch {
			g.Print()
		}
	}()

	buf := make([]byte, 4)
	var n int
	for {
		n, err = os.Stdin.Read(buf)
		if err != nil {
			break
		}
		g.HandleInput(string(buf[:n]))

	}
}

func (g *Gallary) HandleInput(s string) {
	switch s {
	case "q", "Q", "\x03", "\x04":
		g.Quit()
	case "\x1b[C":
		g.Next()
	case "\x1b[D":
		g.Prev()
	case "\x1b[A", "+":
		g.ChangeScale(0.05)
	case "\x1b[B", "-":
		g.ChangeScale(-0.05)
	case "0":
		g.Img.Scale = 1
		g.Print()
	case "\r", " ":
		g.ShowMeta = !g.ShowMeta
		g.Print()
	}
}

func (g *Gallary) ChangeScale(n float64) {
	g.Img.Scale += n
	if g.Img.Scale < 0.05 {
		g.Img.Scale = 0.05
	}
	if g.Img.Scale > 0.95 {
		g.Img.Scale = 0.95
	}
	g.Print()
}

func (g *Gallary) Next() {
	g.Index++
	if g.Index >= len(g.Paths) {
		g.Index = 0
	}
	g.Print()
}
func (g *Gallary) Prev() {
	g.Index--
	if g.Index < 0 {
		g.Index = len(g.Paths) - 1
	}
	g.Print()
}

func (g *Gallary) Quit() {
	g.deferred()
	os.Exit(0)
}

type Gallary struct {
	Paths    []string
	Index    int
	Img      Immg
	ShowMeta bool
}

var oldState *term.State

func (g *Gallary) Load() error {
	de, err := os.ReadDir(".")
	if err != nil {
		return fmt.Errorf("fatal: could not read current directory")
	}
	paths := []string{}
	rx := regexp.MustCompile(`(?i)\.(png|jpe?g)$`).MatchString
	for _, entry := range de {
		if entry.IsDir() {
			continue
		}
		if entry.Type().IsRegular() && rx(entry.Name()) {
			paths = append(paths, entry.Name())
		}
	}
	if len(paths) == 0 {
		return fmt.Errorf("fatal: no images found in current directory")
	}
	g.Paths = paths
	g.Index = 0
	g.Img.Scale = 0.8
	err = g.Img.OpenImg(g.Paths[g.Index])
	if err != nil {
		return fmt.Errorf("error opening image: %s", g.Paths[g.Index])
	}

	oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Print("\033[?1049h")
	return nil
}
func (g *Gallary) deferred() {
	if r := recover(); r != nil {
		term.Restore(int(os.Stdin.Fd()), oldState) // restore terminal when program exits
		fmt.Print("\033[?1049l")
		panic(r)
	} else {
		term.Restore(int(os.Stdin.Fd()), oldState) // restore terminal when program exits
		fmt.Print("\033[?1049l")
	}
}

func humanBytes(n int64) string {
	f := float64(n)
	prefix := "-KMGTPE"
	i := 0
	for f >= 1024 && i < len(prefix)-1 {
		i++
		f /= 1024
	}
	if i == 0 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.1f %ciB", f, prefix[i])
}

func (g *Gallary) Print() {
	err := g.Img.OpenImg(g.Paths[g.Index])
	if err != nil {
		g.Paths = append(g.Paths[:g.Index], g.Paths[g.Index+1:]...)
		g.Next()
		return
	}
	s := ""
	if g.ShowMeta {
		resx := g.Img.Img.Bounds().Dx()
		resy := g.Img.Img.Bounds().Dy()
		name := g.Paths[g.Index]
		st, _ := os.Stat(name)
		size := st.Size()

		s = fmt.Sprintf("%s (%dx%d, %s)", name, resx, resy, humanBytes(size))
	}
	fmt.Println("\033[2J\033[1;1H")
	g.Img.Print()
	y := g.Img.Y/2 + g.Img.Y%2 + 2
	if g.ShowMeta {
		fmt.Printf("\033[%d;1H%s", y, s)
	} else {
		fmt.Printf("\033[%d;1H%s\t%s\t%s\t%s", y, "[Q] quit", "[←] [→] navigation", "[↑] [↓] scale", "[Space] toggle info")
	}
}

var uwu map[string]func(s []string)
