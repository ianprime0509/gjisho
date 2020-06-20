package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
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

// Search is a wrapper around the search-related components of the app.
type Search struct {
	toggle         *gtk.ToggleButton
	revealer       *gtk.Revealer
	entry          *gtk.SearchEntry
	results        *SearchResultList
	cancelPrevious context.CancelFunc // a function to cancel the previous search
}

// Search searches using the given query.
func (s *Search) Search(query string) {
	ctx := s.startSearch()
	ch := make(chan []jmdict.LookupResult)

	go func() {
		if results, err := dict.Lookup(query); err == nil {
			ch <- results
		} else {
			log.Printf("Lookup query error: %v", err)
		}
		close(ch)
	}()

	go func() {
		select {
		case results := <-ch:
			glib.IdleAdd(func() { s.results.Set(results) })
		case <-ctx.Done():
		}
	}()
}

// Toggle toggles whether the search pane is open.
func (s *Search) Toggle() {
	s.revealer.SetRevealChild(!s.revealer.GetRevealChild())
}

// Activate activates and focuses the search entry.
func (s *Search) Activate() {
	s.toggle.SetActive(true)
	s.revealer.SetRevealChild(true)
	s.entry.GrabFocus()
}

// Deactivate deactivates the search entry.
func (s *Search) Deactivate() {
	s.toggle.SetActive(false)
	s.revealer.SetRevealChild(false)
}

