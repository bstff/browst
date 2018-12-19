package ui

import (
	"browst/gocui"
	"strings"
)

func kbcommon(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, switchFocus); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func switchFocus(g *gocui.Gui, v *gocui.View) error {

	viewName := ""

	if v.Name() == inputView {
		viewName = renderView
		g.Cursor = false

	} else if v.Name() == renderView {
		viewName = inputView
		g.Cursor = true
	}

	if _, err := g.SetCurrentView(viewName); err != nil {
		return err
	}
	curView := g.CurrentView()
	curView.SetCursor(0, 0)
	return nil
}

func kbinput(g *gocui.Gui) error {
	if err := g.SetKeybinding(inputView, gocui.KeyEnter, gocui.ModNone, input); err != nil {
		return err
	}
	return nil
}

func input(g *gocui.Gui, v *gocui.View) error {
	vbuf := v.ViewBuffer()
	if strings.Index(vbuf, "http://") != 0 {
		vbuf = "http://" + vbuf
	}

	ev := InputEvent{
		ID:      "navigate",
		Payload: vbuf,
	}
	handlerFunc(ev)
	v.Clear()

	switchFocus(g, v)

	return nil
}

func kbrender(g *gocui.Gui) error {
	if err := g.SetKeybinding(renderView, gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveRenderView(g, v, -viewMoveDelta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveRenderView(g, v, viewMoveDelta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveRenderView(g, v, 0, viewMoveDelta)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding(renderView, gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveRenderView(g, v, 0, -viewMoveDelta)
		}); err != nil {
		return err
	}

	return nil
}

func moveRenderView(g *gocui.Gui, v *gocui.View, dx, dy int) error {

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
		ev := InputEvent{
			ID: "wheel",
			Payload: Mouse{
				Y: scrollY,
			},
		}
		handlerFunc(ev)

	}

	return nil
}

func kbmouse(g *gocui.Gui) error {
	err := g.SetKeybinding(renderView, gocui.MouseLeft, gocui.ModNone, clickRender)
	if err != nil {
		return err
	}

	return nil
}

func clickRender(g *gocui.Gui, v *gocui.View) error {

	if v.Name() != renderView {
		return nil
	}

	x, y := v.Cursor()
	xPix, yPix := getCellPixel()
	l, t := (renderX+x)*xPix, (renderY+y)*yPix
	r, b := l+xPix, t+yPix

	ev := InputEvent{
		ID: "click",
		Payload: Mouse{
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

func keybindings(g *gocui.Gui) error {
	err := kbcommon(g)
	if err != nil {
		return err
	}
	err = kbinput(g)
	if err != nil {
		return err
	}
	err = kbrender(g)
	if err != nil {
		return err
	}
	err = kbmouse(g)
	if err != nil {
		return err
	}
	return nil
}
