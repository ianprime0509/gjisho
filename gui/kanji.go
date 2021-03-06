package gui

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/ianprime0509/gjisho/kanjidic"
	"github.com/ianprime0509/gjisho/kanjivg"
)

// kanjiDetailsModal is a modal window showing additional details about a kanji.
type kanjiDetailsModal struct {
	window                    *gtk.Window
	scrolledWindow            *gtk.ScrolledWindow
	strokeOrderScrolledWindow *gtk.ScrolledWindow
	strokeOrder               *gtk.Box
	charLabel                 *gtk.Label
	subtitleLabel             *gtk.Label
	readingMeanings           *gtk.Label
	dictRefsLabel             *gtk.Label
	queryCodesLabel           *gtk.Label
	cancelPrevious            context.CancelFunc
}

// fetchAndDisplay fetches additional information about the given kanji and
// displays it in the window (which is then shown).
func (kd *kanjiDetailsModal) fetchAndDisplay(c kanjidic.Character) {
	ctx := kd.startDisplay()
	ch := make(chan kanjivg.Kanji)
	go func() {
		if k, err := db.KanjiVG.Fetch(c.Literal); err == nil {
			ch <- k
		} else {
			log.Printf("Could not fetch kanji stroke information for %q: %v", c.Literal, err)
		}
	}()

	go func() {
		select {
		case k := <-ch:
			glib.IdleAdd(func() { kd.display(c, k) })
		case <-ctx.Done():
		}
	}()
}

func (kd *kanjiDetailsModal) display(c kanjidic.Character, k kanjivg.Kanji) {
	kd.charLabel.SetText(c.Literal)
	kd.drawStrokes(k)
	kd.subtitleLabel.SetMarkup(fmtSubtitle(c))
	sb := new(strings.Builder)
	for _, rm := range c.ReadingMeaningGroups {
		sb.WriteString(fmtReadingMeaning(rm))
		sb.WriteRune('\n')
	}
	kd.readingMeanings.SetMarkup(strings.TrimSpace(sb.String()))
	kd.readingMeanings.ShowAll()
	kd.dictRefsLabel.SetMarkup(fmtDictRefs(c.DictRefs))
	kd.queryCodesLabel.SetMarkup(fmtQueryCodes(c.QueryCodes))
	scrollToStart(kd.scrolledWindow)
	scrollToStart(kd.strokeOrderScrolledWindow)
	kd.window.Present()
}

func (kd *kanjiDetailsModal) startDisplay() context.Context {
	if kd.cancelPrevious != nil {
		kd.cancelPrevious()
	}
	ctx, cancel := context.WithCancel(context.Background())
	kd.cancelPrevious = cancel
	return ctx
}

func (kd *kanjiDetailsModal) drawStrokes(kanji kanjivg.Kanji) {
	drawTo := func(n int) *gtk.DrawingArea {
		da, _ := gtk.DrawingAreaNew()
		da.SetSizeRequest(54, 54)
		da.Connect("draw", func(_ *gtk.DrawingArea, ctx *cairo.Context) {
			num := strconv.Itoa(n + 1)
			extents := ctx.TextExtents(num)
			ctx.MoveTo(-extents.XBearing, -extents.YBearing)
			ctx.ShowText(strconv.Itoa(n + 1))

			ctx.Scale(0.5, 0.5)
			ctx.SetSourceRGB(0.5, 0.5, 0.5)
			for i := 0; i < n; i++ {
				kanji.Strokes[i].DrawTo(ctx, false)
			}
			ctx.SetSourceRGB(0, 0, 0)
			kanji.Strokes[n].DrawTo(ctx, true)
		})
		return da
	}

	removeChildren(&kd.strokeOrder.Container)
	for i := range kanji.Strokes {
		kd.strokeOrder.Add(drawTo(i))
	}
	kd.strokeOrder.ShowAll()
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

func fmtReadingMeaning(rm kanjidic.ReadingMeaningGroup) string {
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

	return sb.String()
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