func (s *Search) startSearch() context.Context {
	if s.cancelPrevious != nil {
		s.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelPrevious = cancel
	return ctx
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

// Set sets the currently displayed search results.
func (lst *SearchResultList) Set(results []jmdict.LookupResult) {
	lst.results = results
	removeChildren(&lst.list.Container)
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

	heading, _ := gtk.LabelNew(fmt.Sprintf(`<big>%s</big>`, entry.Heading))
	heading.SetUseMarkup(true)
	heading.SetXAlign(0)
	heading.SetEllipsize(pango.ELLIPSIZE_END)
	box.Add(heading)
	if entry.Heading != entry.PrimaryReading {
		reading, _ := gtk.LabelNew(entry.PrimaryReading)
		reading.SetXAlign(0)
		reading.SetEllipsize(pango.ELLIPSIZE_END)
		box.Add(reading)
	}
	gloss, _ := gtk.LabelNew(fmt.Sprintf(`<small>%s</small>`, entry.GlossSummary))
	gloss.SetUseMarkup(true)
	gloss.SetXAlign(0)
	gloss.SetEllipsize(pango.ELLIPSIZE_END)
	box.Add(gloss)

	return box
}

// EntryNavigation is a wrapper around an EntryDisplay that supports maintaining
// forwards and backwards navigation in a history of entries.
type EntryNavigation struct {
	disp           *EntryDisplay
	backButton     *gtk.Button
	forwardButton  *gtk.Button
	current        int
	backStack      []int
	forwardStack   []int
	cancelPrevious context.CancelFunc // a function to cancel the previous navigation operation
}

// FollowLink attempts to follow the given link and returns whether it was able
// to do so.
func (n *EntryNavigation) FollowLink(link *url.URL) bool {
	if link.Scheme != "entry" {
		return false
	}

	if match, err := dict.LookupByRef(link.Host); err == nil {
		// If I try to call GoTo directly, then for some reason the program
		// crashes (probably because the link text gets freed or otherwise
		// corrupted after navigation)
		glib.IdleAdd(func() { n.GoTo(match.ID) })
	} else {
		log.Printf("Error fetching entry for link %v: %v", link, err)
	}
	return true
}

// GoTo navigates to the entry with the given ID.
func (n *EntryNavigation) GoTo(id int) {
	ctx := n.startNavigation()

	if n.current != 0 {
		n.backStack = append(n.backStack, n.current)
	}
	n.current = id
	n.forwardStack = nil
	n.updateSensitivity()

	n.disp.FetchAndDisplay(ctx, id)
}

// GoBack navigates to the previous entry.
func (n *EntryNavigation) GoBack() {
	if len(n.backStack) == 0 {
		return
	}

	ctx := n.startNavigation()

	if n.current != 0 {
		n.forwardStack = append(n.forwardStack, n.current)
	}
	n.current = n.backStack[len(n.backStack)-1]
	n.backStack = n.backStack[:len(n.backStack)-1]
	n.updateSensitivity()

	n.disp.FetchAndDisplay(ctx, n.current)
}

// GoForward navigates to the next entry.
func (n *EntryNavigation) GoForward() {
	if len(n.forwardStack) == 0 {
		return
	}

	ctx := n.startNavigation()

	if n.current != 0 {
		n.backStack = append(n.backStack, n.current)
	}
	n.current = n.forwardStack[len(n.forwardStack)-1]
	n.forwardStack = n.forwardStack[:len(n.forwardStack)-1]
	n.updateSensitivity()

	n.disp.FetchAndDisplay(ctx, n.current)
}

// startNavigation cancels any previous navigation in progress and returns a
// context for a new one.
func (n *EntryNavigation) startNavigation() context.Context {
	if n.cancelPrevious != nil {
		n.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	n.cancelPrevious = cancel
	return ctx
}

func (n *EntryNavigation) updateSensitivity() {
	n.backButton.SetSensitive(len(n.backStack) > 0)
	n.forwardButton.SetSensitive(len(n.forwardStack) > 0)
}

// EntryDisplay is the main display area for a dictionary entry.
type EntryDisplay struct {
	kanjiList          *KanjiList
	exampleList        *ExampleList
	primaryKanaLabel   *gtk.Label
	primaryKanjiLabel  *gtk.Label
	detailsLabel       *gtk.Label
	kanjiWritingsLabel *gtk.Label
	kanaWritingsLabel  *gtk.Label
	cancelDisplay      context.CancelFunc
}

// FetchAndDisplay fetches and displays the dictionary entry with the given ID
// in the display area.
func (disp *EntryDisplay) FetchAndDisplay(ctx context.Context, id int) {
	ch := make(chan jmdict.Entry)
	go func() {
		if entry, err := dict.Fetch(id); err == nil {
			ch <- entry
		} else {
			log.Printf("Error fetching entry with ID %v: %v", id, err)
		}
		close(ch)
	}()

	go func() {
		select {
		case entry := <-ch:
			glib.IdleAdd(func() { disp.display(entry) })

			disp.kanjiList.FetchAndDisplay(ctx, entry.AssociatedKanji())
			disp.exampleList.FetchAndDisplay(ctx, entry.Heading())
		case <-ctx.Done():
		}
	}()
}

func (disp *EntryDisplay) display(entry jmdict.Entry) {
	disp.primaryKanjiLabel.SetText(entry.Heading())
	if entry.Heading() != entry.PrimaryReading() {
		disp.primaryKanaLabel.SetText(entry.PrimaryReading())
		disp.primaryKanaLabel.Show()
	} else {
		disp.primaryKanaLabel.SetText("")
		disp.primaryKanaLabel.Hide()
	}
	disp.detailsLabel.SetMarkup(fmtSenses(entry.Senses))
	disp.kanjiWritingsLabel.SetMarkup(fmtKanjiWritings(entry.KanjiWritings))
	disp.kanaWritingsLabel.SetMarkup(fmtKanaReadings(entry.KanaWritings))
}

func fmtKanjiWritings(kanji []jmdict.KanjiWriting) string {
	if len(kanji) == 0 {
		return "<i>None</i>"
	}

	var forms []string
	for _, w := range kanji {
		sb := new(strings.Builder)
		sb.WriteString(w.Writing)
		info := strings.Join(w.Info, ", ")
		if info != "" {
			fmt.Fprintf(sb, " <i>%v</i>", info)
		}
		forms = append(forms, sb.String())
	}
	return strings.Join(forms, "\n")
}

func fmtKanaReadings(kana []jmdict.KanaWriting) string {
	var forms []string
	for _, w := range kana {
		sb := new(strings.Builder)
		sb.WriteString(w.Writing)
		var details []string
		info := strings.Join(w.Info, ", ")
		if info != "" {
			details = append(details, info)
		}
		restr := strings.Join(w.Restrictions, ", ")
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

// FetchAndDisplay fetches and displays information about the given kanji in the
// list.
func (lst *KanjiList) FetchAndDisplay(ctx context.Context, kanji []string) {
	ch := make(chan []kanjidic.Character)
	go func() {
		var results []kanjidic.Character
		for _, c := range kanji {
			if res, err := kanjiDict.Fetch(c); err == nil {
				results = append(results, res)
			} else {
				log.Printf("Error fetching kanji information for %q: %v", c, err)
			}
		}
		ch <- results
		close(ch)
	}()

	go func() {
		select {
		case kanji := <-ch:
			glib.IdleAdd(func() { lst.display(kanji) })
		case <-ctx.Done():
		}
	}()
}

func (lst *KanjiList) display(kanji []kanjidic.Character) {
	removeChildren(&lst.list.Container)
	lst.kanji = kanji
	for _, result := range lst.kanji {
		lst.list.Add(newKanjiListRow(result))
	}
	lst.list.ShowAll()
}

func newKanjiListRow(c kanjidic.Character) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	rowBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)

	header, _ := gtk.LabelNew(fmt.Sprintf(`<span size="xx-large">%s</span>`, c.Literal))
	header.SetUseMarkup(true)
	rowBox.PackStart(header, false, false, 5)

	details, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	on, _ := gtk.LabelNew(strings.Join(c.Readings(kanjidic.OnReading), ", "))
	on.SetLineWrap(true)
	on.SetJustify(gtk.JUSTIFY_CENTER)
	details.Add(on)
	kun, _ := gtk.LabelNew(strings.Join(c.Readings(kanjidic.KunReading), ", "))
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

// ExampleList is a list of examples associated with an entry.
type ExampleList struct {
	list       *gtk.ListBox
	examples   []tatoeba.Example
	nDisplayed int
}

// FetchAndDisplay fetches and displays examples for the given word in the list.
func (lst *ExampleList) FetchAndDisplay(ctx context.Context, word string) {
	ch := make(chan []tatoeba.Example)
	go func() {
		if examples, err := exampleDict.FetchByWord(word); err == nil {
			ch <- examples
		} else {
			log.Printf("Error fetching examples for %q: %v", word, err)
		}
		close(ch)
	}()

	go func() {
		select {
		case examples := <-ch:
			glib.IdleAdd(func() { lst.display(examples) })
		case <-ctx.Done():
		}
	}()
}

// ShowMore displays more examples in the list.
func (lst *ExampleList) ShowMore() {
	maxIndex := lst.nDisplayed + 20
	for ; lst.nDisplayed < maxIndex && lst.nDisplayed < len(lst.examples); lst.nDisplayed++ {
		lst.list.Add(newExampleListRow(lst.examples[lst.nDisplayed]))
	}
	lst.list.ShowAll()
}

func (lst *ExampleList) display(examples []tatoeba.Example) {
	removeChildren(&lst.list.Container)
	lst.nDisplayed = 0
	lst.examples = examples
	lst.ShowMore()
}

func newExampleListRow(ex tatoeba.Example) *gtk.ListBoxRow {
	row, _ := gtk.ListBoxRowNew()
	rowBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	jpnLabel, _ := gtk.LabelNew(ex.Japanese)
	jpnLabel.SetLineWrap(true)
	jpnLabel.SetXAlign(0)
	rowBox.Add(jpnLabel)
	engLabel, _ := gtk.LabelNew(ex.English)
	engLabel.SetLineWrap(true)
	engLabel.SetXAlign(0)
	rowBox.Add(engLabel)

	row.Add(rowBox)
	return row
}

// KanjiDetails is a modal window showing additional details about a kanji.
type KanjiDetails struct {
	window          *gtk.Window
	charLabel       *gtk.Label
	subtitleLabel   *gtk.Label
	readingMeanings *gtk.Box
	dictRefsLabel   *gtk.Label
	queryCodesLabel *gtk.Label
}

// Display displays the given kanji in the window (but does not immediately show
// the window).
func (kd *KanjiDetails) Display(c kanjidic.Character) {
	kd.charLabel.SetText(c.Literal)
	kd.subtitleLabel.SetMarkup(fmtSubtitle(c))
	removeChildren(&kd.readingMeanings.Container)
	for _, rm := range c.ReadingMeaningGroups {
		kd.readingMeanings.Add(newReadingMeaningLabel(rm))
	}
	kd.readingMeanings.ShowAll()
	kd.dictRefsLabel.SetMarkup(fmtDictRefs(c.DictRefs))
	kd.queryCodesLabel.SetMarkup(fmtQueryCodes(c.QueryCodes))
}

// Present presents the kanji details window.
func (kd *KanjiDetails) Present() {
	kd.window.Present()
}

func fmtSubtitle(c kanjidic.Character) string {
	fmtStrokes := func(count int) string {
		if count == 1 {
			return "1 stroke"
		}
		return fmt.Sprintf("%v strokes", count)
	}

	strokes := fmtStrokes(c.Misc.StrokeCounts[0])
	extraStrokes := make([]string, 0, len(c.Misc.StrokeCounts)-1)
	for _, count := range c.Misc.StrokeCounts[1:] {
		extraStrokes = append(extraStrokes, fmtStrokes(count))
	}

	sb := new(strings.Builder)
	sb.WriteString(strokes)
	if len(extraStrokes) > 0 {
		fmt.Fprintf(sb, " (alternate %v)", strings.Join(extraStrokes, ", "))
	}
	sb.WriteRune('\n')
	sb.WriteString(c.Misc.Grade.String())
	return sb.String()
}

func newReadingMeaningLabel(rm kanjidic.ReadingMeaningGroup) *gtk.Label {
	var kun, on []string
	for _, r := range rm.Readings {
		switch r.Type {
		case kanjidic.KunReading:
			if r.Jouyou {
				kun = append(kun, fmt.Sprintf("<b>%v</b>", r.Reading))
			} else {
				kun = append(kun, r.Reading)
			}
		case kanjidic.OnReading:
			s := r.Reading
			if r.OnType != "" {
				s += fmt.Sprintf("(%v)", r.OnType)
			}
			if r.Jouyou {
				on = append(on, fmt.Sprintf("<b>%v</b>", s))
			} else {
				on = append(on, s)
			}
		}
	}

	sb := new(strings.Builder)
	if len(kun) > 0 {
		fmt.Fprintf(sb, "<b>Kun:</b> %v\n", strings.Join(kun, ", "))
	}
	if len(on) > 0 {
		fmt.Fprintf(sb, "<b>On:</b> %v\n", strings.Join(on, ", "))
	}
	if sb.Len() > 0 && len(rm.Meanings) > 0 {
		sb.WriteRune('\n')
	}

	item := 1
	for _, m := range rm.Meanings {
		if m.Language != "" {
			continue
		}
		fmt.Fprintf(sb, "%v. %v\n", item, m.Meaning)
		item++
	}

	lbl, _ := gtk.LabelNew(sb.String())
	lbl.SetUseMarkup(true)
	lbl.SetXAlign(0)
	lbl.SetLineWrap(true)
	return lbl
}

func fmtDictRefs(refs []kanjidic.DictRef) string {
	sb := new(strings.Builder)
	for _, ref := range refs {
		fmt.Fprintf(sb, "<b>%v</b> %v\n", ref.Type, ref.Index)
	}
	return sb.String()
}

func fmtQueryCodes(codes []kanjidic.QueryCode) string {
	sb := new(strings.Builder)
	for _, code := range codes {
		fmt.Fprintf(sb, "<b>%v</b> %v\n", code.Type, code.Code)
	}
	return sb.String()
}
