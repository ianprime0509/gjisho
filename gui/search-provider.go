package gui

// #include "application.h"
import "C"

import (
	"log"
	"strconv"
	"unsafe"
)

//export gjisho_search_fetch_result_ids
func gjisho_search_fetch_result_ids(query *C.gchar, data *C.GJishoSearchCallbackData) {
	ids := fetchResultIds(C.GoString(query))
	cIDs := make([]*C.gchar, 0, len(ids))
	for _, id := range ids {
		cIDs = append(cIDs, C.CString(id))
	}
	cIDs = append(cIDs, nil)
	C.gjisho_search_result_ids_cb((**C.gchar)(unsafe.Pointer(&cIDs[0])), data)
}

//export gjisho_search_fetch_result_metas
func gjisho_search_fetch_result_metas(cIDs **C.gchar, data *C.GJishoSearchCallbackData) {
	var ids []string
	for i := uintptr(0); ; i++ {
		cID := *(**C.gchar)(unsafe.Pointer(uintptr(unsafe.Pointer(cIDs)) + i*unsafe.Sizeof(*cIDs)))
		if cID == nil {
			break
		}
		ids = append(ids, C.GoString(cID))
	}

	metas := fetchResultMetas(ids)

	cMetas := make([]*C.GJishoSearchResultMeta, 0, len(metas))
	for _, meta := range metas {
		cMeta := (*C.GJishoSearchResultMeta)(C.g_malloc(C.gsize(unsafe.Sizeof(*cMetas[0]))))
		cMeta.id = C.CString(meta.id)
		cMeta.name = C.CString(meta.name)
		cMeta.description = C.CString(meta.description)
		cMetas = append(cMetas, cMeta)
	}
	cMetas = append(cMetas, nil)

	C.gjisho_search_result_metas_cb((**C.GJishoSearchResultMeta)(unsafe.Pointer(&cMetas[0])), data)
}

//export gjisho_launch_for_result_id
func gjisho_launch_for_result_id(id *C.gchar) {
	app.Emit("activate")
	onNavigate(C.GoString(id))
}

//export gjisho_launch_for_search
func gjisho_launch_for_search(query *C.gchar) {
	app.Emit("activate")
	onSearch(C.GoString(query))
}

type searchResultMeta struct {
	id          string
	name        string
	description string
}

func fetchResultIds(query string) []string {
	res, err := db.JMdict.Lookup(query, 0, 10)
	if err != nil {
		log.Printf("Error fetching result IDs for search: %v", err)
		return nil
	}

	ids := make([]string, 0, len(res))
	for _, r := range res {
		ids = append(ids, strconv.Itoa(r.ID))
	}
	return ids
}

func fetchResultMetas(ids []string) []searchResultMeta {
	metas := make([]searchResultMeta, 0, len(ids))
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Invalid result ID: %v", err)
			continue
		}

		ent, err := db.JMdict.Fetch(id)
		if err != nil {
			log.Printf("Error fetching entry data for ID %v: %v", id, err)
			continue
		}

		var desc string
		if ent.PrimaryReading() != ent.Heading() {
			desc += ent.PrimaryReading()
		}
		if len(desc) > 0 {
			desc += " — "
		}
		desc += ent.GlossSummary()

		metas = append(metas, searchResultMeta{
			id:          idStr,
			name:        ent.Heading(),
			description: desc,
		})
	}
	return metas
}
