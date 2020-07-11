package main

import (
	"log"
	"net/url"
	"os"
	"reflect"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/kanjivg"
	"github.com/ianprime0509/gjisho/kradfile"
	"github.com/ianprime0509/gjisho/tatoeba"
)

const appID = "com.github.ianprime0509.gjisho"

var aboutDialog *gtk.AboutDialog
var moreInfoRevealer *gtk.Revealer

var searchResults = new(SearchResultList)
var searchResultsKanji = new(SearchResultsKanji)
var kanjiInput = new(KanjiInput)
var search = &Search{results: searchResults, resultsKanji: searchResultsKanji, kanjiInput: kanjiInput}
var kanjiList = new(KanjiList)
var kanjiDetails = new(KanjiDetails)
var exampleList = new(ExampleList)
var exampleDetails = new(ExampleDetails)
var entryDisplay = &EntryDisplay{kanjiList: kanjiList, exampleList: exampleList}
var navigation = &EntryNavigation{disp: entryDisplay}

var appComponents = map[string]interface{}{
	"aboutDialog":                      &aboutDialog,
	"backButton":                       &navigation.backButton,
	"entryDetailsLabel":                &entryDisplay.detailsLabel,
	"entryScrolledWindow":              &entryDisplay.scrolledWindow,
	"exampleDetailsEnglishLabel":       &exampleDetails.englishLabel,
	"exampleDetailsJapaneseLabel":      &exampleDetails.japaneseLabel,
	"exampleDetailsScrolledWindow":     &exampleDetails.scrolledWindow,
	"exampleDetailsWindow":             &exampleDetails.window,
	"exampleDetailsWordsList":          &exampleDetails.wordsList,
	"examplesList":                     &exampleList.list,
	"examplesScrolledWindow":           &exampleList.scrolledWindow,
	"forwardButton":                    &navigation.forwardButton,
	"kanjiDetailsCharacterLabel":       &kanjiDetails.charLabel,
	"kanjiDetailsDictRefsLabel":        &kanjiDetails.dictRefsLabel,
	"kanjiDetailsReadingMeanings":      &kanjiDetails.readingMeanings,
	"kanjiDetailsScrolledWindow":       &kanjiDetails.scrolledWindow,
	"kanjiDetailsStrokeOrder":          &kanjiDetails.strokeOrder,
	"kanjiDetailsSubtitleLabel":        &kanjiDetails.subtitleLabel,
	"kanjiDetailsQueryCodesLabel":      &kanjiDetails.queryCodesLabel,
	"kanjiDetailsWindow":               &kanjiDetails.window,
	"kanjiInputButton":                 &kanjiInput.button,
	"kanjiInputButtonIcon":             &kanjiInput.buttonIcon,
	"kanjiInputPopover":                &kanjiInput.popover,
	"kanjiInputRadicals":               &kanjiInput.radicalsBox,
	"kanjiInputRadicalsScrolledWindow": &kanjiInput.radicalsScrolledWindow,
	"kanjiInputResults":                &kanjiInput.resultsBox,
	"kanjiInputResultsScrolledWindow":  &kanjiInput.resultsScrolledWindow,
	"kanaWritingsLabel":                &entryDisplay.kanaWritingsLabel,
	"kanjiList":                        &kanjiList.list,
	"kanjiScrolledWindow":              &kanjiList.scrolledWindow,
	"kanjiWritingsLabel":               &entryDisplay.kanjiWritingsLabel,
	"moreInfoRevealer":                 &moreInfoRevealer,
	"primaryKanaLabel":                 &entryDisplay.primaryKanaLabel,
	"primaryKanjiLabel":                &entryDisplay.primaryKanjiLabel,
	"searchEntry":                      &search.entry,
	"searchRevealer":                   &search.revealer,
	"searchResults":                    &searchResults.list,
	"searchResultsKanji":               &searchResultsKanji.box,
	"searchResultsScrolledWindow":      &searchResults.scrolledWindow,
	"searchToggleButton":               &search.toggle,
	"strokeOrderScrolledWindow":        &kanjiDetails.strokeOrderScrolledWindow,
	"writingsScrolledWindow":           &entryDisplay.writingsScrolledWindow,
}

