.POSIX:
.SUFFIXES:
.PHONY: all clean fetch install

PREFIX=/usr/local

CURL=curl
GO=go

TATOEBA_FILE=raw/examples.utf
TATOEBA_URL=ftp://ftp.monash.edu/pub/nihongo/examples.utf.gz
JMDICT_FILE=raw/JMdict
JMDICT_URL=http://ftp.monash.edu/pub/nihongo/JMdict.gz
KANJIDIC_FILE=raw/kanjidic2.xml
KANJIDIC_URL=http://www.edrdg.org/kanjidic/kanjidic2.xml.gz
KANJIVG_FILE=raw/kanjivg-20160426.xml
KANJIVG_URL=https://github.com/KanjiVG/kanjivg/releases/download/r20160426/kanjivg-20160426.xml.gz

SOURCE_FILES=\
	bindata.go \
	entry.go example.go gui.go kanji.go main.go search.go \
	internal/util/util.go \
	jmdict/jmdict.go \
	kanjidic/kanjidic.go \
	kanjivg/kanjivg.go \
	tatoeba/tatoeba.go

all: gjisho gjisho.sqlite

clean:
	rm -f gjisho gjisho.sqlite

fetch:
	curl -L ${TATOEBA_URL} | gzip -d >${TATOEBA_FILE}
	curl -L ${JMDICT_URL} | gzip -d >${JMDICT_FILE}
	curl -L ${KANJIDIC_URL} | gzip -d >${KANJIDIC_FILE}
	curl -L ${KANJIVG_URL} | gzip -d >${KANJIVG_FILE}

install: gjisho gjisho.sqlite app/gjisho.desktop
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp gjisho ${DESTDIR}${PREFIX}/bin
	mkdir -p ${DESTDIR}${PREFIX}/share/gjisho
	cp gjisho.sqlite ${DESTDIR}${PREFIX}/share/gjisho
	mkdir -p ${DESTDIR}${PREFIX}/share/applications
	cp app/gjisho.desktop ${DESTDIR}${PREFIX}/share/applications

bindata.go: data/gjisho.glade
	${GO} generate

gjisho: ${SOURCE_FILES}
	${GO} build -tags fts5

gjisho.sqlite: gjisho ${TATOEBA_FILE} ${JMDICT_FILE} ${KANJIDIC_FILE} ${KANJIVG_FILE}
	./gjisho -conv gjisho.sqlite \
		-tatoeba ${TATOEBA_FILE} \
		-jmdict ${JMDICT_FILE} \
		-kanjidic ${KANJIDIC_FILE} \
		-kanjivg ${KANJIVG_FILE}
