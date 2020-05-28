package main

import (
	"log"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appID = "com.github.ianprime0509.gjisho"

type hideable interface {
	Hide()
}

var aboutDialog *gtk.AboutDialog
var searchRevealer *gtk.Revealer

var appComponents = map[string]interface{}{
	"aboutDialog":    &aboutDialog,
	"searchRevealer": &searchRevealer,
}

var signals = map[string]interface{}{
	"hideWidget":    hideable.Hide,
	"inhibitNext":   func() bool { return true },
	"searchChanged": searchChanged,
	"toggleSearch": func() {
		searchRevealer.SetRevealChild(!searchRevealer.GetRevealChild())
	},
}

// LaunchGUI launches the application user interface, passing the given
// arguments to GTK. It does not return an error; if any errors occur here, the
// program will terminate.
func LaunchGUI(args []string) {
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatalf("could not create application: %v", err)
	}

	_, err = app.Connect("activate", onActivate, app)
	if err != nil {
		log.Fatalf("could not connect activation signal: %v", err)
	}

	os.Exit(app.Run(args))
}

func onActivate(app *gtk.Application) {
	builder, err := gtk.BuilderNewFromFile("gjisho.glade")
	if err != nil {
		log.Fatalf("could not create application builder: %v", err)
	}
	windowObj, err := builder.GetObject("appWindow")
	if err != nil {
		log.Fatalf("could not get application window: %v", err)
	}
	getAppComponents(builder)
	builder.ConnectSignals(signals)
	window := windowObj.(*gtk.ApplicationWindow)
	app.AddWindow(window)

	aboutAction := glib.SimpleActionNew("about", nil)
	aboutAction.Connect("activate", func() { aboutDialog.Present() })
	app.AddAction(aboutAction)

	window.Show()
}

func getAppComponents(builder *gtk.Builder) {
	for name, ptr := range appComponents {
		comp, err := builder.GetObject(name)
		if err != nil {
			log.Fatalf("could not get application component %v: %v", name, err)
		}
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(comp))
	}
}

func searchChanged(entry *gtk.SearchEntry) {

}
