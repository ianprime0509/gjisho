package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/kanjidic"
)

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
