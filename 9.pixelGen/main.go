package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"sync"
)

const (
	width       = 128
	height      = 72
	totalFrames = 100
)

func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a) + (float64(b)-float64(a))*t)
}

func mixColor(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: lerp(a.R, b.R, t),
		G: lerp(a.G, b.G, t),
		B: lerp(a.B, b.B, t),
		A: 255,
	}
}

func renderFrame(frame int) *image.RGBA {
	img := image.NewRGBA(
		image.Rect(0, 0, width, height),
	)

	// 0.0 → 1.0 over the animation
	t := float64(frame) / float64(totalFrames-1)

	// --------------------------------------------------
	// SKY
	// --------------------------------------------------

	dayTop := color.RGBA{40, 80, 140, 255}
	dayBottom := color.RGBA{100, 150, 190, 255}

	sunsetTop := color.RGBA{180, 70, 70, 255}
	sunsetBottom := color.RGBA{245, 130, 70, 255}

	nightTop := color.RGBA{10, 15, 40, 255}
	nightBottom := color.RGBA{35, 30, 55, 255}

	for y := 0; y < 45; y++ {
		heightT := float64(y) / 45.0

		var top, bottom color.RGBA

		if t < 0.55 {
			// Day → sunset
			p := t / 0.55

			top = mixColor(dayTop, sunsetTop, p)
			bottom = mixColor(dayBottom, sunsetBottom, p)
		} else {
			// Sunset → night
			p := (t - 0.55) / 0.45

			top = mixColor(sunsetTop, nightTop, p)
			bottom = mixColor(sunsetBottom, nightBottom, p)
		}

		c := mixColor(top, bottom, heightT)

		for x := 0; x < width; x++ {
			img.Set(x, y, c)
		}
	}

	// --------------------------------------------------
	// STARS
	// --------------------------------------------------

	if t > 0.55 {
		brightness := (t - 0.55) / 0.45

		stars := [][2]int{
			{10, 8},
			{25, 15},
			{38, 7},
			{52, 12},
			{67, 6},
			{78, 17},
			{101, 8},
			{112, 14},
			{120, 5},
			{31, 24},
			{88, 23},
			{117, 27},
		}

		for i, star := range stars {
			// Stars appear gradually.
			if float64(i)/float64(len(stars)) < brightness {
				img.Set(
					star[0],
					star[1],
					color.RGBA{235, 235, 190, 255},
				)
			}
		}
	}

	// --------------------------------------------------
	// SUN
	// --------------------------------------------------

	sunX := 90

	// Sun moves downward.
	sunY := 10 + int(t*40)

	sunColor := color.RGBA{255, 190, 70, 255}

	for y := sunY - 5; y <= sunY+5; y++ {
		for x := sunX - 5; x <= sunX+5; x++ {

			dx := x - sunX
			dy := y - sunY

			if dx*dx+dy*dy <= 25 {
				img.Set(x, y, sunColor)
			}
		}
	}

	// --------------------------------------------------
	// CLOUDS
	// --------------------------------------------------

	cloudColor := color.RGBA{55, 60, 75, 255}

	for i := 0; i < 3; i++ {

		x := (i*45 + frame/2) % 150
		y := 12 + i*6

		for yy := y; yy < y+3; yy++ {
			for xx := x; xx < x+14; xx++ {

				if xx >= 0 && xx < width &&
					yy >= 0 && yy < height {
					img.Set(xx, yy, cloudColor)
				}
			}
		}
	}

	// --------------------------------------------------
	// FAR MOUNTAINS
	// --------------------------------------------------

	farMountain := color.RGBA{50, 65, 80, 255}

	for x := 0; x < width; x++ {

		mountainY := 38 -
			int(5*float64(x%20)/20.0)

		for y := mountainY; y < 52; y++ {
			img.Set(x, y, farMountain)
		}
	}

	// --------------------------------------------------
	// MAIN MOUNTAINS
	// --------------------------------------------------

	mountainColor := color.RGBA{25, 30, 35, 255}

	for x := 0; x < width; x++ {

		// Creates a few peaks.
		mountainY := 45

		if x < 20 {
			mountainY = 45 - x
		} else if x < 35 {
			mountainY = 25 + (x - 20)
		} else if x < 55 {
			mountainY = 40 - (x-35)/2
		} else if x < 75 {
			mountainY = 30 + (x-55)/2
		} else if x < 95 {
			mountainY = 40 - (x-75)/2
		} else {
			mountainY = 30 + (x-95)/2
		}

		for y := mountainY; y < 55; y++ {
			img.Set(x, y, mountainColor)
		}
	}

	// --------------------------------------------------
	// WATER
	// --------------------------------------------------

	waterColor := color.RGBA{20, 55, 70, 255}

	for y := 55; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, waterColor)
		}
	}

	// --------------------------------------------------
	// SUN REFLECTION
	// --------------------------------------------------

	if t < 0.8 {

		reflectionColor := color.RGBA{
			230, 120, 65, 255,
		}

		for y := 56; y < height; y++ {

			width := 8 - (y-56)/4

			if width < 2 {
				width = 2
			}

			for x := sunX - width; x <= sunX+width; x++ {
				img.Set(x, y, reflectionColor)
			}
		}
	}

	// --------------------------------------------------
	// TREES
	// --------------------------------------------------

	treeColor := color.RGBA{12, 18, 15, 255}

	for i := 0; i < 12; i++ {

		x := 5 + i*11
		y := 53

		// trunk
		img.Set(x, y, treeColor)
		img.Set(x, y+1, treeColor)
		img.Set(x, y+2, treeColor)

		// leaves
		for yy := 0; yy < 6; yy++ {
			for xx := -yy; xx <= yy; xx++ {

				px := x + xx
				py := y - yy

				if px >= 0 && px < width {
					img.Set(px, py, treeColor)
				}
			}
		}
	}

	// --------------------------------------------------
	// CABIN
	// --------------------------------------------------

	cabinColor := color.RGBA{25, 25, 22, 255}

	for y := 48; y < 57; y++ {
		for x := 22; x < 37; x++ {
			img.Set(x, y, cabinColor)
		}
	}

	// Roof
	roofColor := color.RGBA{15, 17, 16, 255}

	for i := 0; i < 9; i++ {
		img.Set(22+i, 47-i/3, roofColor)
		img.Set(36-i, 47-i/3, roofColor)
	}

	// Cabin lights turn on during sunset/night.
	if t > 0.6 {

		windowColor := color.RGBA{255, 180, 60, 255}

		for y := 51; y < 54; y++ {
			for x := 25; x < 29; x++ {
				img.Set(x, y, windowColor)
			}

			for x := 31; x < 35; x++ {
				img.Set(x, y, windowColor)
			}
		}
	}

	return img
}

