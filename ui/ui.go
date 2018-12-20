package ui

import (
	"browst/common"
	"browst/gocui"
	"browst/sdump"
	"fmt"
	"log"
	"strings"
	"time"
)

type HandlerInput func(ev common.Event)

const (
	renderXAnchor = 2
	renderYAnchor = 2

	inputView  = "input"
	renderView = "render"

	viewMoveDelta = 1
)

var (
	tui  *gocui.Gui
	done = make(chan struct{})

	handlerFunc HandlerInput

	renderX = 0
	renderY = 0

	imgWidth  = 1280
	imgHeight = 720

	renderXMax = 0
	renderYMax = 0

	waitInput = false
	waitPosL  = 0
	waitPosT  = 0
	waitPosR  = 0
	waitPosB  = 0
)

func SetHandlerInputFunc(handler HandlerInput) {
	handlerFunc = handler
}

func sixelCropPrint(img_data []byte, l, t, r, b int) {

	cursorTo := fmt.Sprintf("\033[%d;%dH", renderXAnchor, renderYAnchor)

	buf := sdump.EncodeCropImage(img_data, l, t, r, b)
	fmt.Print(cursorTo, string(buf))
}

func waitInputOrNot(i bool, l, t, r, b int) {
	waitInput = i

	waitPosL = l
	waitPosT = t
	waitPosR = r
	waitPosB = b
}

func WaitInput(l, t, r, b int) {
	waitInputOrNot(true, l, t, r, b)
	switchView(tui, inputView, true)
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

	tui = g

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

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func switchView(g *gocui.Gui, name string, showCursor bool) error {
	if _, err := g.SetCurrentView(name); err != nil {
		return err
	}
	curView := g.CurrentView()
	curView.SetCursor(0, 0)

	g.Cursor = showCursor

	return nil
}

func switchFocus(g *gocui.Gui, v *gocui.View) error {
	if waitInput {
		return nil
	}

	viewName := ""
	showCursor := false

	if v.Name() == inputView {
		viewName = renderView
	} else if v.Name() == renderView {
		viewName = inputView
		showCursor = true
	}

	return switchView(g, viewName, showCursor)
}

func input(g *gocui.Gui, v *gocui.View) error {

	vbuf := v.ViewBuffer()

	var ev common.Event
	if waitInput {
		ev = common.Event{
			ID: common.InputWaited,
			Payload: common.Region{
				Cont:   []byte(vbuf),
				Left:   waitPosL,
				Top:    waitPosT,
				Right:  waitPosR,
				Bottom: waitPosB,
			},
		}
		waitInputOrNot(false, 0, 0, 0, 0)

	} else {
		if strings.Index(vbuf, "http://") != 0 {
			vbuf = "http://" + vbuf
		}

		ev = common.Event{
			ID:      common.InputURL,
			Payload: vbuf,
		}
	}

	handlerFunc(ev)
	v.Clear()

	switchFocus(g, v)

	return nil
}

func moveRenderView(g *gocui.Gui, v *gocui.View, dx, dy int) error {
	if waitInput {
		return nil
	}

	x := renderX + dx
	y := renderY + dy
	if x > -1 && x < renderXMax {
		renderX = x
	}
	if y > -1 && y < renderYMax {
		renderY = y
	}

	_, scrollY := getCellPixel()

	if y < 0 || y > renderYMax-1 {
		if y < 0 {
			scrollY = -scrollY
		}
		ev := common.Event{
			ID: common.InputMouseWheel,
			Payload: common.Region{
				Y: scrollY,
			},
		}
		handlerFunc(ev)

	}

	return nil
}

func clickRender(g *gocui.Gui, v *gocui.View) error {
	if waitInput {
		return nil
	}

	if v.Name() != renderView {
		return nil
	}

	x, y := v.Cursor()
	xPix, yPix := getCellPixel()
	l, t := (renderX+x)*xPix, (renderY+y)*yPix
	r, b := l+xPix, t+yPix

	ev := common.Event{
		ID: common.InputMouseClick,
		Payload: common.Region{
			X:      x,
			Y:      y,
			Left:   l,
			Top:    t,
			Right:  r,
			Bottom: b,
		},
	}

	handlerFunc(ev)
	return nil
}
