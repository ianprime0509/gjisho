.POSIX:
.SUFFIXES:
.PHONY: all clean fetch install install-program install-database

PREFIX=/usr/local

CURL=curl
GO=go
GZIP=gzip

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
	${CURL} -L ${TATOEBA_URL} | ${GZIP} -d >'${TATOEBA_FILE}'
	${CURL} -L ${JMDICT_URL} | ${GZIP} -d >'${JMDICT_FILE}'
	${CURL} -L ${KANJIDIC_URL} | ${GZIP} -d >'${KANJIDIC_FILE}'
	${CURL} -L ${KANJIVG_URL} | ${GZIP} -d >'${KANJIVG_FILE}'

install: install-database install-program

install-database: gjisho.sqlite
	mkdir -p '${DESTDIR}${PREFIX}/share/gjisho'
	cp gjisho.sqlite '${DESTDIR}${PREFIX}/share/gjisho'

install-program: gjisho gjisho.desktop
	mkdir -p '${DESTDIR}${PREFIX}/bin'
	cp gjisho '${DESTDIR}${PREFIX}/bin'
	mkdir -p '${DESTDIR}${PREFIX}/share/applications'
	cp gjisho.desktop '${DESTDIR}${PREFIX}/share/applications'

bindata.go: data/gjisho.glade
	${GO} generate

gjisho: ${SOURCE_FILES}
	${GO} build -tags fts5

gjisho.sqlite: gjisho ${TATOEBA_FILE} ${JMDICT_FILE} ${KANJIDIC_FILE} ${KANJIVG_FILE}
	./gjisho -conv gjisho.sqlite \
		-tatoeba '${TATOEBA_FILE}' \
		-jmdict '${JMDICT_FILE}' \
		-kanjidic '${KANJIDIC_FILE}' \
		-kanjivg '${KANJIVG_FILE}'
