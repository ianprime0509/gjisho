package gui

import (
	"github.com/gotk3/gotk3/gtk"
)

// removeChildren removes all children from the given container.
func removeChildren(c *gtk.Container) {
	c.GetChildren().Foreach(func(item interface{}) {
		c.Remove(item.(gtk.IWidget))
	})
}

// scrollToStart scrolls the given scrolled window vertically to the top and
// horizontally to the start.
func scrollToStart(w *gtk.ScrolledWindow) {
	w.GetVAdjustment().SetValue(w.GetVAdjustment().GetLower())
	w.GetHAdjustment().SetValue(w.GetHAdjustment().GetLower())
}
