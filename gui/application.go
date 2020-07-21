package gui

// #cgo pkg-config: gtk+-3.0
// #include "application.h"
import "C"

import (
	"errors"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

var errNilPtr = errors.New("cgo returned unexpected nil pointer")

// application is a wrapper around GJishoApplication. Much of the code for
// handling it is copied from the implementation of gtk.Application.
type application struct {
	gtk.Application
}

func init() {
	tm := []glib.TypeMarshaler{
		{glib.Type(C.gjisho_application_get_type()), marshalApplication},
	}
	glib.RegisterGValueMarshalers(tm)
}

func marshalApplication(p uintptr) (interface{}, error) {
	c := C.g_value_get_object((*C.GValue)(unsafe.Pointer(p)))
	obj := glib.Take(unsafe.Pointer(c))
	return wrapApplication(obj), nil
}

func wrapApplication(obj *glib.Object) *application {
	am := &glib.ActionMap{obj}
	ag := &glib.ActionGroup{obj}
	return &application{gtk.Application{glib.Application{obj, am, ag}}}
}

func applicationNew() (*application, error) {
	capp := C.gjisho_application_new()
	if capp == nil {
		// I don't really like this style of returning an error rather than just
		// nil, but I want to be consistent with gtk.ApplicationNew
		return nil, errNilPtr
	}
	// Basically a re-implementation of gtk.ApplicationNew
	return wrapApplication(glib.Take(unsafe.Pointer(capp))), nil
}
