package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
	"github.com/ianprime0509/gjisho/jmdict"
)

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