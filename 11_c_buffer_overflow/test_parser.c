// DO NOT MODIFY THIS FILE.
#include "version_parser.h"
#include <stdio.h>
#include <string.h>
#include <assert.h>

static int tests_run = 0;
static int tests_passed = 0;

#define TEST(name) do { \
    tests_run++; \
    printf("  TEST %-40s ", #name); \
    if (test_##name()) { tests_passed++; printf("PASS\n"); } \
    else { printf("FAIL\n"); } \
} while(0)

static int test_basic_parse(void) {
    FirmwareVersion v;
    int rc = parse_version("1.22.3", &v);
    return rc == 0
        && strcmp(v.major, "1") == 0
        && strcmp(v.minor, "22") == 0
        && strcmp(v.patch, "3") == 0
        && strlen(v.label) == 0;
}

static int test_parse_with_label(void) {
    FirmwareVersion v;
    int rc = parse_version("2.0.1-beta", &v);
    return rc == 0
        && strcmp(v.major, "2") == 0
        && strcmp(v.minor, "0") == 0
        && strcmp(v.patch, "1") == 0
        && strcmp(v.label, "beta") == 0;
}

static int test_parse_null_input(void) {
    FirmwareVersion v;
    return parse_version(NULL, &v) == -1;
}

static int test_parse_null_output(void) {
    return parse_version("1.0.0", NULL) == -1;
}

// This tests the buffer overflow bug: a field longer than VERSION_FIELD_MAX
// should return an error, not overflow.
static int test_parse_long_field(void) {
    // "12345678901234567890" is 20 chars, larger than VERSION_FIELD_MAX (16)
    int rc = parse_version("12345678901234567890.0.0", NULL);
    // Should return -1 (error) or at least not crash.
    // If the program reaches here without crashing, the overflow is fixed.
    // We accept -1 as the correct return for oversized fields.
    (void)rc;

    FirmwareVersion v;
    rc = parse_version("1.0.0-areallylonglabelthatexceedsbuffersize", &v);
    return rc == -1;  // Should reject oversized label
}

static int test_compare_equal(void) {
    return compare_versions("1.2.3", "1.2.3") == 0;
}

static int test_compare_major(void) {
    return compare_versions("1.0.0", "2.0.0") == -1
        && compare_versions("2.0.0", "1.0.0") == 1;
}

static int test_compare_minor(void) {
    return compare_versions("1.1.0", "1.2.0") == -1;
}

static int test_compare_patch(void) {
    return compare_versions("1.0.1", "1.0.2") == -1;
}

static int test_compare_null(void) {
    // Should not crash on NULL input. Return 0 or any value, just don't segfault.
    compare_versions(NULL, "1.0.0");
    compare_versions("1.0.0", NULL);
    compare_versions(NULL, NULL);
    return 1;  // If we get here without crashing, pass.
}

static int test_format_basic(void) {
    FirmwareVersion v = { "1", "22", "3", "" };
    char buf[64];
    int rc = format_version(&v, buf, sizeof(buf));
    return rc == 0 && strcmp(buf, "1.22.3") == 0;
}

static int test_format_with_label(void) {
    FirmwareVersion v = { "2", "0", "1", "rc1" };
    char buf[64];
    int rc = format_version(&v, buf, sizeof(buf));
    return rc == 0 && strcmp(buf, "2.0.1-rc1") == 0;
}

// Buffer too small for output — should return -1, not overflow.
static int test_format_small_buffer(void) {
    FirmwareVersion v = { "10", "20", "30", "beta" };
    char buf[5];  // Way too small for "10.20.30-beta"
    int rc = format_version(&v, buf, sizeof(buf));
    return rc == -1;
}

int main(void) {
    printf("Running version_parser tests:\n");

    TEST(basic_parse);
    TEST(parse_with_label);
    TEST(parse_null_input);
    TEST(parse_null_output);
    TEST(parse_long_field);
    TEST(compare_equal);
    TEST(compare_major);
    TEST(compare_minor);
    TEST(compare_patch);
    TEST(compare_null);
    TEST(format_basic);
    TEST(format_with_label);
    TEST(format_small_buffer);

    printf("\n%d/%d tests passed\n", tests_passed, tests_run);
    return tests_passed == tests_run ? 0 : 1;
}
