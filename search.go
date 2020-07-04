package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"unicode/utf8"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/ianprime0509/gjisho/internal/util"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kradfile"
)

// Search is a wrapper around the search-related components of the app.
type Search struct {
	toggle         *gtk.ToggleButton
	revealer       *gtk.Revealer
	entry          *gtk.SearchEntry
	kanjiInput     *KanjiInput
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
	// If we grab focus immediately, then if this is the first time the search
	// entry is being opened, we'll get "gtk_widget_event: assertion
	// 'WIDGET_REALIZED_FOR_EVENT (widget, event)' failed", which seems to imply
	// that the focus is being grabbed before the widget actually exists (even
	// though it works fine anyways). To get rid of the message, we just wait
	// until the next opportunity to grab the focus.
	glib.IdleAdd(s.entry.GrabFocus)
}

// Deactivate deactivates the search entry.
func (s *Search) Deactivate() {
	s.toggle.SetActive(false)
	s.revealer.SetRevealChild(false)
}

// InsertEntryText inserts text in the search entry buffer at the current
// position.
func (s *Search) InsertEntryText(text string) {
	old, _ := search.entry.GetText()
	rs := []rune(old)
	pos := search.entry.GetPosition()
	search.entry.SetText(string(rs[:pos]) + text + string(rs[pos:]))
	search.entry.SetPosition(pos + utf8.RuneCountInString(text))
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
	list           *gtk.ListBox
	scrolledWindow *gtk.ScrolledWindow
	results        []jmdict.LookupResult
	nDisplayed     int
}

// ClearSelection clears the currently selected result.
func (lst *SearchResultList) ClearSelection() {
	lst.list.SelectRow(nil)
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
	util.RemoveChildren(&lst.list.Container)
	lst.nDisplayed = 0
	lst.ShowMore()
	util.ScrollToStart(lst.scrolledWindow)
}

// ShowMore displays more search results in the list.
func (lst *SearchResultList) ShowMore() {
	maxIndex := lst.nDisplayed + 50
	for ; lst.nDisplayed < len(lst.results) && lst.nDisplayed < maxIndex; lst.nDisplayed++ {
		lst.list.Add(NewSearchResult(lst.results[lst.nDisplayed]))
	}
	lst.list.ShowAll()
}

// NewSearchResult creates a search result widget for display.
func NewSearchResult(entry jmdict.LookupResult) *gtk.Box {
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)

	headingBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	heading, _ := gtk.LabelNew(fmt.Sprintf(`<big>%s</big>`, entry.Heading))
	heading.SetUseMarkup(true)
	heading.SetXAlign(0)
	heading.SetEllipsize(pango.ELLIPSIZE_END)
	headingBox.Add(heading)
	if entry.Priority > 0 {
		commonLabel, _ := gtk.LabelNew(" <small><i>common</i></small>")
		commonLabel.SetUseMarkup(true)
		headingBox.PackEnd(commonLabel, false, false, 0)
	}
	box.Add(headingBox)
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

// KanjiInput is a special input popover to make it easier to input kanji.
type KanjiInput struct {
	button                 *gtk.ToggleButton
	buttonIcon             *gtk.Image
	popover                *gtk.Popover
	radicalsScrolledWindow *gtk.ScrolledWindow
	radicals               *gtk.Box
	radicalButtons         map[string]*gtk.FlowBoxChild
	selectedRadicals       map[string]struct{}
	resultsScrolledWindow  *gtk.ScrolledWindow
	results                *gtk.FlowBox
	resultKanji            []kradfile.Kanji
}

