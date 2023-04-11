package main

import (
	"image"
	"log"
	"os"
	"sync"
	"time"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

const Width = 800
const Height = 600

func main() {
	go func() {
		w := app.NewWindow()
		err := run(w)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops
	i := 0.0
	start := time.Now()

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			rect := image.Rect(0, 0, Width, Height)
			pix := make([]uint8, rect.Dx()*rect.Dy())
			mandelbrot(-0.7+i, 2.1-i, -1.2+i, 1.5-i, 256, pix)
			myImage := image.Gray{
				Pix:    pix,
				Stride: rect.Dx(),
				Rect:   rect,
			}
			if i < 0.8 {
				op.InvalidateOp{}.Add(&ops)
			} else {
				elapsed := time.Since(start)
				log.Printf("took %s", elapsed)
			}
			i += 0.01
			drawImage(&ops, &myImage)
			e.Frame(gtx.Ops)
		}
	}
}

func drawImage(ops *op.Ops, img image.Image) {
	imageOp := paint.NewImageOp(img)
	imageOp.Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(1, 1))).Add(ops)
	paint.PaintOp{}.Add(ops)
}

func mandelbrot(minX float64, maxX float64, minY float64, maxY float64, maxIteration int, pixels []uint8) {
	dx := (maxX - minX) / Width
	dy := (maxY - minY) / Height
	yk := minY
	var waitGroup sync.WaitGroup
	for y := 0; y < Height; y++ {
		waitGroup.Add(1)
		scanY := y * Width
		xk := minX
		go mandelbrotLine(scanY, dx, dy, xk, yk, maxIteration, pixels, &waitGroup)
		yk += dy
	}
	waitGroup.Wait()
}

func mandelbrotLine(scanY int, dx float64, dy float64, xk float64, yk float64, maxIteration int, pixels []uint8, waitGroup *sync.WaitGroup) {
	for x := 0; x < Width; x++ {
		iteration := 0
		var xx float64 = 0
		var yy float64 = 0
		var valueX float64 = 0
		var valueY float64 = 0
		for (iteration < maxIteration) && ((xx + yy) < 4) {
			valueY = 2.0*valueX*valueY - yk
			valueX = xx - yy - xk
			xx = valueX * valueX
			yy = valueY * valueY
			iteration += 1
		}
		if iteration < maxIteration {
			pixels[scanY+x] = (uint8)(256 - iteration%256)
		} else {
			pixels[scanY+x] = 0
		}
		xk += dx
	}
	waitGroup.Done()
}
