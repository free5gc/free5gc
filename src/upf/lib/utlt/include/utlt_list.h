#ifndef __LIST_H__
#define __LIST_H__

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _ListNode {
    struct _ListNode *prev;
    struct _ListNode *next;
} ListNode;

// Initialize head of list
#define ListInit(__namePtr) do {\
    (__namePtr)->prev = NULL; \
    (__namePtr)->next = NULL; \
} while (0)

#define ListFirst(__namePtr) ((void *)((__namePtr)->next))

#define ListLast(__namePtr) ((void *)((__namePtr)->prev))

#define ListPrev(__nodePtr) ((void *)(((ListNode *)(__nodePtr))->prev))

#define ListNext(__nodePtr) ((void *)(((ListNode *)(__nodePtr))->next))

#define ListIsEmpty(__namePtr) ((__namePtr)->next == NULL)

#define ListAppend(__namePtr, __newPtr) do { \
    ListNode *iter = (__namePtr); \
    while (iter->next) \
        iter = ListNext(iter); \
    ((ListNode *)(__newPtr))->prev = iter; \
    ((ListNode *)(__newPtr))->next = NULL; \
    iter->next = (ListNode *)(__newPtr); \
} while (0)

#define ListInsertToPrev(__namePtr, __nodePtr, __newPtr) do { \
    ((ListNode *)(__newPtr))->prev = ((ListNode *)(__nodePtr))->prev; \
    ((ListNode *)(__newPtr))->next = (ListNode *)(__nodePtr); \
    if (((ListNode *)(__nodePtr))->prev) \
        ((ListNode *)(__nodePtr))->prev->next = (ListNode *)(__newPtr); \
    else \
        (__namePtr)->next = (ListNode *)(__newPtr); \
    ((ListNode *)(__nodePtr))->prev = ((ListNode *)(__newPtr)); \
} while (0)

#define ListRemove(__namePtr, __nodePtr) do { \
    ListNode *iter = (__namePtr); \
    while (iter->next) { \
        if (iter->next == (ListNode *)(__nodePtr)) { \
            iter->next = ((ListNode *)(__nodePtr))->next; \
            if (iter->next) \
                iter->next->prev = (ListNode *)iter; \
            break; \
        } \
        iter = iter->next; \
    } \
} while (0)

typedef int (*ListNodeCompare)(ListNode *pnode1, ListNode *pnode2);

#define ListInsertSorted(__namePtr, __newPtr, __cmpCallback) do { \
    ListNodeCompare callbackPtr = (ListNodeCompare)__cmpCallback; \
    ListNode *iter = ListFirst(__namePtr); \
    while (iter) { \
        if ((*callbackPtr)((ListNode *)(__newPtr), iter) < 0) { \
            ListInsertToPrev(__namePtr, iter, __newPtr); \
            break; \
        } \
        iter = ListNext(iter); \
    } \
    if (iter == NULL) \
        ListAppend(__namePtr, __newPtr); \
} while (0)

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __LIST_H__ */