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
var searchToggleButton *gtk.ToggleButton

var searchResults = new(SearchResultList)

var appComponents = map[string]interface{}{
	"aboutDialog":        &aboutDialog,
	"searchEntry":        &searchEntry,
	"searchRevealer":     &searchRevealer,
	"searchResults":      &searchResults.listBox,
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
	"searchResultsEndReached":  func() { searchResults.ShowMore() },
	"searchResultsRowSelected": func() { log.Print(searchResults.Selected()) },
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

// SearchResultList is a list of search results displayed in the GUI.
type SearchResultList struct {
	listBox    *gtk.ListBox
	results    []LookupResult
	nDisplayed int
}

// Selected returns the currently selected search result, or nil if none is
// selected.
func (lst *SearchResultList) Selected() *LookupResult {
	if row := lst.listBox.GetSelectedRow(); row != nil {
		return &lst.results[row.GetIndex()]
	}
	return nil
}

// SetResults sets the currently displayed search results.
func (lst *SearchResultList) SetResults(results []LookupResult) {
	lst.results = results
	lst.listBox.GetChildren().Foreach(func(e interface{}) {
		lst.listBox.Remove(e.(gtk.IWidget))
	})
	lst.nDisplayed = 0
	lst.ShowMore()
}

// ShowMore displays more search results in the list.
func (lst *SearchResultList) ShowMore() {
	maxIndex := lst.nDisplayed + 50
	for ; lst.nDisplayed < len(lst.results) && lst.nDisplayed < maxIndex; lst.nDisplayed++ {
		lst.listBox.Add(newSearchResult(lst.results[lst.nDisplayed]))
	}
	lst.listBox.ShowAll()
}

// newSearchResult creates a search result widget for display.
func newSearchResult(entry LookupResult) gtk.IWidget {
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

func searchChanged(entry *gtk.SearchEntry) {
	query, _ := entry.GetText()
	go func() {
		results, err := dict.Lookup(query)
		if err == nil {
			glib.IdleAdd(func() { searchResults.SetResults(results) })
		} else {
			log.Printf("Lookup query error: %v", err)
		}
	}()
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
