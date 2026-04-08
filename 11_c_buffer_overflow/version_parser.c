#include "version_parser.h"
#include <string.h>
#include <stdio.h>
#include <stdlib.h>

// BUG(1): strcpy does not check destination buffer size.
//         A version field longer than VERSION_FIELD_MAX causes a buffer overflow.
// BUG(2): compare_versions doesn't handle NULL input.
// BUG(3): format_version uses sprintf without checking buffer size.

int parse_version(const char *version_str, FirmwareVersion *out) {
    if (!version_str || !out)
        return -1;

    memset(out, 0, sizeof(FirmwareVersion));

    // Copy input so we can tokenize
    char buf[256];
    strcpy(buf, version_str);

    // Split on '-' to separate label
    char *dash = strchr(buf, '-');
    if (dash) {
        *dash = '\0';
        strcpy(out->label, dash + 1);
    }

    // Split major.minor.patch on '.'
    char *token = strtok(buf, ".");
    if (!token) return -1;
    strcpy(out->major, token);

    token = strtok(NULL, ".");
    if (!token) return -1;
    strcpy(out->minor, token);

    token = strtok(NULL, ".");
    if (token) {
        strcpy(out->patch, token);
    }

    return 0;
}

int compare_versions(const char *a, const char *b) {
    FirmwareVersion va, vb;
    parse_version(a, &va);
    parse_version(b, &vb);

    int a_major = atoi(va.major);
    int b_major = atoi(vb.major);
    if (a_major != b_major) return a_major < b_major ? -1 : 1;

    int a_minor = atoi(va.minor);
    int b_minor = atoi(vb.minor);
    if (a_minor != b_minor) return a_minor < b_minor ? -1 : 1;

    int a_patch = atoi(va.patch);
    int b_patch = atoi(vb.patch);
    if (a_patch != b_patch) return a_patch < b_patch ? -1 : 1;

    return 0;
}

int format_version(const FirmwareVersion *v, char *buf, int buf_size) {
    if (!v || !buf)
        return -1;

    if (strlen(v->label) > 0) {
        sprintf(buf, "%s.%s.%s-%s", v->major, v->minor, v->patch, v->label);
    } else {
        sprintf(buf, "%s.%s.%s", v->major, v->minor, v->patch);
    }

    return 0;
}
