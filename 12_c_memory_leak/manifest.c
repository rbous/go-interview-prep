#include "manifest.h"
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

// BUG(1): On parse error mid-way through the manifest, previously allocated
//         entries are not freed (memory leak on error path).
// BUG(2): free_manifest doesn't free the individual entry fields.
// BUG(3): find_entry doesn't check for NULL manifest.

Manifest *parse_manifest(const char *data) {
    if (!data)
        return NULL;

    Manifest *m = malloc(sizeof(Manifest));
    m->head = NULL;
    m->count = 0;

    // Work on a copy since strtok modifies the string
    char *copy = strdup(data);
    char *line = strtok(copy, "\n");

    while (line) {
        // Parse "package:version:size"
        char *name_str = strtok_r(line, ":", &line);
        char *ver_str = strtok_r(NULL, ":", &line);
        char *size_str = strtok_r(NULL, ":", &line);

        if (!name_str || !ver_str || !size_str) {
            // Parse error — but we don't clean up entries already added!
            free(copy);
            free(m);
            return NULL;
        }

        ManifestEntry *entry = malloc(sizeof(ManifestEntry));
        entry->package_name = strdup(name_str);
        entry->version = strdup(ver_str);
        entry->size_bytes = atoi(size_str);
        entry->next = m->head;
        m->head = entry;
        m->count++;

        line = strtok(NULL, "\n");
    }

    free(copy);
    return m;
}

void free_manifest(Manifest *m) {
    if (!m) return;

    ManifestEntry *curr = m->head;
    while (curr) {
        ManifestEntry *next = curr->next;
        // BUG: doesn't free curr->package_name and curr->version
        free(curr);
        curr = next;
    }
    free(m);
}

ManifestEntry *find_entry(const Manifest *m, const char *package_name) {
    ManifestEntry *curr = m->head;
    while (curr) {
        if (strcmp(curr->package_name, package_name) == 0)
            return curr;
        curr = curr->next;
    }
    return NULL;
}
