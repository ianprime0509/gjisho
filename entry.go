package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/internal/util"
	"github.com/ianprime0509/gjisho/jmdict"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/tatoeba"
)

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
	util.RemoveChildren(&lst.list.Container)
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
	util.RemoveChildren(&lst.list.Container)
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
