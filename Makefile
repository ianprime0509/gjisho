.POSIX:
.SUFFIXES:
.PHONY: all check clean fetch install install-program install-database

PREFIX=/usr/local

GJISHO=cmd/gjisho/gjisho
GJISHO_CLI=cmd/gjisho-cli/gjisho-cli

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
KRADFILE_FILE=raw/kradfile
KRADFILE_URL=ftp://ftp.monash.edu/pub/nihongo/kradfile.gz

COMMON_SOURCES=\
	datautil/datautil.go \
	dictdb/dictdb.go \
	jmdict/jmdict.go \
	kanjidic/kanjidic.go \
	kanjivg/kanjivg.go \
	kradfile/kradfile.go kradfile/radstrokes.go \
	tatoeba/tatoeba.go

GJISHO_SOURCES=\
	cmd/gjisho/gjisho.go \
	${COMMON_SOURCES} \
	gui/bindata.go \
	gui/entry.go \
	gui/event.go \
	gui/example.go \
	gui/gui.go \
	gui/kanji.go \
	gui/search.go \
	gui/util.go

GJISHO_CLI_SOURCES=\
	cmd/gjisho-cli/gjisho-cli.go \
	${COMMON_SOURCES}

all: ${GJISHO} ${GJISHO_CLI} gjisho.sqlite

check: ${GJISHO_SOURCES} ${GJISHO_CLI_SOURCES}
	${GO} test ./...

clean:
	rm -f ${GJISHO} ${GJISHO_CLI} gjisho.sqlite

fetch:
	${CURL} -L ${TATOEBA_URL} | ${GZIP} -d >'${TATOEBA_FILE}'
	${CURL} -L ${JMDICT_URL} | ${GZIP} -d >'${JMDICT_FILE}'
	${CURL} -L ${KANJIDIC_URL} | ${GZIP} -d >'${KANJIDIC_FILE}'
	${CURL} -L ${KANJIVG_URL} | ${GZIP} -d >'${KANJIVG_FILE}'
	${CURL} -L ${KRADFILE_URL} | ${GZIP} -d | iconv -f euc-jp -t utf-8 >'${KRADFILE_FILE}'

install: install-database install-program

install-database: gjisho.sqlite
	mkdir -p '${DESTDIR}${PREFIX}/share/gjisho'
	cp gjisho.sqlite '${DESTDIR}${PREFIX}/share/gjisho'

install-programs: ${GJISHO} ${GJISHO_CLI} gjisho.desktop
	mkdir -p '${DESTDIR}${PREFIX}/bin'
	cp ${GJISHO} '${DESTDIR}${PREFIX}/bin'
	cp ${GJISHO_CLI} '${DESTDIR}${PREFIX}/bin'
	mkdir -p '${DESTDIR}${PREFIX}/share/applications'
	cp gjisho.desktop '${DESTDIR}${PREFIX}/share/applications'

bindata.go: data/gjisho.glade
	${GO} generate

${GJISHO}: ${GJISHO_SOURCES}
	cd cmd/gjisho && ${GO} build -tags fts5

${GJISHO_CLI}: ${GJISHO_CLI_SOURCES}
	cd cmd/gjisho-cli && ${GO} build -tags fts5

gjisho.sqlite: ${GJISHO_CLI} ${TATOEBA_FILE} ${JMDICT_FILE} ${KANJIDIC_FILE} ${KANJIVG_FILE} ${KRADFILE_FILE}
	${GJISHO_CLI} convert \
		-tatoeba '${TATOEBA_FILE}' \
		-jmdict '${JMDICT_FILE}' \
		-kanjidic '${KANJIDIC_FILE}' \
		-kanjivg '${KANJIVG_FILE}' \
		-kradfile '${KRADFILE_FILE}' \
		gjisho.sqlite
