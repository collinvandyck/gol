package main

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
)

type grid [][]bool

const (
	alive  = true
	dead   = false
	width  = 200
	height = 200
)

var (
	pixelSize  = 2
	img        = newBaseImage(width, height, pixelSize)
	colorAlive = color.RGBA{R: 0, G: 204, B: 255}
	colorDead  = color.Black
	tickDelay  = 10 * time.Millisecond

	initialNoiseFactor = 0.1
)

func main() {
	rand.Seed(time.Now().Unix())
	app := app.New()
	window := app.NewWindow("Game of Life")
	window.SetFullScreen(false)
	window.Resize(fyne.Size{
		Width:  img.Bounds().Size().X,
		Height: img.Bounds().Size().Y,
	})
	image := canvas.NewImageFromImage(img)
	image.FillMode = canvas.ImageFillContain
	canvas := window.Canvas()
	canvas.SetContent(image)
	grid := newGrid(width, height)
	go simulate(grid, image)

	window.ShowAndRun()
}

// liveNeighbors returns the number of neighbors to this point that are alive
func (g grid) liveNeighbors(x, y int) int {
	res := 0
	check := func(x int, y int) {
		if x >= 0 && x < width && y >= 0 && y < height {
			if g[x][y] == alive {
				res++
			}
		}
	}
	check(x-1, y-1)
	check(x-1, y)
	check(x-1, y+1)
	check(x, y-1)
	check(x, y+1)
	check(x+1, y-1)
	check(x+1, y)
	check(x+1, y+1)
	return res
}

func tick(g grid) grid {
	res := g.copy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			neighbors := g.liveNeighbors(x, y)
			switch g[x][y] {
			case alive:
				if !(neighbors >= 2 && neighbors <= 3) {
					res[x][y] = dead
				}
			default:
				if neighbors == 3 {
					res[x][y] = alive
				}
			}
		}
	}
	return res
}

func simulate(g grid, ci *canvas.Image) {
	for range time.Tick(tickDelay) {
		g = tick(g)
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				var c color.Color = colorDead
				if g[x][y] {
					c = colorAlive
				}
				//log.Printf("x:%d y:%d color:%v", x, y, c)
				for psx := 0; psx < pixelSize; psx++ {
					for psy := 0; psy < pixelSize; psy++ {
						px := x*pixelSize + psx
						py := y*pixelSize + psy
						//log.Printf("set x:%d y:%d px:%d py:%d", x, y, px, py)
						img.Set(px, py, c)
						//img.Set(x*pixelSize+i, y*pixelSize+i, c)
					}
				}
				/*
					for i := 0; i < pixelSize; i++ {
						px := x*pixelSize + i
						py := y*pixelSize + i
						log.Printf("set x:%d y:%d px:%d py:%d", x, y, px, py)
						img.Set(x*pixelSize+i, y*pixelSize+i, c)
						//img.Set(x+i, y+i, c)
						//fmt.Println("Set", x+i, y+i)
						/*
							img.Set(x*pixelSize+i, y*pixelSize, c)
							img.Set(x*pixelSize, y*pixelSize+1, c)
						}
				*/
			}
		}
		ci.Refresh()
	}
}

func newGrid(width int, height int) grid {
	res := make(grid, 0)
	for x := 0; x < width; x++ {
		col := make([]bool, height)
		res = append(res, col)
		for y := 0; y < height; y++ {
			if rand.Float64() < initialNoiseFactor {
				col[y] = alive
			}
		}
	}
	return res
}

func (g grid) copy() grid {
	res := make(grid, 0)
	for x := 0; x < width; x++ {
		col := make([]bool, height)
		res = append(res, col)
		for y := 0; y < height; y++ {
			res[x][y] = g[x][y]
		}
	}
	return res
}

func newBaseImage(width int, height int, pixelSize int) *image.RGBA {
	i := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width * pixelSize, Y: height * pixelSize},
	})
	for x := 0; x < i.Bounds().Max.X; x++ {
		for y := 0; y < i.Bounds().Max.Y; y++ {
			i.Set(x, y, colorDead)
		}
	}
	return i
}