var signals = map[string]interface{}{
	"activateLink": func(_ *gtk.Label, uri string) bool {
		url, err := url.Parse(uri)
		if err != nil {
			log.Printf("Invalid URL: %v", uri)
			return true
		}
		search.results.ClearSelection()
		return navigation.FollowLink(url)
	},
	"exampleDetailsWordActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		search.results.ClearSelection()
		navigation.GoTo(exampleDetails.words[row.GetIndex()].ID)
		exampleDetails.Close()
	},
	"examplesEdgeReached": func(_ *gtk.ScrolledWindow, pos gtk.PositionType) {
		if pos == gtk.POS_BOTTOM {
			exampleList.ShowMore()
		}
	},
	"exampleListRowActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		exampleDetails.FetchAndDisplay(exampleList.examples[row.GetIndex()])
	},
	"hideWidget":  func(w interface{ Hide() }) { w.Hide() },
	"inhibitNext": func() bool { return true },
	"kanjiListRowActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		kanjiDetails.FetchAndDisplay(kanjiList.kanji[row.GetIndex()])
	},
	"kanjiInputButtonToggled": func(b *gtk.ToggleButton) {
		if b.GetActive() {
			kanjiInput.Display()
		}
	},
	"kanjiInputPopoverClosed": func() {
		kanjiInput.button.SetActive(false)
	},
	"kanjiInputResultsChildActivated": func(_ *gtk.FlowBox, child *gtk.FlowBoxChild) {
		search.InsertEntryText(kanjiInput.results[child.GetIndex()].Literal)
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
	"searchEntryKeyPress": AdaptKeyHandler(searchEntryKeyMap),
	"searchResultsEdgeReached": func(_ *gtk.ScrolledWindow, pos gtk.PositionType) {
		if pos == gtk.POS_BOTTOM {
			searchResults.ShowMore()
		}
	},
	"searchResultsKanjiChildActivated": func(_ *gtk.FlowBox, child *gtk.FlowBoxChild) {
		kanjiDetails.FetchAndDisplay(searchResultsKanji.kanji[child.GetIndex()])
	},
	"searchResultsRowSelected": func() {
		sel := searchResults.Selected()
		if sel == nil {
			return
		}
		navigation.GoTo(sel.ID)
	},
	"searchToggle":      search.Toggle,
	"windowButtonPress": AdaptButtonHandler(windowButtonMap),
	"windowKeyPress":    AdaptKeyHandler(windowKeyMap),
}

var searchEntryKeyMap = KeyMap{
	Key{gdk.KEY_Escape, 0}: search.Deactivate,
}

var windowKeyMap = KeyMap{
	Key{gdk.KEY_f, gdk.GDK_CONTROL_MASK}: search.Activate,
}

var windowButtonMap = ButtonMap{
	Button{8, 0}: func() {
		search.results.ClearSelection()
		navigation.GoBack()
	},
	Button{9, 0}: func() {
		search.results.ClearSelection()
		navigation.GoForward()
	},
}

var dict *jmdict.JMdict
var kanjiDict *kanjidic.KANJIDIC
var radicalDict *kradfile.KRADFILE
var exampleDict *tatoeba.Tatoeba
var strokeDict *kanjivg.KanjiVG

// LaunchGUI launches the application user interface, passing the given
// arguments to GTK. It does not return an error; if any errors occur here, the
// program will terminate.
func LaunchGUI(args []string) {
	db, err := OpenDB()
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}

	dict, err = jmdict.New(db)
	if err != nil {
		log.Fatalf("Could not open JMdict handler: %v", err)
	}

	kanjiDict, err = kanjidic.New(db)
	if err != nil {
		log.Fatalf("Could not open KANJIDIC handler: %v", err)
	}

	radicalDict, err = kradfile.New(db)
	if err != nil {
		log.Fatalf("Could not open KRADFILE handler: %v", err)
	}

	exampleDict, err = tatoeba.New(db)
	if err != nil {
		log.Fatalf("Could not open Tatoeba handler: %v", err)
	}

	strokeDict, err = kanjivg.New(db)
	if err != nil {
		log.Fatalf("Could not open KanjiVG handler: %v", err)
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

	kanjiIconData, err := Asset("data/kanji-icon.png")
	if err != nil {
		log.Fatalf("Could not load kanji icon data: %v", err)
	}
	pbLoader, _ := gdk.PixbufLoaderNew()
	kanjiIcon, err := pbLoader.WriteAndReturnPixbuf(kanjiIconData)
	if err != nil {
		log.Fatalf("Could not process kanji icon data: %v", err)
	}

	search.kanjiInput.InitRadicals()

	window.Show()

	height := search.entry.GetAllocatedHeight() * 3 / 5
	kanjiIcon, _ = kanjiIcon.ScaleSimple(height, height, gdk.INTERP_BILINEAR)
	kanjiInput.buttonIcon.SetFromPixbuf(kanjiIcon)
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
