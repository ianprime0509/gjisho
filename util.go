package main

import "github.com/gotk3/gotk3/gtk"

// removeChildren removes all children from the given container.
func removeChildren(lst *gtk.Container) {
	lst.GetChildren().Foreach(func(item interface{}) {
		lst.Remove(item.(gtk.IWidget))
	})
}
