package main

import (
	"github.com/gotk3/gotk3/gdk"
)

// Key is a keyboard key along with any modifiers.
type Key struct {
	Key       uint
	Modifiers uint
}

// KeyMap is a map of keys to corresponding event handlers.
type KeyMap map[Key]func()

// Handle handles the given event.
func (m KeyMap) Handle(ev *gdk.EventKey) {
	key := Key{ev.KeyVal(), ev.State() & gdk.GDK_MODIFIER_MASK}
	if h := m[key]; h != nil {
		h()
	}
}

// AdaptKeyHandler returns a function suitable for handling a key press or
// release event using the given key map.
func AdaptKeyHandler(m KeyMap) func(interface{}, *gdk.Event) {
	return func(_ interface{}, ev *gdk.Event) {
		m.Handle(&gdk.EventKey{Event: ev})
	}
}

// Button is a mouse button along with any modifiers.
type Button struct {
	Key       uint
	Modifiers uint
}

// ButtonMap is a map of mouse buttons to corresponding event handlers.
type ButtonMap map[Button]func()

// Handle handles the given event.
func (m ButtonMap) Handle(ev *gdk.EventButton) {
	but := Button{ev.ButtonVal(), ev.State() & gdk.GDK_MODIFIER_MASK}
	if h := m[but]; h != nil {
		h()
	}
}

// AdaptButtonHandler returns a function suitable for handling a button press or
// release event using the given button map.
func AdaptButtonHandler(m ButtonMap) func(interface{}, *gdk.Event) {
	return func(_ interface{}, ev *gdk.Event) {
		m.Handle(&gdk.EventButton{Event: ev})
	}
}