// InitRadicals initializes the radical input buttons.
func (ki *KanjiInput) InitRadicals() {
	ki.radicalButtons = make(map[string]*gtk.FlowBoxChild, len(kradfile.RadicalStrokes))
	ki.selectedRadicals = make(map[string]struct{})

	rads := make(map[int][]string)
	for rad, strokes := range kradfile.RadicalStrokes {
		rads[strokes] = append(rads[strokes], rad)
	}
	for _, lst := range rads {
		sort.Strings(lst)
	}

	strokes := make([]int, 0, len(rads))
	for s := range rads {
		strokes = append(strokes, s)
	}
	sort.Ints(strokes)

	for _, s := range strokes {
		var heading string
		if s == 1 {
			heading = "1 stroke"
		} else {
			heading = fmt.Sprintf("%v strokes", s)
		}
		lbl, _ := gtk.LabelNew(heading)
		ki.radicals.Add(lbl)

		fb, _ := gtk.FlowBoxNew()
		fb.SetMinChildrenPerLine(10)
		fb.SetMaxChildrenPerLine(10)
		fb.SetSelectionMode(gtk.SELECTION_NONE)
		for _, rad := range rads[s] {
			rad := rad
			lbl, _ := gtk.LabelNew(kanjiLabelMarkup(rad, false))
			lbl.SetUseMarkup(true)
			b, _ := gtk.FlowBoxChildNew()
			b.Add(lbl)
			ki.radicalButtons[rad] = b
			fb.Add(b)
			// Why can't I just use the activate signal on b itself? GTK is so
			// stupid sometimes.
			fb.Connect("child-activated", func(_ interface{}, child *gtk.FlowBoxChild) {
				if child.GetIndex() == b.GetIndex() && b.GetSensitive() {
					ki.toggleRadical(rad)
					ki.updateResults()
				}
			})
		}
		ki.radicals.Add(fb)
	}
}

// Display displays the kanji input popover.
func (ki *KanjiInput) Display() {
	util.ScrollToStart(ki.radicalsScrolledWindow)
	util.ScrollToStart(ki.resultsScrolledWindow)
	for rad := range ki.selectedRadicals {
		ki.unselectRadical(rad)
	}
	for _, b := range ki.radicalButtons {
		b.SetSensitive(true)
	}
	util.RemoveChildren(&ki.results.Container)
	ki.popover.ShowAll()
}

func (ki *KanjiInput) selectRadical(rad string) {
	ki.selectedRadicals[rad] = struct{}{}
	b := ki.radicalButtons[rad]
	child, _ := b.GetChild()
	lbl := (*gtk.Label)(unsafe.Pointer(child))
	lbl.SetMarkup(kanjiLabelMarkup(rad, true))
	b.ShowAll()
}

func (ki *KanjiInput) unselectRadical(rad string) {
	delete(ki.selectedRadicals, rad)
	b := ki.radicalButtons[rad]
	child, _ := b.GetChild()
	lbl := (*gtk.Label)(unsafe.Pointer(child))
	lbl.SetMarkup(kanjiLabelMarkup(rad, false))
	b.ShowAll()
}

func (ki *KanjiInput) toggleRadical(rad string) {
	if _, ok := ki.selectedRadicals[rad]; ok {
		ki.unselectRadical(rad)
	} else {
		ki.selectRadical(rad)
	}
}

func (ki *KanjiInput) updateResults() {
	rads := make([]string, 0, len(ki.selectedRadicals))
	for rad := range ki.selectedRadicals {
		rads = append(rads, rad)
	}

	kanji, krads, err := radicalDict.FetchByRadicals(rads)
	if err != nil {
		log.Printf("Error fetching kanji by radicals: %v", err)
		return
	}
	ki.setResults(kanji)
	ki.setSensitivity(krads)
}

func (ki *KanjiInput) setResults(kanji []kradfile.Kanji) {
	util.RemoveChildren(&ki.results.Container)
	sort.Slice(kanji, func(i, j int) bool {
		ki := kanji[i]
		kj := kanji[j]
		return ki.StrokeCount < kj.StrokeCount ||
			(ki.StrokeCount == kj.StrokeCount && ki.Literal < kj.Literal)
	})

	for _, k := range kanji {
		lbl, _ := gtk.LabelNew(kanjiLabelMarkup(k.Literal, false))
		lbl.SetUseMarkup(true)
		b, _ := gtk.FlowBoxChildNew()
		b.Add(lbl)
		ki.results.Add(b)
	}
	ki.resultKanji = kanji
	ki.results.ShowAll()
}

func (ki *KanjiInput) setSensitivity(krads []string) {
	radSet := make(map[string]struct{}, len(krads))
	for _, rad := range krads {
		radSet[rad] = struct{}{}
	}

	for rad, b := range ki.radicalButtons {
		_, ok := radSet[rad]
		b.SetSensitive(ok)
	}
}

func kanjiLabelMarkup(k string, selected bool) string {
	txt := fmt.Sprintf(`<span size="x-large">%v</span>`, k)
	if selected {
		txt = "<b>" + txt + "</b>"
	}
	return txt
}
