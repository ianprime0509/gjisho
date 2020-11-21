.POSIX:
.SUFFIXES:
.PHONY: all check clean fetch install install-program install-database programs

PREFIX=/usr/local
ICON_DIR=${PREFIX}/share/icons

APP_ID=xyz.ianjohnson.GJisho

GJISHO=cmd/gjisho/gjisho
GJISHO_CLI=cmd/gjisho-cli/gjisho-cli

CONVERT=convert
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

SEARCH_PROVIDER_SKELETON=\
	gui/shell-search-provider2.h \
	gui/shell-search-provider2.c

GJISHO_SOURCES=\
	cmd/gjisho/gjisho.go \
	${COMMON_SOURCES} \
	${SEARCH_PROVIDER_SKELETON} \
	gui/bindata.go \
	gui/entry.go \
	gui/event.go \
	gui/example.go \
	gui/gui.go \
	gui/kanji.go \
	gui/search.go \
	gui/search-provider.go \
	gui/util.go \
	gui/application.h \
	gui/application.c

GUI_BINDATA_SOURCES=\
	gui/data/gjisho.glade \
	gui/data/kanji-icon.svg \
	gui/data/logo.svg

GJISHO_CLI_SOURCES=\
	cmd/gjisho-cli/gjisho-cli.go \
	${COMMON_SOURCES}

all: ${GJISHO} ${GJISHO_CLI} gjisho.sqlite

check: ${GJISHO_SOURCES} ${GJISHO_CLI_SOURCES}
	env CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' ${GO} test ./...

clean:
	rm -f ${GJISHO} ${GJISHO_CLI} gjisho.sqlite \
		${SEARCH_PROVIDER_SKELETON} gui/bindata.go gui/data/logo.svg

fetch:
	${CURL} -L ${TATOEBA_URL} | ${GZIP} -d >'${TATOEBA_FILE}'
	${CURL} -L ${JMDICT_URL} | ${GZIP} -d >'${JMDICT_FILE}'
	${CURL} -L ${KANJIDIC_URL} | ${GZIP} -d >'${KANJIDIC_FILE}'
	${CURL} -L ${KANJIVG_URL} | ${GZIP} -d >'${KANJIVG_FILE}'
	${CURL} -L ${KRADFILE_URL} | ${GZIP} -d | iconv -f euc-jp -t utf-8 >'${KRADFILE_FILE}'

install: install-database install-programs

install-database: gjisho.sqlite
	mkdir -p '${DESTDIR}${PREFIX}/share/gjisho'
	cp gjisho.sqlite '${DESTDIR}${PREFIX}/share/gjisho'

install-programs: ${GJISHO} ${GJISHO_CLI} ${APP_ID}.desktop ${APP_ID}.service ${APP_ID}.search-provider.ini
	mkdir -p '${DESTDIR}${PREFIX}/bin'
	cp ${GJISHO} '${DESTDIR}${PREFIX}/bin'
	cp ${GJISHO_CLI} '${DESTDIR}${PREFIX}/bin'
	mkdir -p '${DESTDIR}${PREFIX}/share/applications'
	sed 's_%BIN%_${DESTDIR}${PREFIX}/bin/gjisho_' ${APP_ID}.desktop >'${DESTDIR}${PREFIX}/share/applications/${APP_ID}.desktop'
	mkdir -p '${DESTDIR}${ICON_DIR}/hicolor/scalable/apps'
	cp logo.svg '${DESTDIR}${ICON_DIR}/hicolor/scalable/apps/${APP_ID}.svg'
	for size in 48x48 128x128 192x192 256x256 512x512; do \
		mkdir -p '${DESTDIR}${ICON_DIR}'/hicolor/$$size/apps; \
		${CONVERT} -size $$size -background none logo.svg '${DESTDIR}${ICON_DIR}'/hicolor/$$size/apps/${APP_ID}.png; \
	done
	mkdir -p '${DESTDIR}${PREFIX}/share/dbus-1/services'
	sed 's_%BIN%_${DESTDIR}${PREFIX}/bin/gjisho_' ${APP_ID}.service >'${DESTDIR}${PREFIX}/share/dbus-1/services/${APP_ID}.service'
	mkdir -p '${DESTDIR}${PREFIX}/share/gnome-shell/search-providers'
	cp ${APP_ID}.search-provider.ini '${DESTDIR}${PREFIX}/share/gnome-shell/search-providers'

programs: ${GJISHO} ${GJISHO_CLI}

gui/bindata.go ${SEARCH_PROVIDER_SKELETON}: ${GUI_BINDATA_SOURCES}
	${GO} generate ./gui

gui/data/logo.svg: logo.svg
	cp logo.svg gui/data/logo.svg

${GJISHO}: ${GJISHO_SOURCES}
	cd cmd/gjisho && env CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' ${GO} build -tags fts5

${GJISHO_CLI}: ${GJISHO_CLI_SOURCES}
	cd cmd/gjisho-cli && env CGO_CFLAGS='${CGO_CFLAGS}' CGO_LDFLAGS='${CGO_LDFLAGS}' ${GO} build -tags fts5

gjisho.sqlite: ${GJISHO_CLI} ${TATOEBA_FILE} ${JMDICT_FILE} ${KANJIDIC_FILE} ${KANJIVG_FILE} ${KRADFILE_FILE}
	${GJISHO_CLI} convert \
		-tatoeba '${TATOEBA_FILE}' \
		-jmdict '${JMDICT_FILE}' \
		-kanjidic '${KANJIDIC_FILE}' \
		-kanjivg '${KANJIVG_FILE}' \
		-kradfile '${KRADFILE_FILE}' \
		gjisho.sqlite
