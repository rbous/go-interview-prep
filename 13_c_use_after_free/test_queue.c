// DO NOT MODIFY THIS FILE.
#include "update_queue.h"
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

static int test_create_destroy(void) {
    UpdateQueue *q = queue_create();
    if (!q) return 0;
    queue_destroy(q);
    return 1;
}

static int test_push_and_count(void) {
    UpdateQueue *q = queue_create();
    queue_push(q, "curl", "8.0");
    queue_push(q, "wget", "1.21");
    queue_push(q, "git", "2.40");
    int ok = (q->count == 3);
    queue_destroy(q);
    return ok;
}

static int test_remove_middle(void) {
    UpdateQueue *q = queue_create();
    queue_push(q, "a", "1.0");
    UpdateNode *mid = queue_push(q, "b", "2.0");
    queue_push(q, "c", "3.0");

    queue_remove(q, mid);
    int ok = (q->count == 2)
        && strcmp(q->head->package_name, "a") == 0
        && strcmp(q->tail->package_name, "c") == 0
        && q->head->next == q->tail
        && q->tail->prev == q->head;
    queue_destroy(q);
    return ok;
}

static int test_remove_head(void) {
    UpdateQueue *q = queue_create();
    UpdateNode *head = queue_push(q, "a", "1.0");
    queue_push(q, "b", "2.0");

    queue_remove(q, head);
    int ok = (q->count == 1)
        && strcmp(q->head->package_name, "b") == 0
        && q->head->prev == NULL;
    queue_destroy(q);
    return ok;
}

static int test_remove_tail(void) {
    UpdateQueue *q = queue_create();
    queue_push(q, "a", "1.0");
    UpdateNode *tail = queue_push(q, "b", "2.0");

    queue_remove(q, tail);
    int ok = (q->count == 1)
        && strcmp(q->tail->package_name, "a") == 0
        && q->tail->next == NULL;
    queue_destroy(q);
    return ok;
}

static int test_remove_only(void) {
    UpdateQueue *q = queue_create();
    UpdateNode *only = queue_push(q, "a", "1.0");

    queue_remove(q, only);
    int ok = (q->count == 0) && (q->head == NULL) && (q->tail == NULL);
    queue_destroy(q);
    return ok;
}

static int test_start_all(void) {
    UpdateQueue *q = queue_create();
    queue_push(q, "a", "1.0");
    queue_push(q, "b", "2.0");

    int started = queue_start_all(q);
    int ok = (started == 2)
        && q->head->status == STATUS_DOWNLOADING
        && q->tail->status == STATUS_DOWNLOADING;
    queue_destroy(q);
    return ok;
}

static int test_drain_finished(void) {
    UpdateQueue *q = queue_create();
    UpdateNode *a = queue_push(q, "a", "1.0");
    queue_push(q, "b", "2.0");
    UpdateNode *c = queue_push(q, "c", "3.0");

    a->status = STATUS_COMPLETE;
    c->status = STATUS_FAILED;

    int drained = queue_drain_finished(q);
    int ok = (drained == 2)
        && (q->count == 1)
        && strcmp(q->head->package_name, "b") == 0;
    queue_destroy(q);
    return ok;
}

static int test_drain_all_finished(void) {
    UpdateQueue *q = queue_create();
    UpdateNode *a = queue_push(q, "a", "1.0");
    UpdateNode *b = queue_push(q, "b", "2.0");

    a->status = STATUS_COMPLETE;
    b->status = STATUS_COMPLETE;

    int drained = queue_drain_finished(q);
    int ok = (drained == 2) && (q->count == 0)
        && (q->head == NULL) && (q->tail == NULL);
    queue_destroy(q);
    return ok;
}

static int test_destroy_frees_fields(void) {
    // ASAN will catch leaks here
    UpdateQueue *q = queue_create();
    queue_push(q, "longpackagename", "1.0.0-beta.1");
    queue_push(q, "another-package", "2.5.3-rc2");
    queue_destroy(q);
    return 1;
}

int main(void) {
    printf("Running update_queue tests:\n");

    TEST(create_destroy);
    TEST(push_and_count);
    TEST(remove_middle);
    TEST(remove_head);
    TEST(remove_tail);
    TEST(remove_only);
    TEST(start_all);
    TEST(drain_finished);
    TEST(drain_all_finished);
    TEST(destroy_frees_fields);

    printf("\n%d/%d tests passed\n", tests_passed, tests_run);
    return tests_passed == tests_run ? 0 : 1;
}
