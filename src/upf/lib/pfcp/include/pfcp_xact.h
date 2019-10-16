#ifndef __PFCP_XACT_H__
#define __PFCP_XACT_H__

#include <stdint.h>

#include "utlt_index.h"
#include "utlt_list.h"
#include "utlt_index.h"
#include "utlt_buff.h"

#include "pfcp_message.h"

#include "pfcp_node.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

typedef struct _PfcpXact {
    ListNode    node;
    uint32_t    index;

#define PFCP_LOCAL_ORIGINATOR  0
#define PFCP_REMOTE_ORIGINATOR 1
    uint8_t     origin;
    uint32_t    transactionId;
    PfcpNode    *gnode;

    int         step;               // 1: Init, 2: Trigger, 3: Trigger Reply
    struct {
        uint8_t type;
        Bufblk  *bufBlk;
    } seq[3];

    TimerBlkID  timerResponse;
    uint8_t     responseReCount;
    TimerBlkID  timerHolding;
    uint8_t     holdingReCount;

    struct _PfcpXact    *associatedXact;

#define PfcpXactStoreSession(xact, session) \
    do { \
        UTLT_Assert((xact), break, "xact error"); \
        UTLT_Assert((session), break, "session error"); \
        ((xact)->session) = (session); \
    } while(0)

    void        *session;
    void        *gtpXact;
    Bufblk      *gtpBuf;
} PfcpXact;

typedef enum {
    PFCP_XACT_UNKNOWN_STAGE,
    PFCP_XACT_INITIAL_STAGE,
    PFCP_XACT_INTERMEDIATE_STAGE,
    PFCP_XACT_FINAL_STAGE,
} PfcpXactStage;

Status PfcpXactInit(TimerList *timerList, uint32_t responseEvent, uint32_t holdingEvent);
Status PfcpXactTerminate(void);
PfcpXact *PfcpXactLocalCreate(PfcpNode *gnode, PfcpHeader *header, Bufblk *bufBlk);
PfcpXact *PfcpXactRemoteCreate(PfcpNode *gnode, uint32_t sqn);
void PfcpXactDeleteAll(PfcpNode *gnode);
Status PfcpXactUpdateTx(PfcpXact *xact, PfcpHeader *header, Bufblk *bufBlk);
Status PfcpXactUpdateRx(PfcpXact *xact, uint8_t type);
Status PfcpXactCommit(PfcpXact *xact);
Status PfcpXactTimeout(uint32_t index, uint32_t event, uint8_t *type);
Status PfcpXactReceive(PfcpNode *gnode, PfcpHeader *header, PfcpXact **xact);
PfcpXact *PfcpXactFind(uint32_t index);
//static PfcpXactStage PfcpXactGetStage(uint8_t type, uint32_t transactionId);
PfcpXact *PfcpXactFindByTransactionId(PfcpNode *gnode, uint8_t type, uint32_t transactionId);
//void PfcpXactAssociate(PfcpXact *xact1, PfcpXact *xact2);
//void PfcpXactDeassociate(PfcpXact *xact1, PfcpXact *xact2);
//static Status PfcpXactDelete(PfcpXact *xact);


#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif