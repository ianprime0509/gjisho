/*
 * General idea largely copied from https://developer.gnome.org/SearchProvider/.
 * I wish it were possible to do this directly from Go, but such support hasn't
 * been added yet, and I'm not yet knowledgeable enough to know how to add it.
 *
 * Help was also taken from the Nautilus implementation, which is linked on the
 * above page:
 * https://gitlab.gnome.org/GNOME/nautilus/blob/master/src/nautilus-shell-search-provider.c
 */

#include "application.h"
#include "shell-search-provider2.h"

#define GJISHO_APP_ID "xyz.ianjohnson.GJisho"

/*
 * Implemented in search-provider.go. These functions will do the work of
 * performing the relevant searches or lookups and will call the applicable
 * callback functions (below) once done.
 */
extern void gjisho_search_fetch_result_ids(gchar *query, GJishoSearchCallbackData *data);
extern void gjisho_search_fetch_result_metas(gchar **ids, GJishoSearchCallbackData *data);
/*
 * Implemented in search-provider.go. These functions do not have callbacks,
 * since they launch the application directly.
 */
extern void gjisho_launch_for_result_id(gchar *id);
extern void gjisho_launch_for_search(gchar *query);

struct GJishoSearchCallbackData {
	gboolean is_subsearch;
	ShellSearchProvider2 *provider;
	GDBusMethodInvocation *invocation;
};

void
gjisho_search_result_ids_cb(gchar **ids, GJishoSearchCallbackData *data)
{
	int i;

	if (data->is_subsearch)
		shell_search_provider2_complete_get_subsearch_result_set(
			data->provider, data->invocation, (const gchar *const *)ids);
	else
		shell_search_provider2_complete_get_initial_result_set(
			data->provider, data->invocation, (const gchar *const *)ids);


	for (i = 0; ids[i] != NULL; i++)
		g_free(ids[i]);
	g_application_release(g_application_get_default());
	g_object_unref(data->invocation);
	g_free(data);
}

void
gjisho_search_result_metas_cb(GJishoSearchResultMeta **metas, GJishoSearchCallbackData *data)
{
	GVariantBuilder builder;
	gsize i;

	g_variant_builder_init(&builder, G_VARIANT_TYPE("aa{sv}"));
	for (i = 0; metas[i] != NULL; i++) {
		g_variant_builder_open(&builder, G_VARIANT_TYPE("a{sv}"));
		g_variant_builder_add(&builder, "{sv}", "id", g_variant_new_string(metas[i]->id));
		g_free(metas[i]->id);
		g_variant_builder_add(&builder, "{sv}", "name", g_variant_new_string(metas[i]->name));
		g_free(metas[i]->name);
		g_variant_builder_add(&builder, "{sv}", "description", g_variant_new_string(metas[i]->description));
		g_free(metas[i]->description);
		g_variant_builder_close(&builder);
		g_free(metas[i]);
	}

	shell_search_provider2_complete_get_result_metas(
		data->provider, data->invocation, g_variant_builder_end(&builder));

	g_application_release(g_application_get_default());
	g_object_unref(data->invocation);
	g_free(data);
}

#define GJISHO_TYPE_SEARCH_PROVIDER gjisho_search_provider_get_type()

G_DECLARE_FINAL_TYPE(GJishoSearchProvider, gjisho_search_provider, GJISHO, SEARCH_PROVIDER, GObject);

struct _GJishoSearchProvider {
	GObject parent_instance;
	ShellSearchProvider2 *skeleton;
};

G_DEFINE_TYPE(GJishoSearchProvider, gjisho_search_provider, G_TYPE_OBJECT);

static gboolean
gjisho_search_provider_get_initial_result_set(GJishoSearchProvider *self,
	GDBusMethodInvocation *invocation, gchar **terms, gpointer user_data)
{
	gchar *query = g_strjoinv(" ", terms);
	GJishoSearchCallbackData *data = g_malloc(sizeof(*data));

	data->is_subsearch = FALSE;
	data->provider = self->skeleton;
	data->invocation = g_object_ref(invocation);
	g_application_hold(g_application_get_default());
	gjisho_search_fetch_result_ids(query, data);
	g_free(query);
	return TRUE;
}

static gboolean
gjisho_search_provider_get_subsearch_result_set(GJishoSearchProvider *self,
	GDBusMethodInvocation *invocation, gchar **previous_results,
	gchar **terms, gpointer user_data)
{
	gchar *query = g_strjoinv(" ", terms);
	GJishoSearchCallbackData *data = g_malloc(sizeof(*data));

	data->is_subsearch = TRUE;
	data->provider = self->skeleton;
	data->invocation = g_object_ref(invocation);
	g_application_hold(g_application_get_default());
	gjisho_search_fetch_result_ids(query, data);
	g_free(query);
	return TRUE;
}

