#ifndef UPDATE_QUEUE_H
#define UPDATE_QUEUE_H

typedef enum {
    STATUS_PENDING,
    STATUS_DOWNLOADING,
    STATUS_APPLYING,
    STATUS_COMPLETE,
    STATUS_FAILED
} UpdateStatus;

typedef struct UpdateNode {
    char *package_name;
    char *version;
    UpdateStatus status;
    struct UpdateNode *next;
    struct UpdateNode *prev;
} UpdateNode;

typedef struct {
    UpdateNode *head;
    UpdateNode *tail;
    int count;
} UpdateQueue;

UpdateQueue *queue_create(void);
void queue_destroy(UpdateQueue *q);

// Adds an update to the end of the queue. Returns the new node.
UpdateNode *queue_push(UpdateQueue *q, const char *package_name, const char *version);

// Removes a node from the queue and frees it.
void queue_remove(UpdateQueue *q, UpdateNode *node);

// Advances all PENDING items to DOWNLOADING status.
// Returns the number of items advanced.
int queue_start_all(UpdateQueue *q);

// Removes all completed or failed items from the queue.
// Returns the number removed.
int queue_drain_finished(UpdateQueue *q);

#endif
