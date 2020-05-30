package main

import (
	"log"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

const appID = "com.github.ianprime0509.gjisho"

type hideable interface {
	Hide()
}

var aboutDialog *gtk.AboutDialog
var searchEntry *gtk.SearchEntry
var searchRevealer *gtk.Revealer
var searchResults *gtk.ListBox
var searchToggleButton *gtk.ToggleButton

var appComponents = map[string]interface{}{
	"aboutDialog":        &aboutDialog,
	"searchEntry":        &searchEntry,
	"searchRevealer":     &searchRevealer,
	"searchResults":      &searchResults,
	"searchToggleButton": &searchToggleButton,
}

var signals = map[string]interface{}{
	"hideWidget":    hideable.Hide,
	"inhibitNext":   func() bool { return true },
	"searchChanged": searchChanged,
	"searchEntryKeyPress": func(_ interface{}, ev *gdk.Event) {
		keyEv := &gdk.EventKey{Event: ev}
		if keyEv.KeyVal() == gdk.KEY_Escape {
			stopSearch()
		}
	},
	"toggleSearch": func() {
		searchRevealer.SetRevealChild(!searchRevealer.GetRevealChild())
	},
	"windowKeyPress": func(_ interface{}, ev *gdk.Event) {
		keyEv := &gdk.EventKey{Event: ev}
		if keyEv.KeyVal() == gdk.KEY_f && keyEv.State()&gdk.GDK_CONTROL_MASK != 0 {
			startSearch()
		}
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

	css, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatalf("could not create CSS provider: %v", err)
	}
	if err := css.LoadFromPath("gui.css"); err != nil {
		log.Fatalf("could not load CSS: %v", err)
	}
	gtk.AddProviderForScreen(window.GetScreen(), css, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

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
	query, _ := entry.GetText()
	go func() {
		entries, err := dict.Lookup(query)
		if err == nil {
			glib.IdleAdd(func() { showLookupEntries(entries) })
		} else {
			log.Printf("Lookup query error: %v", err)
		}
	}()
}

func showLookupEntries(entries []LookupEntry) {
	searchResults.GetChildren().Foreach(func(item interface{}) {
		searchResults.Remove(item.(gtk.IWidget))
	})

	for i, entry := range entries {
		if i > 50 {
			break
		}
		searchResults.Add(newSearchResult(&entry))
	}
	searchResults.ShowAll()
}

func newSearchResult(entry *LookupEntry) gtk.IWidget {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	ctx, _ := box.GetStyleContext()
	ctx.AddClass("search-result")

	box.Add(newSimpleLabel(entry.Heading, "search-result-heading"))
	if entry.Heading != entry.PrimaryReading {
		box.Add(newSimpleLabel(entry.PrimaryReading, "search-result-reading"))
	}
	box.Add(newSimpleLabel(entry.GlossSummary, "search-result-gloss"))

	return box
}

func newSimpleLabel(text string, classes ...string) *gtk.Label {
	label, _ := gtk.LabelNew(text)
	label.SetXAlign(0)
	label.SetEllipsize(pango.ELLIPSIZE_END)
	ctx, _ := label.GetStyleContext()
	for _, class := range classes {
		ctx.AddClass(class)
	}
	return label
}

func startSearch() {
	searchToggleButton.SetActive(true)
	searchRevealer.SetRevealChild(true)
	searchEntry.GrabFocus()
}

func stopSearch() {
	searchToggleButton.SetActive(false)
	searchRevealer.SetRevealChild(false)
}
