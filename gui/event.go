package gui

import (
	"github.com/gotk3/gotk3/gdk"
)

// key is a keyboard key along with any modifiers.
type key struct {
	key       uint
	modifiers uint
}

// keyMap is a map of keys to corresponding event handlers.
type keyMap map[key]func()

// handle handles the given event.
func (m keyMap) handle(ev *gdk.EventKey) {
	k := key{ev.KeyVal(), ev.State() & gdk.GDK_MODIFIER_MASK}
	if h := m[k]; h != nil {
		h()
	}
}

// adaptKeyHandler returns a function suitable for handling a key press or
// release event using the given key map.
func adaptKeyHandler(m keyMap) func(interface{}, *gdk.Event) {
	return func(_ interface{}, ev *gdk.Event) {
		m.handle(&gdk.EventKey{Event: ev})
	}
}

// button is a mouse button along with any modifiers.
type button struct {
	button    uint
	modifiers uint
}

// buttonMap is a map of mouse buttons to corresponding event handlers.
type buttonMap map[button]func()

// handle handles the given event.
func (m buttonMap) handle(ev *gdk.EventButton) {
	b := button{ev.ButtonVal(), ev.State() & gdk.GDK_MODIFIER_MASK}
	if h := m[b]; h != nil {
		h()
	}
}

// adaptButtonHandler returns a function suitable for handling a button press or
// release event using the given button map.
func adaptButtonHandler(m buttonMap) func(interface{}, *gdk.Event) {
	return func(_ interface{}, ev *gdk.Event) {
		m.handle(&gdk.EventButton{Event: ev})
	}
}
