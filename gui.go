package main

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/tatoeba"
)

const appID = "com.github.ianprime0509.gjisho"

var aboutDialog *gtk.AboutDialog
var moreInfoRevealer *gtk.Revealer

var searchResults = new(SearchResultList)
var search = &Search{results: searchResults}
var kanjiList = new(KanjiList)
var kanjiDetails = new(KanjiDetails)
var exampleList = new(ExampleList)
var entryDisplay = &EntryDisplay{kanjiList: kanjiList, exampleList: exampleList}
var navigation = &EntryNavigation{disp: entryDisplay}

var appComponents = map[string]interface{}{
	"aboutDialog":                 &aboutDialog,
	"backButton":                  &navigation.backButton,
	"entryDetailsLabel":           &entryDisplay.detailsLabel,
	"examplesList":                &exampleList.list,
	"forwardButton":               &navigation.forwardButton,
	"kanjiDetailsCharacterLabel":  &kanjiDetails.charLabel,
	"kanjiDetailsDictRefsLabel":   &kanjiDetails.dictRefsLabel,
	"kanjiDetailsReadingMeanings": &kanjiDetails.readingMeanings,
	"kanjiDetailsSubtitleLabel":   &kanjiDetails.subtitleLabel,
	"kanjiDetailsQueryCodesLabel": &kanjiDetails.queryCodesLabel,
	"kanjiDetailsWindow":          &kanjiDetails.window,
	"kanaWritingsLabel":           &entryDisplay.kanaWritingsLabel,
	"kanjiList":                   &kanjiList.list,
	"kanjiWritingsLabel":          &entryDisplay.kanjiWritingsLabel,
	"moreInfoRevealer":            &moreInfoRevealer,
	"primaryKanaLabel":            &entryDisplay.primaryKanaLabel,
	"primaryKanjiLabel":           &entryDisplay.primaryKanjiLabel,
	"searchEntry":                 &search.entry,
	"searchRevealer":              &search.revealer,
	"searchResults":               &searchResults.list,
	"searchToggleButton":          &search.toggle,
}

var signals = map[string]interface{}{
	"activateLink": func(_ *gtk.Label, uri string) bool {
		url, err := url.Parse(uri)
		if err != nil {
			log.Printf("Invalid URL: %v", uri)
			return true
		}
		return navigation.FollowLink(url)
	},
	"examplesEdgeReached": func(_ *gtk.ScrolledWindow, pos gtk.PositionType) {
		if pos == gtk.POS_BOTTOM {
			exampleList.ShowMore()
		}
	},
	"hideWidget":  func(w interface{ Hide() }) { w.Hide() },
	"inhibitNext": func() bool { return true },
	"kanjiListRowActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		kanjiDetails.Display(kanjiList.kanji[row.GetIndex()])
		kanjiDetails.Present()
	},
	"moreInfoToggle": func() {
		moreInfoRevealer.SetRevealChild(!moreInfoRevealer.GetRevealChild())
	},
	"navigateBack":    func() { navigation.GoBack() },
	"navigateForward": func() { navigation.GoForward() },
	"searchChanged": func(entry *gtk.SearchEntry) {
		query, _ := entry.GetText()
		search.Search(query)
	},
	"searchEntryKeyPress": func(_ interface{}, ev *gdk.Event) {
		keyEv := &gdk.EventKey{Event: ev}
		if keyEv.KeyVal() == gdk.KEY_Escape {
			search.Deactivate()
		}
	},
	"searchResultsEdgeReached": func(_ *gtk.ScrolledWindow, pos gtk.PositionType) {
		if pos == gtk.POS_BOTTOM {
			searchResults.ShowMore()
		}
	},
	"searchResultsRowSelected": func() {
		sel := searchResults.Selected()
		if sel == nil {
			return
		}
		navigation.GoTo(sel.ID)
	},
	"searchToggle": search.Toggle,
	"windowButtonPress": func(_ interface{}, ev *gdk.Event) {
		buttonEv := &gdk.EventButton{Event: ev}
		switch buttonEv.Button() {
		case 8:
			navigation.GoBack()
		case 9:
			navigation.GoForward()
		}
	},
	"windowKeyPress": func(_ interface{}, ev *gdk.Event) {
		keyEv := &gdk.EventKey{Event: ev}
		if keyEv.KeyVal() == gdk.KEY_f && keyEv.State()&gdk.GDK_CONTROL_MASK != 0 {
			search.Activate()
		}
	},
}

var dict *jmdict.JMdict
var kanjiDict *kanjidic.Kanjidic
var exampleDict *tatoeba.Tatoeba

// LaunchGUI launches the application user interface, passing the given
// arguments to GTK. It does not return an error; if any errors occur here, the
// program will terminate.
func LaunchGUI(args []string) {
	db, err := sql.Open("sqlite3", "gjisho.sqlite")
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}

	dict, err = jmdict.New(db)
	if err != nil {
		log.Fatalf("Could not open JMdict handler: %v", err)
	}

	kanjiDict, err = kanjidic.New(db)
	if err != nil {
		log.Fatalf("Could not open Kanjidic handler: %v", err)
	}

	exampleDict, err = tatoeba.New(db)
	if err != nil {
		log.Fatalf("Could not open Tatoeba handler: %v", err)
	}

	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	if err != nil {
		log.Fatalf("Could not create application: %v", err)
	}

	_, err = app.Connect("activate", onActivate, app)
	if err != nil {
		log.Fatalf("Could not connect activation signal: %v", err)
	}

	os.Exit(app.Run(args))
}

func onActivate(app *gtk.Application) {
	builder, err := gtk.BuilderNewFromFile("gjisho.glade")
	if err != nil {
		log.Fatalf("Could not create application builder: %v", err)
	}
	windowObj, err := builder.GetObject("appWindow")
	if err != nil {
		log.Fatalf("Could not get application window: %v", err)
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
			log.Fatalf("Could not get application component %v: %v", name, err)
		}
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(comp))
	}
}