func saveFrame(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
func formatFrameName(frame int) string {
	return fmt.Sprintf("frame_%03d.png", frame)
}

type frameResult struct {
	frame int
	img   *image.RGBA
}

func worker(jobs <-chan int, res_images chan<- frameResult, wg *sync.WaitGroup, completedFrames *int, mu *sync.Mutex, ctx context.Context) {
	defer wg.Done()
	for frame := range jobs {

		select {
		case <-ctx.Done():
			return
		default:

			res_images <- frameResult{frame: frame, img: renderFrame(frame)}
			mu.Lock()
			*completedFrames++
			mu.Unlock()

		}
	}

}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	completedFrames := 0

	bufferSize := flag.Int("buffer", 10, "size of the jobs buffer")
	workerCount := flag.Int("workers", 4, "number of workers")

	flag.Parse()

	jobs := make(chan int, *bufferSize)
	res_images := make(chan frameResult)

	err := os.MkdirAll("frames", 0755)
	if err != nil {
		panic(err)
	}

	rootCtx := context.Background()

	for range *workerCount {
		wg.Add(1)
		go worker(
			jobs,
			res_images,
			&wg,
			&completedFrames,
			&mu,
			rootCtx,
		)
	}

	go func() {
		for frame := range totalFrames {
			jobs <- frame
		}

		close(jobs)

		wg.Wait()
		close(res_images)
	}()

	buffer := make([]frameResult, totalFrames)

	for result := range res_images {
		buffer[result.frame] = result
		println("rendered frame", result.frame)
	}

	fmt.Println("completed:", completedFrames)

	for i, result := range buffer {
		filename := filepath.Join("frames", formatFrameName(i))

		err := saveFrame(result.img, filename)
		if err != nil {
			panic(err)
		}
	}
}
