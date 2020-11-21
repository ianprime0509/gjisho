package gui

import (
	"context"

	"github.com/gotk3/gotk3/gtk"
)

// A fetchOperation represents an operation to fetch some data, where only one
// fetch (for a particular operation) should be active at one time. That is, if
// a fetch is in progress when another starts, the new fetch should supersede
// and cancel the old one.
type fetchOperation struct {
	cancel context.CancelFunc
}

func (o *fetchOperation) start(ctx context.Context) context.Context {
	if o.cancel != nil {
		o.cancel()
	}
	ctx, o.cancel = context.WithCancel(ctx)
	return ctx
}

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
