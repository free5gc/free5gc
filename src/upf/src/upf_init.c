#include "upf.h"

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <signal.h>
#include <errno.h>

#include "utlt_lib.h"
#include "utlt_debug.h"
#include "utlt_buff.h"
#include "utlt_thread.h"
#include "utlt_timer.h"
#include "utlt_network.h"
#include "gtp_path.h"
#include "upf_context.h"
#include "upf_config.h"
#include "up/up_gtp_path.h"
#include "n4/n4_pfcp_path.h"
#include "pfcp_xact.h"

static Status SignalRegister();
static void SignalHandler(int sigval);

void PacketReceiverThread(ThreadID id, void *data);

// TODO : Add other init in this part for UPF
Status UpfInit(char *configPath) {
    Status status;

    // Resolve config path
    UTLT_Assert(GetAbsPath(configPath) == STATUS_OK, 
                return STATUS_ERROR, "Invalid config path: %s", configPath);
    UTLT_Info("Config: %s", configPath);

    status = UtltLibInit();
    if (status != STATUS_OK) return status;

    status = SignalRegister();
    if (status != STATUS_OK) return status;

    status = UpfContextInit();
    if (status != STATUS_OK) return status;

    status = UpfLoadConfigFile(configPath);
    if (status != STATUS_OK) return status;

    status = UpfConfigParse();
    if (status != STATUS_OK) return status;

    Self()->epfd = EpollCreate();
    if (Self()->epfd < 0) return STATUS_ERROR;

    Self()->eventQ = EventQueueCreate(O_RDWR | O_NONBLOCK);
    if (Self()->eventQ <= 0) return STATUS_ERROR;

    status = ThreadCreate(&Self()->pktRecvThread, PacketReceiverThread, NULL);
    if (status != STATUS_OK) return status;

    status = GTPv1ServerInit();
    if (status != STATUS_OK) return status;

    status = PfcpServerInit();
    if (status != STATUS_OK) return status;

    status = PfcpXactInit(&Self()->timerServiceList, UPF_EVENT_N4_T3_RESPONSE, UPF_EVENT_N4_T3_HOLDING); // init pfcp xact context
    if (status != STATUS_OK) return status;

    status = GtpRouteInit();
    if (status != STATUS_OK) return status;

    UTLT_Info("UPF initialized");

    return STATUS_OK;
}

Status UpfTerminate() {
    Status status = STATUS_OK;

    UTLT_Info("Terminating UPF...");

    UTLT_Assert(GtpRouteTerminate() == STATUS_OK, status |= STATUS_ERROR,
                "GTP routes removal failed");

    UTLT_Assert(PfcpServerTerminate() == STATUS_OK, status |= STATUS_ERROR,
                "PFCP server terminate failed");

    UTLT_Assert(GTPv1ServerTerminate() == STATUS_OK, status |= STATUS_ERROR,
                "GTPv1 server terminate failed");

    UTLT_Assert(ThreadDelete(Self()->pktRecvThread) == STATUS_OK, status |= STATUS_ERROR,
                "Thread PacketReceiverThread delete failed");

    UTLT_Assert(EventQueueDelete(Self()->eventQ) == STATUS_OK, status |= STATUS_ERROR,
                "Event queue delete failed");

    close(Self()->epfd);

    UTLT_Assert(UpfContextTerminate() == STATUS_OK, status |= STATUS_ERROR,
                "UPF context terminate failed");

    UTLT_Assert(UtltLibTerminate() == STATUS_OK, status |= STATUS_ERROR,
                "UPF library terminate failed");

    if (status == STATUS_OK)
        UTLT_Info("UPF terminated");
    else
        UTLT_Info("UPF failed to terminate");

    return status;
}

Status UtltLibInit() {
    Status status;

    status = BufblkPoolInit();
    if (status != STATUS_OK) return status;

    status = ThreadInit();
    if (status != STATUS_OK) return status;

    status = TimerPoolInit();
    if (status != STATUS_OK) return status;

    status = SockPoolInit();
    if (status != STATUS_OK) return status;

    status = Gtpv1DevPoolInit();
    if (status != STATUS_OK) return status;

    status = TimerPoolInit();
    if (status != STATUS_OK) return status;

    return STATUS_OK;
} 

Status UtltLibTerminate() {
    Status status = STATUS_OK;

    UTLT_Assert(Gtpv1DevPoolFinal() == STATUS_OK, status |= STATUS_ERROR,
                "Gtpv1Dev pool terminate failed");

    UTLT_Assert(SockPoolFinal() == STATUS_OK, status |= STATUS_ERROR,
                "Socket pool terminate failed");

    UTLT_Assert(TimerFinal() == STATUS_OK, status |= STATUS_ERROR,
                "Timer pool terminate failed");

    UTLT_Assert(ThreadFinal() == STATUS_OK, status |= STATUS_ERROR,
                "Thread terminate failed");

    UTLT_Assert(BufblkPoolFinal() == STATUS_OK, status |= STATUS_ERROR,
                "Bufblk pool terminate failed");

    return status;
}

static Status SignalRegister() {
    signal(SIGINT, SignalHandler);
    signal(SIGTERM, SignalHandler);

    return STATUS_OK;
}

static void SignalHandler(int sigval) {
    switch(sigval) {
        case SIGINT :
            UTLT_Assert(UpfTerminate() == STATUS_OK, , "Handle Ctrl-C fail");
            break;
        case SIGTERM :
            UTLT_Assert(UpfTerminate() == STATUS_OK, , "Handle Ctrl-C fail");
            break;
        default :
            break;
    }
    exit(0);
}

void PacketReceiverThread(ThreadID id, void *data) {
    Status status;

    int nfds;
    Sock *sockPtr;
    struct epoll_event events[MAX_NUM_OF_EVENT];
    while (!ThreadStop()) {
        nfds = EpollWait(Self()->epfd, events, 1);
        UTLT_Assert(nfds >= 0, , "Epoll Wait error : %s", strerror(errno));

        for (int i = 0; i < nfds; i++) {
            sockPtr = events[i].data.ptr;
            status = sockPtr->handler(sockPtr, sockPtr->data);
            // TODO : Log may show which socket
            UTLT_Assert(status == STATUS_OK, , "Error handling UP socket");
        }
    }

    sem_post(((Thread *)id)->semaphore);
    UTLT_Info("Packet receiver thread terminated");

    return;
}
