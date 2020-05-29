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
var searchResults *gtk.ListBox

var appComponents = map[string]interface{}{
	"aboutDialog":    &aboutDialog,
	"searchRevealer": &searchRevealer,
	"searchResults":  &searchResults,
}

var signals = map[string]interface{}{
	"hideWidget":    hideable.Hide,
	"inhibitNext":   func() bool { return true },
	"searchChanged": searchChanged,
	"toggleSearch": func() {
		searchRevealer.SetRevealChild(!searchRevealer.GetRevealChild())
	},
}

var dict *JMdict

// LaunchGUI launches the application user interface, passing the given
// arguments to GTK. It does not return an error; if any errors occur here, the
// program will terminate.
func LaunchGUI(args []string) {
	var err error
	dict, err = OpenJMdict("jmdict.sqlite")
	if err != nil {
		log.Fatalf("could not open JMdict database: %v", err)
	}

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
	if query, err := entry.GetText(); err == nil {
		go func() {
			entries, err := dict.Lookup(query)
			if err == nil {
				glib.IdleAdd(func() { showLookupEntries(entries) })
			} else {
				log.Printf("Lookup query error: %v", err)
			}
		}()
	} else {
		log.Printf("Error getting search text: %v", err)
	}
}

func showLookupEntries(entries []LookupEntry) {
	searchResults.GetChildren().Foreach(func(item interface{}) {
		searchResults.Remove(item.(gtk.IWidget))
	})

	for i, entry := range entries {
		if i > 50 {
			break
		}

		if label, err := gtk.LabelNew(entry.Heading); err == nil {
			searchResults.Add(label)
		} else {
			log.Printf("Error creating label: %v", err)
		}
	}
	searchResults.ShowAll()
}
