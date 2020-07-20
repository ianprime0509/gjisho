//go:generate go run github.com/go-bindata/go-bindata/go-bindata -ignore .*~ -nometadata -pkg gui data/

// Package gui contains the GUI interface to GJisho.
package gui

import (
	"log"
	"net/url"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/dictdb"
)

const appID = "com.github.ianprime0509.gjisho"

var aboutDialog *gtk.AboutDialog

var search = &appSearch{
	results:      new(searchResultList),
	resultsKanji: new(searchResultsKanji),
	kanjiInput:   new(kanjiInput),
}
var navigation = &entryNavigation{disp: &entryDisplay{
	kanji:    new(kanjiList),
	examples: new(exampleList),
}}
var kanjiDetails = new(kanjiDetailsModal)
var exampleDetails = new(exampleDetailsModal)

var appComponents = map[string]interface{}{
	"aboutDialog":                      &aboutDialog,
	"backButton":                       &navigation.backButton,
	"entryContentStack":                &navigation.contentStack,
	"entryDetailsLabel":                &navigation.disp.detailsLabel,
	"entryScrolledWindow":              &navigation.disp.scrolledWindow,
	"exampleDetailsEnglishLabel":       &exampleDetails.englishLabel,
	"exampleDetailsJapaneseLabel":      &exampleDetails.japaneseLabel,
	"exampleDetailsScrolledWindow":     &exampleDetails.scrolledWindow,
	"exampleDetailsWindow":             &exampleDetails.window,
	"exampleDetailsWordsList":          &exampleDetails.wordsList,
	"examplesList":                     &navigation.disp.examples.list,
	"examplesScrolledWindow":           &navigation.disp.examples.scrolledWindow,
	"forwardButton":                    &navigation.forwardButton,
	"kanjiDetailsCharacterLabel":       &kanjiDetails.charLabel,
	"kanjiDetailsDictRefsLabel":        &kanjiDetails.dictRefsLabel,
	"kanjiDetailsReadingMeanings":      &kanjiDetails.readingMeanings,
	"kanjiDetailsScrolledWindow":       &kanjiDetails.scrolledWindow,
	"kanjiDetailsStrokeOrder":          &kanjiDetails.strokeOrder,
	"kanjiDetailsSubtitleLabel":        &kanjiDetails.subtitleLabel,
	"kanjiDetailsQueryCodesLabel":      &kanjiDetails.queryCodesLabel,
	"kanjiDetailsWindow":               &kanjiDetails.window,
	"kanjiInputButton":                 &search.kanjiInput.button,
	"kanjiInputButtonIcon":             &search.kanjiInput.buttonIcon,
	"kanjiInputPopover":                &search.kanjiInput.popover,
	"kanjiInputRadicals":               &search.kanjiInput.radicalsBox,
	"kanjiInputRadicalsScrolledWindow": &search.kanjiInput.radicalsScrolledWindow,
	"kanjiInputResults":                &search.kanjiInput.resultsBox,
	"kanjiInputResultsScrolledWindow":  &search.kanjiInput.resultsScrolledWindow,
	"kanaWritingsLabel":                &navigation.disp.kanaWritingsLabel,
	"kanjiList":                        &navigation.disp.kanji.list,
	"kanjiScrolledWindow":              &navigation.disp.kanji.scrolledWindow,
	"kanjiWritingsLabel":               &navigation.disp.kanjiWritingsLabel,
	"moreInfoRevealer":                 &navigation.moreInfoRevealer,
	"moreInfoToggleButton":             &navigation.moreInfoToggleButton,
	"primaryKanaLabel":                 &navigation.disp.primaryKanaLabel,
	"primaryKanjiLabel":                &navigation.disp.primaryKanjiLabel,
	"searchEntry":                      &search.entry,
	"searchRevealer":                   &search.revealer,
	"searchResults":                    &search.results.list,
	"searchResultsKanji":               &search.resultsKanji.box,
	"searchResultsScrolledWindow":      &search.results.scrolledWindow,
	"searchToggleButton":               &search.toggle,
	"strokeOrderScrolledWindow":        &kanjiDetails.strokeOrderScrolledWindow,
	"writingsScrolledWindow":           &navigation.disp.writingsScrolledWindow,
}