static gboolean
gjisho_search_provider_get_result_metas(GJishoSearchProvider *self,
	GDBusMethodInvocation *invocation, gchar **ids, gpointer user_data)
{
	GJishoSearchCallbackData *data = g_malloc(sizeof(*data));

	data->is_subsearch = FALSE;
	data->provider = self->skeleton;
	data->invocation = g_object_ref(invocation);
	g_application_hold(g_application_get_default());
	gjisho_search_fetch_result_metas(ids, data);
	return TRUE;
}

static gboolean
gjisho_search_provider_activate_result(GJishoSearchProvider *self,
	GDBusMethodInvocation *invocation, gchar *id, gchar **terms,
	guint32 timestamp)
{
	gjisho_launch_for_result_id(id);
	shell_search_provider2_complete_activate_result(self->skeleton, invocation);
	return TRUE;
}

static gboolean
gjisho_search_provider_launch_search(GJishoSearchProvider *self,
	GDBusMethodInvocation *invocation, gchar **terms, guint32 timestamp)
{
	gchar *query = g_strjoinv(" ", terms);
	gjisho_launch_for_search(query);
	g_free(query);
	shell_search_provider2_complete_launch_search(self->skeleton, invocation);
	return TRUE;
}

static gboolean
gjisho_search_provider_dbus_export(GJishoSearchProvider *self, GDBusConnection *connection,
	const gchar *object_path, GError **error)
{
	return g_dbus_interface_skeleton_export(G_DBUS_INTERFACE_SKELETON(self->skeleton),
		connection, object_path, error);
}

static void
gjisho_search_provider_dbus_unexport(GJishoSearchProvider *self, GDBusConnection *connection,
	const gchar *object_path)
{
	GDBusInterfaceSkeleton *skeleton = G_DBUS_INTERFACE_SKELETON(self->skeleton);

	if (g_dbus_interface_skeleton_has_connection(skeleton, connection))
		g_dbus_interface_skeleton_unexport_from_connection(skeleton, connection);
}

static void
gjisho_search_provider_class_init(GJishoSearchProviderClass *class)
{
}

static void
gjisho_search_provider_init(GJishoSearchProvider *self)
{
	self->skeleton = shell_search_provider2_skeleton_new();

	g_signal_connect_swapped(self->skeleton, "handle-get-initial-result-set",
		G_CALLBACK(gjisho_search_provider_get_initial_result_set), self);
	g_signal_connect_swapped(self->skeleton, "handle-get-subsearch-result-set",
		G_CALLBACK(gjisho_search_provider_get_subsearch_result_set), self);
	g_signal_connect_swapped(self->skeleton, "handle-get-result-metas",
		G_CALLBACK(gjisho_search_provider_get_result_metas), self);
	g_signal_connect_swapped(self->skeleton, "handle-activate-result",
		G_CALLBACK(gjisho_search_provider_activate_result), self);
	g_signal_connect_swapped(self->skeleton, "handle-launch-search",
		G_CALLBACK(gjisho_search_provider_launch_search), self);
}

struct _GJishoApplication {
	GtkApplication parent_instance;
	GJishoSearchProvider *search_provider;
};

G_DEFINE_TYPE(GJishoApplication, gjisho_application, GTK_TYPE_APPLICATION);

static gboolean
gjisho_application_dbus_register(GApplication *application, GDBusConnection *connection,
	const gchar *object_path, GError **error)
{
	GJishoApplication *self = GJISHO_APPLICATION(application);
	gboolean retval = FALSE;
	gchar *search_provider_path = NULL;

	if (!G_APPLICATION_CLASS(gjisho_application_parent_class)->dbus_register(
		application, connection, object_path, error))
		goto OUT;

	search_provider_path = g_strconcat(object_path, "/SearchProvider", NULL);
	if (!gjisho_search_provider_dbus_export(
		self->search_provider, connection, search_provider_path, error))
		goto OUT;

	retval = TRUE;

OUT:
	g_free(search_provider_path);
	return retval;
}

static void
gjisho_application_dbus_unregister(GApplication *application, GDBusConnection *connection,
	const gchar *object_path)
{
	GJishoApplication *self = GJISHO_APPLICATION(application);
	gchar *search_provider_path = NULL;

	search_provider_path = g_strconcat(object_path, "/SearchProvider", NULL);
	gjisho_search_provider_dbus_unexport(self->search_provider, connection, search_provider_path);

	G_APPLICATION_CLASS(gjisho_application_parent_class)->dbus_unregister(
		application, connection, object_path);

	g_free(search_provider_path);
}

static void
gjisho_application_class_init(GJishoApplicationClass *class)
{
	GApplicationClass *app_class = G_APPLICATION_CLASS(class);

	app_class->dbus_register = gjisho_application_dbus_register;
	app_class->dbus_unregister = gjisho_application_dbus_unregister;
}

static void
gjisho_application_init(GJishoApplication *self)
{
	self->search_provider = g_object_new(GJISHO_TYPE_SEARCH_PROVIDER, NULL);
}

GJishoApplication *
gjisho_application_new(void)
{
	return g_object_new(GJISHO_TYPE_APPLICATION,
		"application-id", GJISHO_APP_ID,
		NULL);
}
