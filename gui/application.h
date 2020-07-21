#ifndef GJISHO_APPLICATION_H
#define GJISHO_APPLICATION_H

#include <gtk/gtk.h>

G_BEGIN_DECLS

#define GJISHO_TYPE_APPLICATION gjisho_application_get_type()

G_DECLARE_FINAL_TYPE(GJishoApplication, gjisho_application, GJISHO, APPLICATION, GtkApplication);

GJishoApplication *gjisho_application_new(void);

typedef struct GJishoSearchCallbackData GJishoSearchCallbackData;

typedef struct {
    gchar *id;
    gchar *name;
    gchar *description;
} GJishoSearchResultMeta;

void gjisho_search_result_ids_cb(gchar **ids, GJishoSearchCallbackData *data);
void gjisho_search_result_metas_cb(GJishoSearchResultMeta **metas, GJishoSearchCallbackData *data);

G_END_DECLS

#endif
