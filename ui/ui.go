package ui

import (
	"browst/gocui"
	"browst/sdump"
	"fmt"
	"log"
	"time"
)

type HandlerInput func(ev InputEvent)

const (
	renderAnchorX = 2
	renderAnchorY = 2

	inputView  = "input"
	renderView = "render"

	viewMoveDelta = 1
)

var (
	done = make(chan struct{})

	handlerFunc HandlerInput

	renderX = 0
	renderY = 0

	imgWidth  = 1280
	imgHeight = 720

	renderXMax = 0
	renderYMax = 0
)

func SetHandlerInput(handler HandlerInput) {
	handlerFunc = handler
}

func sixelCropPrint(img_data []byte, l, t, r, b int) {

	cursorTo := fmt.Sprintf("\033[%d;%dH", renderAnchorX, renderAnchorY)

	buf := sdump.EncodeCropImage(img_data, l, t, r, b)
	fmt.Print(cursorTo, string(buf))
}

func updateView(g *gocui.Gui, ch chan []byte) {
	for {
		select {
		case <-done:
			return
		case data := <-ch:

			g.Update(func(g *gocui.Gui) error {
				v, err := g.View(renderView)
				if err != nil {
					return err
				}
				x, y := v.Size()
				xPix, yPix := getCellPixel()

				sixelCropPrint(data,
					renderX*xPix, renderY*yPix, (x)*xPix, (y)*yPix)

				return nil
			})

		case <-time.After(time.Millisecond * 10):
		}
	}
}

func Run(ch chan []byte, w, h int) {
	imgWidth, imgHeight = w, h

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	go updateView(g, ch)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
