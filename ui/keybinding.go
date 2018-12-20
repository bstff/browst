package ui

import (
	"browst/gocui"
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

func kbinput(g *gocui.Gui) error {
	if err := g.SetKeybinding(inputView, gocui.KeyEnter, gocui.ModNone, input); err != nil {
		return err
	}
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

func kbmouse(g *gocui.Gui) error {
	err := g.SetKeybinding(renderView, gocui.MouseLeft, gocui.ModNone, clickRender)
	if err != nil {
		return err
	}

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