var signals = map[string]interface{}{
	"activateLink": func(_ *gtk.Label, uri string) bool {
		url, err := url.Parse(uri)
		if err != nil {
			log.Printf("Invalid URL: %v", uri)
			return true
		}
		search.results.clearSelection()
		return navigation.followLink(url)
	},
	"exampleDetailsWordActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		search.results.clearSelection()
		navigation.goTo(exampleDetails.words[row.GetIndex()].ID)
		exampleDetails.close()
	},
	"examplesEdgeReached": func(_ *gtk.ScrolledWindow, pos gtk.PositionType) {
		if pos == gtk.POS_BOTTOM {
			navigation.disp.examples.showMore()
		}
	},
	"exampleListRowActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		exampleDetails.fetchAndDisplay(navigation.disp.examples.examples[row.GetIndex()])
	},
	"hideWidget":  func(w interface{ Hide() }) { w.Hide() },
	"inhibitNext": func() bool { return true },
	"kanjiListRowActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		kanjiDetails.fetchAndDisplay(navigation.disp.kanji.kanji[row.GetIndex()])
	},
	"kanjiInputButtonToggled": func(b *gtk.ToggleButton) {
		if b.GetActive() {
			search.kanjiInput.display()
		}
	},
	"kanjiInputPopoverClosed": func() {
		search.kanjiInput.button.SetActive(false)
	},
	"kanjiInputResultsChildActivated": func(_ *gtk.FlowBox, child *gtk.FlowBoxChild) {
		search.insertEntryText(search.kanjiInput.results[child.GetIndex()].Literal)
	},
	"moreInfoToggle":  navigation.toggleMoreInfo,
	"navigateBack":    navigation.goBack,
	"navigateForward": navigation.goForward,
	"searchChanged": func(entry *gtk.SearchEntry) {
		query, _ := entry.GetText()
		search.search(query)
	},
	"searchEntryKeyPress": adaptKeyHandler(searchEntryKeyMap),
	"searchResultsEdgeReached": func(_ *gtk.ScrolledWindow, pos gtk.PositionType) {
		if pos == gtk.POS_BOTTOM {
			search.results.showMore()
		}
	},
	"searchResultsKanjiChildActivated": func(_ *gtk.FlowBox, child *gtk.FlowBoxChild) {
		kanjiDetails.fetchAndDisplay(search.resultsKanji.kanji[child.GetIndex()])
	},
	"searchResultsRowSelected": func() {
		sel := search.results.selected()
		if sel == nil {
			return
		}
		navigation.goTo(sel.ID)
	},
	"searchToggle":      search.toggleOpen,
	"windowButtonPress": adaptButtonHandler(windowButtonMap),
	"windowKeyPress":    adaptKeyHandler(windowKeyMap),
}

var searchEntryKeyMap = keyMap{
	key{gdk.KEY_Escape, 0}: search.deactivateEntry,
}

var windowKeyMap = keyMap{
	key{gdk.KEY_f, gdk.GDK_CONTROL_MASK}: search.activateEntry,
}

var windowButtonMap = buttonMap{
	button{8, 0}: func() {
		search.results.clearSelection()
		navigation.goBack()
	},
	button{9, 0}: func() {
		search.results.clearSelection()
		navigation.goForward()
	},
}

var db *dictdb.DB

// LaunchGUI launches the application user interface, passing the given
// arguments to GTK. It does not return an error; if any errors occur here, the
// program will terminate.
func LaunchGUI(args []string) {
	var err error
	db, err = dictdb.Open()
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
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
	builderData, err := Asset("data/gjisho.glade")
	if err != nil {
		log.Fatalf("Could not load GUI builder data: %v", err)
	}
	builder, err := gtk.BuilderNew()
	if err != nil {
		log.Fatalf("Could not create application builder: %v", err)
	}
	if err := builder.AddFromString(string(builderData)); err != nil {
		log.Fatalf("Could not load data for application builder: %v", err)
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

	logoLoader, _ := gdk.PixbufLoaderNew()
	logoLoader.SetSize(192, 192)
	logoData, err := Asset("data/logo.svg")
	if err != nil {
		log.Fatalf("Could not load logo data: %v", err)
	}
	logoPixbuf, err := logoLoader.WriteAndReturnPixbuf(logoData)
	if err != nil {
		log.Fatalf("Could not process logo data: %v", err)
	}
	aboutDialog.SetLogo(logoPixbuf)

	search.kanjiInput.initRadicals()

	window.Present()

	kanjiIconLoader, _ := gdk.PixbufLoaderNew()
	height := search.entry.GetAllocatedHeight() * 3 / 5
	kanjiIconLoader.SetSize(height, height)
	kanjiIconData, err := Asset("data/kanji-icon.svg")
	if err != nil {
		log.Fatalf("Could not load kanji icon data: %v", err)
	}
	kanjiIconPixbuf, err := kanjiIconLoader.WriteAndReturnPixbuf(kanjiIconData)
	if err != nil {
		log.Fatalf("Could not process kanji icon data: %v", err)
	}
	search.kanjiInput.buttonIcon.SetFromPixbuf(kanjiIconPixbuf)
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
