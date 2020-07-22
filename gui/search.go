package gui

import (
	"context"
	"fmt"
	"log"
	"sort"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/kradfile"
)

// appSearch is a wrapper around the search-related components of the app.
type appSearch struct {
	toggle         *gtk.ToggleButton
	revealer       *gtk.Revealer
	entry          *gtk.SearchEntry
	kanjiInput     *kanjiInput
	results        *searchResultList
	resultsKanji   *searchResultsKanji
	cancelPrevious context.CancelFunc // a function to cancel the previous search
}

// search searches using the given query.
func (s *appSearch) search(query string) {
	ctx := s.startSearch()
	s.resultsKanji.fetchAndDisplay(ctx, query)
	ch := make(chan []jmdict.LookupResult)

	go func() {
		if results, err := db.JMdict.Lookup(query); err == nil {
			ch <- results
		} else {
			log.Printf("Lookup query error: %v", err)
		}
	}()

	go func() {
		select {
		case results := <-ch:
			glib.IdleAdd(func() { s.results.set(results) })
		case <-ctx.Done():
		}
	}()
}

// toggleOpen toggles whether the search pane is open.
func (s *appSearch) toggleOpen() {
	s.revealer.SetRevealChild(!s.revealer.GetRevealChild())
}

// activateEntry activates and focuses the search entry.
func (s *appSearch) activateEntry() {
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

// deactivateEntry deactivates the search entry.
func (s *appSearch) deactivateEntry() {
	s.toggle.SetActive(false)
	s.revealer.SetRevealChild(false)
}

// insertEntryText inserts text in the search entry buffer at the current
// position.
func (s *appSearch) insertEntryText(text string) {
	old, _ := search.entry.GetText()
	rs := []rune(old)
	pos := search.entry.GetPosition()
	search.entry.SetText(string(rs[:pos]) + text + string(rs[pos:]))
	search.entry.SetPosition(pos + utf8.RuneCountInString(text))
}

func (s *appSearch) startSearch() context.Context {
	if s.cancelPrevious != nil {
		s.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelPrevious = cancel
	return ctx
}

// searchResultList is a list of search results displayed in the GUI.
type searchResultList struct {
	list           *gtk.ListBox
	scrolledWindow *gtk.ScrolledWindow
	results        []jmdict.LookupResult
	nDisplayed     int
}

// clearSelection clears the currently selected result.
func (lst *searchResultList) clearSelection() {
	lst.list.SelectRow(nil)
}

// selected returns the currently selected search result, or nil if none is
// selected.
func (lst *searchResultList) selected() *jmdict.LookupResult {
	if row := lst.list.GetSelectedRow(); row != nil {
		return &lst.results[row.GetIndex()]
	}
	return nil
}

// set sets the currently displayed search results.
func (lst *searchResultList) set(results []jmdict.LookupResult) {
	lst.results = results
	removeChildren(&lst.list.Container)
	lst.nDisplayed = 0
	lst.showMore()
	scrollToStart(lst.scrolledWindow)
}

// showMore displays more search results in the list.
func (lst *searchResultList) showMore() {
	maxIndex := lst.nDisplayed + 50
	for ; lst.nDisplayed < len(lst.results) && lst.nDisplayed < maxIndex; lst.nDisplayed++ {
		lst.list.Add(newSearchResult(lst.results[lst.nDisplayed]))
	}
	lst.list.ShowAll()
}

// newSearchResult creates a search result widget for display.
func newSearchResult(entry jmdict.LookupResult) *gtk.Box {
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

// searchResultsKanji is a display of kanji related to the search query.
type searchResultsKanji struct {
	box   *gtk.FlowBox
	kanji []kanjidic.Character
}

// fetchAndDisplay fetches and displays information about the kanji related to
// the given search query.
func (srk *searchResultsKanji) fetchAndDisplay(ctx context.Context, query string) {
	kanji := associatedKanji(query)
	// Kanji lookup is fast enough that for now I'm just fetching in sequence
	ch := make(chan []kanjidic.Character)

	go func() {
		res := make([]kanjidic.Character, 0, len(kanji))
		for _, k := range kanji {
			if c, err := db.KANJIDIC.Fetch(k); err == nil {
				res = append(res, c)
			} else {
				log.Printf("Error fetching kanji details for %q: %v", k, err)
			}
		}
		ch <- res
	}()

	go func() {
		select {
		case res := <-ch:
			glib.IdleAdd(func() { srk.display(res) })
		case <-ctx.Done():
		}
	}()
}

func (srk *searchResultsKanji) display(kanji []kanjidic.Character) {
	removeChildren(&srk.box.Container)
	for _, k := range kanji {
		lbl, _ := gtk.LabelNew(fmt.Sprintf(`<span size="large">%v</span>`, k.Literal))
		lbl.SetUseMarkup(true)
		srk.box.Add(lbl)
	}
	srk.box.ShowAll()
	srk.kanji = kanji
}

func associatedKanji(query string) []string {
	// Very similar to Entry.AssociatedKanji. We want the resulting list of
	// related kanji to be in order, so we make the values of the map the
	// indices of the kanji in the final list.
	set := make(map[rune]int)
	idx := 0
	for _, c := range query {
		if unicode.Is(unicode.Han, c) {
			if _, ok := set[c]; !ok {
				set[c] = idx
				idx++
			}
		}
	}

	kanji := make([]string, len(set))
	for k, i := range set {
		kanji[i] = string(k)
	}
	return kanji
}

// kanjiInput is a special input popover to make it easier to input kanji.
type kanjiInput struct {
	button                 *gtk.ToggleButton
	buttonIcon             *gtk.Image
	popover                *gtk.Popover
	radicalsScrolledWindow *gtk.ScrolledWindow
	radicalsBox            *gtk.Box
	radicalButtons         map[string]*gtk.FlowBoxChild
	selectedRadicals       map[string]struct{}
	resultsScrolledWindow  *gtk.ScrolledWindow
	resultsBox             *gtk.FlowBox
	results                []kradfile.Kanji
	cancelPrevious         context.CancelFunc
}

// initRadicals initializes the radical input buttons.
func (ki *kanjiInput) initRadicals() {
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
		s := s
		var heading string
		if s == 1 {
			heading = "1 stroke"
		} else {
			heading = fmt.Sprintf("%v strokes", s)
		}
		lbl, _ := gtk.LabelNew(heading)
		ki.radicalsBox.Add(lbl)

		fb, _ := gtk.FlowBoxNew()
		fb.SetMinChildrenPerLine(10)
		fb.SetMaxChildrenPerLine(10)
		fb.SetSelectionMode(gtk.SELECTION_NONE)
		for _, rad := range rads[s] {
			lbl, _ := gtk.LabelNew(kanjiLabelMarkup(rad, false))
			lbl.SetUseMarkup(true)
			b, _ := gtk.FlowBoxChildNew()
			b.Add(lbl)
			ki.radicalButtons[rad] = b
			b.Connect("focus", func() {
				fby := fb.GetAllocation().GetY()
				by := b.GetAllocation().GetY()
				ki.radicalsScrolledWindow.GetVAdjustment().SetValue(float64(fby + by))
			})
			fb.Add(b)
		}
		// For some reason, connecting the activate signal of one of the
		// children directly doesn't seem to work with mouse clicks (it works
		// with keyboard activations, though). To get around this, we can just
		// use the child-activated signal on the FlowBox, which does work.
		fb.Connect("child-activated", func(_ interface{}, c *gtk.FlowBoxChild) {
			if c.GetSensitive() {
				rad := rads[s][c.GetIndex()]
				ki.toggleRadical(rad)
				ki.updateResults()
			}
		})
		ki.radicalsBox.Add(fb)
	}
}

// display displays the kanji input popover.
func (ki *kanjiInput) display() {
	scrollToStart(ki.radicalsScrolledWindow)
	scrollToStart(ki.resultsScrolledWindow)
	for rad := range ki.selectedRadicals {
		ki.unselectRadical(rad)
	}
	for _, b := range ki.radicalButtons {
		b.SetSensitive(true)
	}
	removeChildren(&ki.resultsBox.Container)
	ki.popover.ShowAll()
}

func (ki *kanjiInput) selectRadical(rad string) {
	ki.selectedRadicals[rad] = struct{}{}
	b := ki.radicalButtons[rad]
	child, _ := b.GetChild()
	lbl := (*gtk.Label)(unsafe.Pointer(child))
	lbl.SetMarkup(kanjiLabelMarkup(rad, true))
	b.ShowAll()
}

func (ki *kanjiInput) unselectRadical(rad string) {
	delete(ki.selectedRadicals, rad)
	b := ki.radicalButtons[rad]
	child, _ := b.GetChild()
	lbl := (*gtk.Label)(unsafe.Pointer(child))
	lbl.SetMarkup(kanjiLabelMarkup(rad, false))
	b.ShowAll()
}

func (ki *kanjiInput) toggleRadical(rad string) {
	if _, ok := ki.selectedRadicals[rad]; ok {
		ki.unselectRadical(rad)
	} else {
		ki.selectRadical(rad)
	}
}

func (ki *kanjiInput) updateResults() {
	ctx := ki.startDisplay()
	ch := make(chan struct {
		kanji []kradfile.Kanji
		krads []string
	})

	go func() {
		rads := make([]string, 0, len(ki.selectedRadicals))
		for rad := range ki.selectedRadicals {
			rads = append(rads, rad)
		}

		if kanji, krads, err := db.KRADFILE.FetchByRadicals(rads); err == nil {
			ch <- struct {
				kanji []kradfile.Kanji
				krads []string
			}{kanji, krads}
		} else {
			log.Printf("Error fetching kanji by radicals: %v", err)
		}
	}()

	go func() {
		select {
		case res := <-ch:
			glib.IdleAdd(func() {
				ki.setResults(res.kanji)
				ki.setSensitivity(res.krads)
			})
		case <-ctx.Done():
		}
	}()
}

func (ki *kanjiInput) setResults(kanji []kradfile.Kanji) {
	removeChildren(&ki.resultsBox.Container)
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
		ki.resultsBox.Add(b)
		b.Connect("focus", func() {
			by := b.GetAllocation().GetY()
			ki.resultsScrolledWindow.GetVAdjustment().SetValue(float64(by))
		})
	}
	ki.results = kanji
	ki.resultsBox.ShowAll()
}

func (ki *kanjiInput) setSensitivity(krads []string) {
	radSet := make(map[string]struct{}, len(krads))
	for _, rad := range krads {
		radSet[rad] = struct{}{}
	}

	for rad, b := range ki.radicalButtons {
		_, ok := radSet[rad]
		b.SetSensitive(ok)
	}
}

func (ki *kanjiInput) startDisplay() context.Context {
	if ki.cancelPrevious != nil {
		ki.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	ki.cancelPrevious = cancel
	return ctx
}

func kanjiLabelMarkup(k string, selected bool) string {
	txt := fmt.Sprintf(`<span size="x-large">%v</span>`, k)
	if selected {
		txt = "<b>" + txt + "</b>"
	}
	return txt
}
