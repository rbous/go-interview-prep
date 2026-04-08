// DO NOT MODIFY THIS FILE.
#include "manifest.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>

static int tests_run = 0;
static int tests_passed = 0;

#define TEST(name) do { \
    tests_run++; \
    printf("  TEST %-40s ", #name); \
    if (test_##name()) { tests_passed++; printf("PASS\n"); } \
    else { printf("FAIL\n"); } \
} while(0)

static int test_parse_basic(void) {
    Manifest *m = parse_manifest("curl:8.0.0:1024\nwget:1.21.0:2048");
    if (!m) return 0;
    int ok = m->count == 2;
    free_manifest(m);
    return ok;
}

static int test_parse_single(void) {
    Manifest *m = parse_manifest("openssl:3.1.0:4096");
    if (!m) return 0;
    int ok = (m->count == 1)
        && strcmp(m->head->package_name, "openssl") == 0
        && strcmp(m->head->version, "3.1.0") == 0
        && m->head->size_bytes == 4096;
    free_manifest(m);
    return ok;
}

static int test_parse_null(void) {
    return parse_manifest(NULL) == NULL;
}

// Malformed line mid-manifest: should return NULL and not leak
// the entries parsed before the error.
// Run with -fsanitize=address to detect leaks.
static int test_parse_error_no_leak(void) {
    Manifest *m = parse_manifest("curl:8.0.0:1024\nBADLINE\nwget:1.21.0:2048");
    // Should return NULL due to malformed "BADLINE"
    return m == NULL;
}

static int test_free_null(void) {
    free_manifest(NULL);  // Should not crash
    return 1;
}

// Allocate, free, and let ASAN check for leaks.
static int test_no_leak_on_free(void) {
    Manifest *m = parse_manifest("a:1.0:100\nb:2.0:200\nc:3.0:300");
    if (!m) return 0;
    free_manifest(m);
    return 1;  // ASAN will catch leaks
}

static int test_find_entry(void) {
    Manifest *m = parse_manifest("curl:8.0.0:1024\nwget:1.21.0:2048");
    if (!m) return 0;

    ManifestEntry *e = find_entry(m, "curl");
    int ok = e && strcmp(e->version, "8.0.0") == 0;

    free_manifest(m);
    return ok;
}

static int test_find_entry_not_found(void) {
    Manifest *m = parse_manifest("curl:8.0.0:1024");
    if (!m) return 0;

    ManifestEntry *e = find_entry(m, "nonexistent");
    int ok = (e == NULL);

    free_manifest(m);
    return ok;
}

static int test_find_entry_null_manifest(void) {
    // Should not crash on NULL manifest
    ManifestEntry *e = find_entry(NULL, "curl");
    return e == NULL;
}

int main(void) {
    printf("Running manifest tests:\n");

    TEST(parse_basic);
    TEST(parse_single);
    TEST(parse_null);
    TEST(parse_error_no_leak);
    TEST(free_null);
    TEST(no_leak_on_free);
    TEST(find_entry);
    TEST(find_entry_not_found);
    TEST(find_entry_null_manifest);

    printf("\n%d/%d tests passed\n", tests_passed, tests_run);
    return tests_passed == tests_run ? 0 : 1;
}
