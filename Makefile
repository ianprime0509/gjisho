FILES=\
	bindata.go \
	entry.go example.go gui.go kanji.go main.go search.go \
	internal/util/util.go \
	jmdict/jmdict.go \
	kanjidic/kanjidic.go \
	kanjivg/kanjivg.go \
	tatoeba/tatoeba.go

.PHONY: all clean

all: gjisho

clean:
	rm -f gjisho

bindata.go: data/gjisho.glade
	go generate

gjisho: ${FILES}
	go build -tags fts5
