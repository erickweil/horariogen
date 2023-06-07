package genetic

import (
	//"math"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/erickweil/horariogen/meucanvas"
	"github.com/tfriedel6/canvas"
)

var imageCanvas *canvas.Image
var imagePixels [][]meucanvas.Pixel
var populacao *Populacao

func hillFit(criatura *Cromossomo) float64 {
	imgX := criatura.genoma[0]
	imgY := criatura.genoma[1]

	pixelsw := len(imagePixels)
	pixelsh := len(imagePixels[0])

	if imgX < 0 || imgX >= pixelsw || imgY < 0 || imgY >= pixelsh {
		return 0.0
	}

	var p meucanvas.Pixel = imagePixels[imgY][imgX]
	return 100 + float64(p.R + p.G + p.B) / 3.0
}

func hillInit(cv *canvas.Canvas) {
	pixels, imageObj, err := meucanvas.GetPixels("hill.png")
			
	if err != nil {	panic(err) }

	img, err := cv.LoadImage(imageObj)

	if err != nil {	panic(err) }

	imageCanvas = img
	imagePixels = pixels

	w := len(imagePixels)
	h := len(imagePixels[0])

	fmt.Println("Carregou imagem:",w,h," [0,0] -> ",imagePixels[0][0])

	rand.Seed(time.Now().UnixNano())
	populacao = CriarPopulacao(
		50,2,
		INTEIRO,
		ROLETA,
		ROLETA,
		PERTURBACAO,
		RECOMBINACAO_ARITMETICA,
		func(criatura *Cromossomo) {
		criatura.genoma[0] = rand.Intn(w)
		criatura.genoma[1] = rand.Intn(h)

		//criatura.genoma[0] = rand.Intn(100)
		//criatura.genoma[1] = rand.Intn(100)
	}, hillFit)
}

func hillLoop(cv *canvas.Canvas, w, h, minwh float64) {
	pixelsw := len(imagePixels)
	//pixelsh := len(imagePixels[0])
	imagex, imagey := meucanvas.PosToScreen(-1.0,-1.0,w,h)
	imagew := meucanvas.DimToScreen(2.0,minwh)
	cv.DrawImage(imageCanvas,imagex,imagey,imagew,imagew)

	nCriaturas := len(populacao.criaturas)

	cv.SetStrokeStyle("#FFFFFF")
	cv.SetFillStyle("#FFFFFF")
	cv.SetFont("Arimo-Regular.ttf", 22)
	cv.FillText(fmt.Sprintf("Tamanho da população:%d",nCriaturas),40.0,40.0)
	
	cv.SetFillStyle("#FF0000")
	for i := 0; i < nCriaturas; i++ {
		criatura := &populacao.criaturas[i]
		
		posx, posy := meucanvas.PosToScreen(
			float64(criatura.genoma[0])/float64(pixelsw) * 2.0 - 1.0,
			float64(criatura.genoma[1])/float64(pixelsw) * 2.0 - 1.0,
			w,h)
		cv.BeginPath()
		cv.Ellipse(posx,posy,3,3,0.0,0.0,math.Pi*2,false)
		cv.ClosePath()
		cv.Fill()
	}

	populacao = SimularGeracao(populacao,hillFit)
}

func ExecHillClimbing() {
	meucanvas.RunCanvas(
		hillInit,
		hillLoop,
	)
}