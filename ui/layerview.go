package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/wagoodman/docker-image-explorer/image"
	"github.com/lunixbochs/vtclean"
)

type LayerView struct {
	Name       string
	gui        *gocui.Gui
	view       *gocui.View
	header     *gocui.View
	LayerIndex int
	Layers     []*image.Layer
}

func NewLayerView(name string, gui *gocui.Gui, layers []*image.Layer) (layerview *LayerView) {
	layerview = new(LayerView)

	// populate main fields
	layerview.Name = name
	layerview.gui = gui
	layerview.Layers = layers

	return layerview
}

func (view *LayerView) Setup(v *gocui.View, header *gocui.View) error {

	// set view options
	view.view = v
	view.view.Editable = false
	view.view.Wrap = false
	//view.view.Highlight = true
	//view.view.SelBgColor = gocui.ColorGreen
	//view.view.SelFgColor = gocui.ColorBlack
	view.view.Frame = false

	view.header = header
	view.header.Editable = false
	view.header.Wrap = false
	view.header.Frame = false

	// set keybindings
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyArrowDown, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.CursorDown() }); err != nil {
		return err
	}
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyArrowUp, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.CursorUp() }); err != nil {
		return err
	}
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyCtrlL, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.setCompareMode(CompareLayer) }); err != nil {
		return err
	}
	if err := view.gui.SetKeybinding(view.Name, gocui.KeyCtrlA, gocui.ModNone, func(*gocui.Gui, *gocui.View) error { return view.setCompareMode(CompareAll) }); err != nil {
		return err
	}


	headerStr := fmt.Sprintf(image.LayerFormat, "Image ID", "Size", "Command")
	fmt.Fprintln(view.header, Formatting.Header(vtclean.Clean(headerStr, false)))

	return view.Render()
}

func (view *LayerView) setCompareMode(compareMode CompareType) error {
	Views.Tree.CompareMode = compareMode
	view.Render()
	return Views.Tree.setLayer(Views.Tree.CompareStopIndex)
}

func (view *LayerView) Render() error {
	view.gui.Update(func(g *gocui.Gui) error {
		view.view.Clear()
		for revIdx := len(view.Layers) - 1; revIdx >= 0; revIdx-- {
			layer := view.Layers[revIdx]
			idx := (len(view.Layers)-1) - revIdx

			layerStr := layer.String()
			if idx == 0 {
				// TODO: add size
				layerStr = fmt.Sprintf(image.LayerFormat, layer.History.ID[0:25], "", "FROM "+layer.Id())
			}

			if idx == view.LayerIndex {
				fmt.Fprintln(view.view, Formatting.StatusBar(layerStr))
			} else {
				fmt.Fprintln(view.view, layerStr)
			}

		}
		return nil
	})
	// todo: blerg
	return nil
}

func (view *LayerView) CursorDown() error {
	if view.LayerIndex < len(view.Layers) {
		err := CursorDown(view.gui, view.view)
		if err == nil {
			view.LayerIndex++
			view.Render()
			Views.Tree.setLayer(view.LayerIndex)
		}
	}
	return nil
}

func (view *LayerView) CursorUp() error {
	if view.LayerIndex > 0 {
		err := CursorUp(view.gui, view.view)
		if err == nil {
			view.LayerIndex--
			view.Render()
			Views.Tree.setLayer(view.LayerIndex)
		}
	}
	return nil
}

func (view *LayerView) KeyHelp() string {
	return  Formatting.Control("[^L]") + ": Layer Changes " +
		Formatting.Control("[^A]") + ": All Changes "
}