package meucanvas

import (
	"bytes"
	//"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"math"

	//"os"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/glfwcanvas"
)

// [-1.0 -- 1.0] -> [0.0 -- minwh]
func PosToScreen(x float64, y float64, w float64, h float64) (float64, float64) {
	minwh := w
	if minwh > h { minwh = h }

	sx := (x + 1.0) * 0.5 * minwh
	sy := (y + 1.0) * 0.5 * minwh

	if w > minwh {
		return (sx + (w - minwh) * 0.5), sy
	} else if h > minwh {
		return sx, (sy + (h - minwh) * 0.5)
	} else {
		return sx, sy
	}
}

func DimToScreen(x float64, minwh float64) float64 {
	s := x * 0.5 * minwh
	return s
}


type CanvasInit func(cv *canvas.Canvas)
type CanvasLoop func(cv *canvas.Canvas, w float64, h float64, minwh float64) 

func RunCanvas(init CanvasInit,loop CanvasLoop) {
	wnd, cv, err := glfwcanvas.CreateWindow(1280, 720, "Hello")
	if err != nil {
		panic(err)
	}
	defer wnd.Close()

	init(cv)
	wnd.MainLoop(func() {
		w, h := float64(cv.Width()), float64(cv.Height())
		cv.SetFillStyle("#444444")
		cv.FillRect(0, 0, w, h)

		minwh := w
		if minwh > h { minwh = h }

		
		squarex, squarey := PosToScreen(-1.0,-1.0,w,h)
		cv.SetFillStyle("#000000")
		cv.FillRect(squarex, squarey, minwh, minwh)

		loop(cv,w,h,minwh)
	})
}

// https://stackoverflow.com/questions/33186783/get-a-pixel-array-from-from-golang-image-image
// Get the bi-dimensional pixel array
func GetPixels(filename string) ([][]Pixel, image.Image, error) {
	image.RegisterFormat("png","png",png.Decode, png.DecodeConfig)
			
	bytearr, err := ioutil.ReadFile(filename)

	if err != nil {	
        return nil, nil, err
	}

    img, _, err := image.Decode(bytes.NewReader(bytearr))

    if err != nil {
        return nil, img, err
    }

    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y

    var pixels [][]Pixel
    for y := 0; y < height; y++ {
        var row []Pixel
        for x := 0; x < width; x++ {
            row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
        }
        pixels = append(pixels, row)
    }

    return pixels, img, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
    return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
    R int
    G int
    B int
    A int
}

func ExecCanvas() {
	RunCanvas(
		func(cv *canvas.Canvas) {
		},
		func(cv *canvas.Canvas, w, h, minwh float64) {
			elipsex, elipsey := PosToScreen(0.0,0.0,w,h)
			radius := DimToScreen(0.75,minwh)
			cv.SetFillStyle("#F00000")
			cv.SetStrokeStyle("#FF0000")
			cv.BeginPath()
			cv.Ellipse(elipsex,elipsey,radius,radius,0.0,0.0,math.Pi*2,false)
			cv.ClosePath()
			cv.Fill()
		},
	)
}