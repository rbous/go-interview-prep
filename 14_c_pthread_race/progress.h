#ifndef PROGRESS_H
#define PROGRESS_H

#include <stdint.h>

typedef struct {
    char     package_name[64];
    uint64_t bytes_downloaded;
    uint64_t total_bytes;
    int      complete;
} DownloadProgress;

typedef struct {
    DownloadProgress *entries;
    int count;
    int num_complete;
} ProgressTracker;

// Creates a tracker for `count` packages.
ProgressTracker *tracker_create(int count, const char **package_names, const uint64_t *sizes);

// Frees the tracker.
void tracker_free(ProgressTracker *t);

// Simulates downloading all packages concurrently using pthreads.
// Each "download" increments bytes_downloaded in chunks until it reaches total_bytes.
// Should be thread-safe.
void download_all(ProgressTracker *t);

// Returns overall progress as a percentage (0-100).
int overall_progress(const ProgressTracker *t);

#endif
