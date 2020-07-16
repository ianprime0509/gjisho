package gui

import (
	"context"
	"log"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/tatoeba"
)

// exampleDetailsModal is a modal window showing additional details about an example.
type exampleDetailsModal struct {
	window         *gtk.Window
	japaneseLabel  *gtk.Label
	englishLabel   *gtk.Label
	scrolledWindow *gtk.ScrolledWindow
	words          []jmdict.LookupResult
	wordsList      *gtk.ListBox
	cancelPrevious context.CancelFunc
}

// fetchAndDisplay fetches additional information for the given example,
// displays it in the details modal and shows it.
func (ed *exampleDetailsModal) fetchAndDisplay(ex tatoeba.Example) {
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
			if e, err := db.JMdict.LookupByRef(ref); err == nil {
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

// close closes the example details modal.
func (ed *exampleDetailsModal) close() {
	ed.window.Close()
}

func (ed *exampleDetailsModal) display(ex tatoeba.Example, words []jmdict.LookupResult) {
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
	removeChildren(&ed.wordsList.Container)
	for _, w := range ed.words {
		ed.wordsList.Add(newSearchResult(w))
	}
	ed.wordsList.ShowAll()
	scrollToStart(ed.scrolledWindow)

	ed.window.Present()
}

func (ed *exampleDetailsModal) startDisplay() context.Context {
	if ed.cancelPrevious != nil {
		ed.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	ed.cancelPrevious = cancel
	return ctx
}
