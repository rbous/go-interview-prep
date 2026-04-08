#ifndef VERSION_PARSER_H
#define VERSION_PARSER_H

#define VERSION_FIELD_MAX 16

typedef struct {
    char major[VERSION_FIELD_MAX];
    char minor[VERSION_FIELD_MAX];
    char patch[VERSION_FIELD_MAX];
    char label[VERSION_FIELD_MAX];  // e.g. "beta", "rc1"
} FirmwareVersion;

// Parses a version string like "1.22.3-beta" into a FirmwareVersion struct.
// Returns 0 on success, -1 on error (e.g., malformed input).
int parse_version(const char *version_str, FirmwareVersion *out);

// Compares two version strings. Returns:
//   -1 if a < b, 0 if a == b, 1 if a > b
// Only compares major.minor.patch numerically; ignores label.
int compare_versions(const char *a, const char *b);

// Formats a FirmwareVersion back into a string like "1.22.3-beta".
// Writes into `buf` of size `buf_size`. Returns 0 on success, -1 if buf too small.
int format_version(const FirmwareVersion *v, char *buf, int buf_size);

#endif
