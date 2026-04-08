#include "progress.h"
#include <pthread.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

// BUG(1): Multiple threads update tracker->num_complete without synchronization.
// BUG(2): bytes_downloaded is read and written by different threads without
//         any locking (data race in overall_progress during download).
// BUG(3): The thread function accesses an index variable that may change
//         before the thread reads it (classic loop variable capture bug).
//
// Fix all three bugs. You may use mutexes, atomics, or any other approach.

ProgressTracker *tracker_create(int count, const char **package_names, const uint64_t *sizes) {
    ProgressTracker *t = calloc(1, sizeof(ProgressTracker));
    t->count = count;
    t->entries = calloc(count, sizeof(DownloadProgress));

    for (int i = 0; i < count; i++) {
        strncpy(t->entries[i].package_name, package_names[i], 63);
        t->entries[i].total_bytes = sizes[i];
    }

    return t;
}

void tracker_free(ProgressTracker *t) {
    if (!t) return;
    free(t->entries);
    free(t);
}

static void *download_worker(void *arg) {
    // arg points to a struct { ProgressTracker*, int index }
    // but we're being sloppy...
    void **args = (void **)arg;
    ProgressTracker *t = (ProgressTracker *)args[0];
    int idx = *(int *)args[1];

    DownloadProgress *p = &t->entries[idx];
    uint64_t chunk = p->total_bytes / 10;
    if (chunk == 0) chunk = 1;

    while (p->bytes_downloaded < p->total_bytes) {
        uint64_t remaining = p->total_bytes - p->bytes_downloaded;
        uint64_t to_add = (remaining < chunk) ? remaining : chunk;
        p->bytes_downloaded += to_add;
        usleep(1000); // simulate network delay
    }

    p->complete = 1;
    t->num_complete++;

    return NULL;
}

void download_all(ProgressTracker *t) {
    pthread_t *threads = malloc(sizeof(pthread_t) * t->count);

    for (int i = 0; i < t->count; i++) {
        // BUG: &i is shared across all threads — classic loop variable capture
        void *args[2] = { t, &i };
        pthread_create(&threads[i], NULL, download_worker, args);
    }

    for (int i = 0; i < t->count; i++) {
        pthread_join(threads[i], NULL);
    }

    free(threads);
}

int overall_progress(const ProgressTracker *t) {
    if (!t || t->count == 0) return 0;

    uint64_t total = 0;
    uint64_t downloaded = 0;
    for (int i = 0; i < t->count; i++) {
        total += t->entries[i].total_bytes;
        downloaded += t->entries[i].bytes_downloaded;
    }

    if (total == 0) return 100;
    return (int)((downloaded * 100) / total);
}
