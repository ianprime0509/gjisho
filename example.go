package main

import (
	"context"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/tatoeba"
)

// ExampleDetails is a modal window showing additional details about an example.
type ExampleDetails struct {
	window         *gtk.Window
	japaneseLabel  *gtk.Label
	englishLabel   *gtk.Label
	scrolledWindow *gtk.ScrolledWindow
	words          []jmdict.LookupResult
	wordsList      *gtk.ListBox
	cancelPrevious context.CancelFunc
}

// FetchAndDisplay fetches additional information for the given example,
// displays it in the details modal and shows it.
func (ed *ExampleDetails) FetchAndDisplay(ex tatoeba.Example) {
	ctx := ed.startDisplay()

	indices := ex.UniqueIndices()
	ch := make(chan struct {
		idx tatoeba.Index
		res jmdict.LookupResult
	})
	errCh := make(chan struct{})
	for _, idx := range indices {
		go func(idx tatoeba.Index) {
			ref := idx.Word
			if idx.Disambiguation != "" {
				ref += "・" + idx.Disambiguation
			}
			if e, err := dict.LookupByRef(ref); err == nil {
				ch <- struct {
					idx tatoeba.Index
					res jmdict.LookupResult
				}{idx: idx, res: e}
			} else {
				log.Printf("Error fetching results for reference %q: %v", ref, err)
				errCh <- struct{}{}
			}
		}(idx)
	}

	go func() {
		fetched := make(map[tatoeba.Index]jmdict.LookupResult, len(indices))
		for range indices {
			select {
			case res := <-ch:
				fetched[res.idx] = res.res
			case <-errCh:
			case <-ctx.Done():
				return
			}
		}

		words := make([]jmdict.LookupResult, 0, len(indices))
		for _, idx := range indices {
			if word, ok := fetched[idx]; ok {
				words = append(words, word)
			}
		}
		glib.IdleAdd(func() { ed.display(ex, words) })
	}()
}

// Close closes the example details modal.
func (ed *ExampleDetails) Close() {
	ed.window.Close()
}

func (ed *ExampleDetails) display(ex tatoeba.Example, words []jmdict.LookupResult) {
	jpText := new(strings.Builder)
	for _, seg := range ex.Segments() {
		if jpText.Len() > 0 {
			jpText.WriteRune(' ')
		}
		jpText.WriteString(seg.Text)
	}
	ed.japaneseLabel.SetText(jpText.String())
	ed.englishLabel.SetText(ex.English)

	ed.words = words
	RemoveChildren(&ed.wordsList.Container)
	for _, w := range ed.words {
		ed.wordsList.Add(NewSearchResult(w))
	}
	ed.wordsList.ShowAll()
	ScrollToStart(ed.scrolledWindow)

	ed.window.Present()
}

func (ed *ExampleDetails) startDisplay() context.Context {
	if ed.cancelPrevious != nil {
		ed.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	ed.cancelPrevious = cancel
	return ctx
}
