#ifndef __UTLT_TEST_H__
#define __UTLT_TEST_H__

#include "utlt_debug.h"
#include "utlt_list.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct {
    const char *name;
    Status (*FuncPtr)(void *data);
    void *data;
} TestCase;

typedef struct {
    ListNode node;
    TestCase test;
} TestNode;

typedef struct {
    int totalCase;
    int finishCase;
    
    ListNode node;
} TestContext;

Status TestInit();
Status TestTerminate();

Status TestAdd(TestCase *testCase);
Status TestAddList(TestCase *testList, int size);

Status TestRun();

int TestCaseArrayFindByName(TestCase *array, int size, const char *targetName);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __UTLT_TEST_H__ */
