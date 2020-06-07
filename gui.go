package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync/atomic"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
)

const appID = "com.github.ianprime0509.gjisho"

type hideable interface {
	Hide()
}

var aboutDialog *gtk.AboutDialog
var moreInfoRevealer *gtk.Revealer
var searchEntry *gtk.SearchEntry
var searchRevealer *gtk.Revealer
var searchToggleButton *gtk.ToggleButton

var searchResults = new(SearchResultList)
var entryDisplay = new(EntryDisplay)
var kanjiList = new(KanjiList)

var appComponents = map[string]interface{}{
	"aboutDialog":        &aboutDialog,
	"entryDetailsLabel":  &entryDisplay.detailsLabel,
	"kanaWritingsLabel":  &entryDisplay.kanaWritingsLabel,
	"kanjiList":          &kanjiList.list,
	"kanjiWritingsLabel": &entryDisplay.kanjiWritingsLabel,
	"moreInfoRevealer":   &moreInfoRevealer,
	"primaryKanaLabel":   &entryDisplay.primaryKanaLabel,
	"primaryKanjiLabel":  &entryDisplay.primaryKanjiLabel,
	"searchEntry":        &searchEntry,
	"searchRevealer":     &searchRevealer,
	"searchResults":      &searchResults.list,
	"searchToggleButton": &searchToggleButton,
}

var signals = map[string]interface{}{
	"hideWidget":  hideable.Hide,
	"inhibitNext": func() bool { return true },
	"kanjiListRowActivated": func(_ *gtk.ListBox, row *gtk.ListBoxRow) {
		log.Printf("Kanji row activated: %v", kanjiList.kanji[row.GetIndex()])
	},
	"moreInfoToggle": func() {
		moreInfoRevealer.SetRevealChild(!moreInfoRevealer.GetRevealChild())
	},
	"searchChanged": searchChanged,
	"searchEntryKeyPress": func(_ interface{}, ev *gdk.Event) {
		keyEv := &gdk.EventKey{Event: ev}
		if keyEv.KeyVal() == gdk.KEY_Escape {
			stopSearch()
		}
	},
	"searchResultsEndReached": func() { searchResults.ShowMore() },
	"searchResultsRowSelected": func() {
		sel := searchResults.Selected()
		if sel == nil {
			return
		}
		if entry, err := dict.Fetch(sel.ID); err == nil {
			entryDisplay.Display(entry)
			kanjiList.Display(entry.AssociatedKanji())
		} else {
			log.Printf("Could not fetch entry with ID %v: %v", searchResults.Selected().ID, err)
		}
	},
	"searchToggle": func() {
		searchRevealer.SetRevealChild(!searchRevealer.GetRevealChild())
	},
	"windowKeyPress": func(_ interface{}, ev *gdk.Event) {
		keyEv := &gdk.EventKey{Event: ev}
		if keyEv.KeyVal() == gdk.KEY_f && keyEv.State()&gdk.GDK_CONTROL_MASK != 0 {
			startSearch()
		}
	},
}

var dict *jmdict.JMdict
var kanjiDict *kanjidic.Kanjidic

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

	css, err := gtk.CssProviderNew()
	if err != nil {
		log.Fatalf("Could not create CSS provider: %v", err)
	}
	if err := css.LoadFromPath("gui.css"); err != nil {
		log.Fatalf("Could not load CSS: %v", err)
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
			log.Fatalf("Could not get application component %v: %v", name, err)
		}
		reflect.ValueOf(ptr).Elem().Set(reflect.ValueOf(comp))
	}
}

// SearchResultList is a list of search results displayed in the GUI.
type SearchResultList struct {
	list       *gtk.ListBox
	results    []jmdict.LookupResult
	nDisplayed int
}

// Selected returns the currently selected search result, or nil if none is
// selected.
func (lst *SearchResultList) Selected() *jmdict.LookupResult {
	if row := lst.list.GetSelectedRow(); row != nil {
		return &lst.results[row.GetIndex()]
	}
	return nil
}

