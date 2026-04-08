#ifndef MANIFEST_H
#define MANIFEST_H

typedef struct ManifestEntry {
    char *package_name;   // heap-allocated
    char *version;        // heap-allocated
    int   size_bytes;
    struct ManifestEntry *next;
} ManifestEntry;

typedef struct {
    ManifestEntry *head;
    int count;
} Manifest;

// Parses a manifest string into a linked list of ManifestEntry.
// Format: "package:version:size\npackage:version:size\n..."
// Returns NULL on parse error.
Manifest *parse_manifest(const char *data);

// Frees all memory associated with a manifest.
void free_manifest(Manifest *m);

// Finds an entry by package name. Returns NULL if not found.
ManifestEntry *find_entry(const Manifest *m, const char *package_name);

#endif
