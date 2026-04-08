// DO NOT MODIFY THIS FILE.
#include "progress.h"
#include <stdio.h>
#include <string.h>

static int tests_run = 0;
static int tests_passed = 0;

#define TEST(name) do { \
    tests_run++; \
    printf("  TEST %-40s ", #name); \
    if (test_##name()) { tests_passed++; printf("PASS\n"); } \
    else { printf("FAIL\n"); } \
} while(0)

static int test_create_and_free(void) {
    const char *names[] = { "curl", "wget" };
    uint64_t sizes[] = { 1000, 2000 };
    ProgressTracker *t = tracker_create(2, names, sizes);
    if (!t) return 0;
    int ok = (t->count == 2)
        && (t->entries[0].total_bytes == 1000)
        && (t->entries[1].total_bytes == 2000);
    tracker_free(t);
    return ok;
}

static int test_download_all_completes(void) {
    const char *names[] = { "pkg-a", "pkg-b", "pkg-c" };
    uint64_t sizes[] = { 10000, 20000, 15000 };
    ProgressTracker *t = tracker_create(3, names, sizes);

    download_all(t);

    int ok = 1;
    for (int i = 0; i < 3; i++) {
        if (t->entries[i].bytes_downloaded != t->entries[i].total_bytes) {
            printf("entry %d: got %llu, want %llu  ",
                i, t->entries[i].bytes_downloaded, t->entries[i].total_bytes);
            ok = 0;
        }
        if (!t->entries[i].complete) {
            printf("entry %d not marked complete  ", i);
            ok = 0;
        }
    }

    if (t->num_complete != 3) {
        printf("num_complete=%d, want 3  ", t->num_complete);
        ok = 0;
    }

    tracker_free(t);
    return ok;
}

static int test_progress_100_after_download(void) {
    const char *names[] = { "a", "b" };
    uint64_t sizes[] = { 5000, 5000 };
    ProgressTracker *t = tracker_create(2, names, sizes);

    download_all(t);
    int pct = overall_progress(t);

    tracker_free(t);
    return pct == 100;
}

static int test_progress_0_before_download(void) {
    const char *names[] = { "a", "b" };
    uint64_t sizes[] = { 5000, 5000 };
    ProgressTracker *t = tracker_create(2, names, sizes);

    int pct = overall_progress(t);
    tracker_free(t);
    return pct == 0;
}

static int test_many_concurrent_downloads(void) {
    const int n = 20;
    const char *names[20];
    uint64_t sizes[20];
    char name_bufs[20][16];

    for (int i = 0; i < n; i++) {
        snprintf(name_bufs[i], 16, "pkg-%d", i);
        names[i] = name_bufs[i];
        sizes[i] = 10000 + i * 1000;
    }

    ProgressTracker *t = tracker_create(n, names, sizes);
    download_all(t);

    int ok = (t->num_complete == n);
    for (int i = 0; i < n; i++) {
        if (t->entries[i].bytes_downloaded != t->entries[i].total_bytes) {
            ok = 0;
            break;
        }
    }

    tracker_free(t);
    return ok;
}

static int test_each_package_gets_own_download(void) {
    // This catches the loop variable capture bug.
    // If all threads share &i, some packages won't be downloaded.
    const char *names[] = { "alpha", "bravo", "charlie", "delta", "echo" };
    uint64_t sizes[] = { 1000, 2000, 3000, 4000, 5000 };
    ProgressTracker *t = tracker_create(5, names, sizes);

    download_all(t);

    int ok = 1;
    for (int i = 0; i < 5; i++) {
        if (t->entries[i].bytes_downloaded != t->entries[i].total_bytes) {
            printf("entry %d (%s): got %llu, want %llu  ",
                i, names[i], t->entries[i].bytes_downloaded, t->entries[i].total_bytes);
            ok = 0;
        }
    }

    tracker_free(t);
    return ok;
}

int main(void) {
    printf("Running progress tracker tests:\n");

    TEST(create_and_free);
    TEST(download_all_completes);
    TEST(progress_100_after_download);
    TEST(progress_0_before_download);
    TEST(many_concurrent_downloads);
    TEST(each_package_gets_own_download);

    printf("\n%d/%d tests passed\n", tests_passed, tests_run);
    return tests_passed == tests_run ? 0 : 1;
}
