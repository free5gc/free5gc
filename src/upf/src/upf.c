#include "upf.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include "utlt_lib.h"
#include "utlt_debug.h"
#include "utlt_event.h"
#include "utlt_network.h"
#include "upf_context.h"
#include "n4/n4_dispatcher.h"

static Status parseArgs(int argc, char *argv[]);
static Status checkPermission();
static void eventConsumer();

char configPath[MAX_FILE_PATH_STRLEN] = "./config/upfcfg.yaml";

int main(int argc, char *argv[]) {
    Status status, returnStatus = STATUS_OK;

    UTLT_Assert(parseArgs(argc, argv) == STATUS_OK, return STATUS_ERROR, 
                "Error parsing args");

    if (checkPermission() != STATUS_OK) {
        return STATUS_ERROR;
    }

    status = UpfInit(configPath);
    UTLT_Assert(status == STATUS_OK, returnStatus = STATUS_ERROR,
                "UPF failed to initialize");

    if (status == STATUS_OK) {
        eventConsumer();
    }

    status = UpfTerminate();
    UTLT_Assert(status == STATUS_OK, returnStatus = STATUS_ERROR,
                "UPF terminate error");

    return returnStatus;
}

static Status parseArgs(int argc, char *argv[]) {
    int opt;

    while ((opt = getopt(argc, argv, "f:h")) != -1) {
        switch (opt) {
            case 'f':
                strcpy(configPath, optarg);
                break;

            case 'h':
                printf("Usage: %s [-f CONFIG_PATH]", argv[0]);
                exit(0);
            
            case '?':
                UTLT_Error("Illigal option: %c", optopt); 
                return STATUS_ERROR;
        }
    }

    return STATUS_OK;
}

static Status checkPermission() {
    if (geteuid() != 0) {
        UTLT_Error("Please run UPF as root in order to enable route management "
                   "and communication with gtp5g kernel module.");
        return STATUS_ERROR;
    }
    return STATUS_OK;
}


static void eventConsumer() {
    Status status;
    Event event;

    while (1) {
        status = EventRecv(Self()->eventQ, &event);
        if (status != STATUS_OK) {
            if (status == STATUS_EAGAIN) {
                continue;
            } else {
                UTLT_Assert(0, break, "Event receive fail");
            }
        }

        UpfDispatcher(&event);
    }
}
