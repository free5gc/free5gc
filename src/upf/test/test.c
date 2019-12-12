#include <stdio.h>
#include <string.h>

#include "test_upf.h"
#include "utlt_debug.h"

static TestCase upfTestList[] = {
    {"UPTest", UPTest, NULL},
};

// TODO : Add Alarm to prevent timeout
int main(int argc, char *argv[]) {
    UTLT_Assert(TestInit() == STATUS_OK, return -1, "TestInit fial");

    int sizeOfTestList = sizeof(upfTestList)/sizeof(TestCase);
    if (argc > 1) {
        for (int i = 1; i < argc; i++) {
            int idx = TestCaseArrayFindByName(upfTestList, sizeOfTestList, argv[i]);
            if (idx >= 0) TestAdd(&upfTestList[idx]);
        }
    } else {
        TestAddList(upfTestList, sizeOfTestList);
    }

    TestRun();

    UTLT_Assert(TestTerminate() == STATUS_OK, return -1, "TestTerminate fial");

    return 0;
}