// SetResults sets the currently displayed search results.
func (lst *SearchResultList) SetResults(results []jmdict.LookupResult) {
	lst.results = results
	lst.list.GetChildren().Foreach(func(e interface{}) {
		lst.list.Remove(e.(gtk.IWidget))
	})
	lst.nDisplayed = 0
	lst.ShowMore()
}

// ShowMore displays more search results in the list.
func (lst *SearchResultList) ShowMore() {
	maxIndex := lst.nDisplayed + 50
	for ; lst.nDisplayed < len(lst.results) && lst.nDisplayed < maxIndex; lst.nDisplayed++ {
		lst.list.Add(newSearchResult(lst.results[lst.nDisplayed]))
	}
	lst.list.ShowAll()
}

// newSearchResult creates a search result widget for display.
func newSearchResult(entry jmdict.LookupResult) gtk.IWidget {
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

// Consider the situation where a search query comes in that takes a long time,
// followed by one that takes less time. To avoid having the results of the
// first query come in later and overwrite the results of the second, we use a
// counter to identify the queries in sequence and prevent overwriting newer
// query results with older ones.
var queryCounter uint32
var currentQuery uint32

func searchChanged(entry *gtk.SearchEntry) {
	query, _ := entry.GetText()
	go func() {
		queryNum := atomic.AddUint32(&queryCounter, 1)
		results, err := dict.Lookup(query)
		if err == nil {
			glib.IdleAdd(func() {
				// There is no race condition here since this function will only be
				// executed on the main thread
				if queryNum > currentQuery {
					searchResults.SetResults(results)
					currentQuery = queryNum
				}
			})
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

// EntryDisplay is the main display area for a dictionary entry.
type EntryDisplay struct {
	primaryKanaLabel   *gtk.Label
	primaryKanjiLabel  *gtk.Label
	detailsLabel       *gtk.Label
	kanjiWritingsLabel *gtk.Label
	kanaWritingsLabel  *gtk.Label
}

// Display displays the given dictionary entry in the display area.
func (disp *EntryDisplay) Display(entry jmdict.Entry) {
	disp.primaryKanjiLabel.SetText(entry.Heading())
	if entry.Heading() != entry.PrimaryReading() {
		disp.primaryKanaLabel.SetText(entry.PrimaryReading())
		disp.primaryKanaLabel.Show()
	} else {
		disp.primaryKanaLabel.SetText("")
		disp.primaryKanaLabel.Hide()
	}
	disp.detailsLabel.SetMarkup(fmtSenses(entry.Senses))
	disp.kanjiWritingsLabel.SetMarkup(fmtKanjiWritings(entry.KanjiReadings))
	disp.kanaWritingsLabel.SetMarkup(fmtKanaWritings(entry.KanaReadings))
}

func fmtKanjiWritings(kanji []jmdict.KanjiReading) string {
	if len(kanji) == 0 {
		return "<i>None</i>"
	}

	var forms []string
	for _, reading := range kanji {
		sb := new(strings.Builder)
		sb.WriteString(reading.Reading)
		info := strings.Join(reading.Info, ", ")
		if info != "" {
			fmt.Fprintf(sb, " <i>%v</i>", info)
		}
		forms = append(forms, sb.String())
	}
	return strings.Join(forms, "\n")
}

func fmtKanaWritings(kana []jmdict.KanaReading) string {
	var forms []string
	for _, reading := range kana {
		sb := new(strings.Builder)
		sb.WriteString(reading.Reading)
		var details []string
		info := strings.Join(reading.Info, ", ")
		if info != "" {
			details = append(details, info)
		}
		restr := strings.Join(reading.Restrictions, ", ")
		if restr != "" {
			details = append(details, "restricted to "+restr)
		}
		if len(details) > 0 {
			fmt.Fprintf(sb, " <i>%v</i>", strings.Join(details, "; "))
		}
		forms = append(forms, sb.String())
	}
	return strings.Join(forms, "\n")
}

func fmtSenses(senses []jmdict.Sense) string {
	sb := new(strings.Builder)
	glossIdx := 1
	for _, sense := range senses {
		for _, pos := range sense.PartsOfSpeech {
			fmt.Fprintf(sb, "<i>%v</i>\n", pos)
		}
		for _, info := range sense.Info {
			fmt.Fprintf(sb, "<i>%v</i>\n", info)
		}
		for _, field := range sense.Fields {
			fmt.Fprintf(sb, "<i>%v</i>\n", field)
		}
		for _, misc := range sense.Misc {
			fmt.Fprintf(sb, "<i>%v</i>\n", misc)
		}
		if len(sense.KanjiRestrictions) > 0 {
			fmt.Fprintf(sb, "<i>Restricted to %v</i>\n", strings.Join(sense.KanjiRestrictions, "\n"))
		}
		if len(sense.KanaRestrictions) > 0 {
			fmt.Fprintf(sb, "<i>Restricted to %v</i>\n", strings.Join(sense.KanaRestrictions, ", "))
		}
		if len(sense.Dialects) > 0 {
			fmt.Fprintf(sb, "<i>Dialects: %v</i>\n", strings.Join(sense.Dialects, ", "))
		}
		if len(sense.LoanSources) > 0 {
			var sources []string
			for _, source := range sense.LoanSources {
				text := source.Source
				if source.Wasei {
					text += " (wasei)"
				}
				sources = append(sources, text)
			}
			fmt.Fprintf(sb, "<i>Loanword sources: %v</i>\n", strings.Join(sources, ", "))
		}
		if len(sense.CrossReferences) > 0 {
			var refs []string
			for _, ref := range sense.CrossReferences {
				refs = append(refs, fmtEntryRef(ref))
			}
			fmt.Fprintf(sb, "<i>See also: %v</i>\n", strings.Join(refs, ", "))
		}
		if len(sense.Antonyms) > 0 {
			var ants []string
			for _, ant := range sense.Antonyms {
				ants = append(ants, fmtEntryRef(ant))
			}
			fmt.Fprintf(sb, "<i>Antonyms: %v</i>\n", strings.Join(ants, ", "))
		}

		// We only want to print an extra newline if there were glosses formatted
		// for this sense, so we keep track of that
		foundGloss := false
		for _, gloss := range sense.Glosses {
			// Only consider English glosses for now
			if gloss.Language != "" {
				continue
			}

			fmt.Fprintf(sb, "%v. %v\n", glossIdx, gloss.Gloss)
			glossIdx++
			foundGloss = true
		}
		if foundGloss {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func fmtEntryRef(entry string) string {
	return fmt.Sprintf("<a href=\"entry://%s\">%[1]s</a>", entry)
}

// KanjiList is an overview list of kanji associated with an entry.
type KanjiList struct {
	list  *gtk.ListBox
	kanji []kanjidic.Character
}

// Display displays information about the given kanji in the list.
func (lst *KanjiList) Display(kanji []string) {
	lst.list.GetChildren().Foreach(func(e interface{}) { lst.list.Remove(e.(gtk.IWidget)) })
	lst.kanji = nil

	for _, c := range kanji {
		if result, err := kanjiDict.Fetch(c); err == nil {
			lst.list.Add(newKanjiListRow(result))
			lst.kanji = append(lst.kanji, result)
		} else {
			log.Printf("Error fetching kanji information for %q: %v", c, err)
		}
	}

	lst.list.ShowAll()
}

func newKanjiListRow(c kanjidic.Character) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	rowBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)

	header, _ := gtk.LabelNew(c.Literal)
	ctx, _ := header.GetStyleContext()
	ctx.AddClass("kanji-header")
	rowBox.PackStart(header, false, false, 0)

	details, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	on, _ := gtk.LabelNew(strings.Join(c.Readings(kanjidic.On), ", "))
	on.SetLineWrap(true)
	on.SetJustify(gtk.JUSTIFY_CENTER)
	details.Add(on)
	kun, _ := gtk.LabelNew(strings.Join(c.Readings(kanjidic.Kun), ", "))
	kun.SetLineWrap(true)
	kun.SetJustify(gtk.JUSTIFY_CENTER)
	details.Add(kun)
	meanings, _ := gtk.LabelNew(strings.Join(c.Meanings(), ", "))
	meanings.SetLineWrap(true)
	meanings.SetJustify(gtk.JUSTIFY_CENTER)
	details.Add(meanings)
	rowBox.PackStart(details, true, true, 0)

	row.Add(rowBox)
	return row
}
