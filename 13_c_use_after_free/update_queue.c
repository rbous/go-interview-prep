#include "update_queue.h"
#include <stdlib.h>
#include <string.h>

// BUG(1): queue_drain_finished uses a node pointer after freeing it.
// BUG(2): queue_remove doesn't handle removing the head or tail correctly
//         in all cases (NULL pointer dereference on edge cases).
// BUG(3): queue_destroy doesn't free node fields before freeing nodes.

UpdateQueue *queue_create(void) {
    UpdateQueue *q = calloc(1, sizeof(UpdateQueue));
    return q;
}

void queue_destroy(UpdateQueue *q) {
    if (!q) return;

    UpdateNode *curr = q->head;
    while (curr) {
        UpdateNode *next = curr->next;
        free(curr);
        curr = next;
    }
    free(q);
}

UpdateNode *queue_push(UpdateQueue *q, const char *package_name, const char *version) {
    UpdateNode *node = calloc(1, sizeof(UpdateNode));
    node->package_name = strdup(package_name);
    node->version = strdup(version);
    node->status = STATUS_PENDING;

    if (q->tail) {
        q->tail->next = node;
        node->prev = q->tail;
    } else {
        q->head = node;
    }
    q->tail = node;
    q->count++;
    return node;
}

void queue_remove(UpdateQueue *q, UpdateNode *node) {
    if (!q || !node) return;

    // BUG: doesn't check if prev/next are NULL (crashes when removing head or tail)
    node->prev->next = node->next;
    node->next->prev = node->prev;
    q->count--;

    free(node->package_name);
    free(node->version);
    free(node);
}

int queue_start_all(UpdateQueue *q) {
    if (!q) return 0;
    int count = 0;
    UpdateNode *curr = q->head;
    while (curr) {
        if (curr->status == STATUS_PENDING) {
            curr->status = STATUS_DOWNLOADING;
            count++;
        }
        curr = curr->next;
    }
    return count;
}

int queue_drain_finished(UpdateQueue *q) {
    if (!q) return 0;
    int count = 0;
    UpdateNode *curr = q->head;
    while (curr) {
        // BUG: we free curr via queue_remove, then access curr->next
        if (curr->status == STATUS_COMPLETE || curr->status == STATUS_FAILED) {
            queue_remove(q, curr);
            count++;
        }
        curr = curr->next;
    }
    return count;
}
