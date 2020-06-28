FILES=\
	bindata.go \
	entry.go example.go gui.go kanji.go main.go search.go \
	internal/util/util.go \
	jmdict/jmdict.go \
	kanjidic/kanjidic.go \
	kanjivg/kanjivg.go \
	tatoeba/tatoeba.go

.PHONY: all clean install

all: gjisho

clean:
	rm -f gjisho

install: gjisho app/gjisho.desktop
	mkdir -p /usr/local/bin
	cp gjisho /usr/local/bin/gjisho
	mkdir -p /usr/local/share/applications
	cp app/gjisho.desktop /usr/local/share/applications

bindata.go: data/gjisho.glade
	go generate

gjisho: ${FILES}
	go build -tags fts5
