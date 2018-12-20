package ui

import (
	"browst/gocui"
)

func layout(g *gocui.Gui) error {

	x, y := g.Size()

	if v, err := g.SetView(renderView, 0, 0, x-1, y-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true

		delta := 1
		xPix, yPix := getCellPixel()
		renderXMax, renderYMax = imgWidth/xPix-x+delta, imgHeight/yPix-y+delta

		if _, err := g.SetCurrentView(v.Name()); err != nil {
			return err
		}
	}

	if v, err := g.SetView(inputView, 0, y-3, x-1, y-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = true
	}

	return nil
}
