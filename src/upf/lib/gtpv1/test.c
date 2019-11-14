#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include "utlt_debug.h"
#include "utlt_network.h"
#include "gtp_link.h"


Status testGtpv1_1() {
    printf("[Sudo Testing] testGtpv1 1\n");
    Status status;

    Gtpv1TunDevNode node = {
        .ip = "0.0.0.0",
        .ifname = "GtPvLtEsT",
    };
    
    status = GtpLinkCreate(&node);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "GtpLinkCreate ifdev[%s] fail", node.ifname);
    printf("[Sudo Testing] Use `ip a` to check the device named %s is existed\n",
           node.ifname);
    sleep(5);

    status = GtpLinkFree(&node);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "GtpLinkFree ifdev[%s] fail", node.ifname);
    printf("[Sudo Testing] Use `ip a` to check the device named %s is not existed\n",
           node.ifname);
    sleep(5);

    return STATUS_OK;
}

Status testGtpv1_2() {
    printf("[Sudo Testing] testGtpv1 2\n");
    Status status;

    Gtpv1TunDevNode node = {
        .ip = "0.0.0.0",
        .ifname = "gTpDeV1CeTeSt",
    };
    int iteid = 87, oteid = 78;
    const char destIP[] = "192.168.87.87";
    const char tunIP[] = "172.168.31.87";

    status = GtpLinkCreate(&node);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR,
                "GtpLinkCreate ifdev[%s] fail", node.ifname);
    printf("[Sudo Testing] Use `ip a` to check the device named %s is existed\n",
           node.ifname);
    sleep(5);

    status = GtpTunnelAdd(node.ifname, iteid, oteid, destIP, tunIP);
    if (status != STATUS_OK) {
        UTLT_Assert(0, , "GtpTunnelAdd ifname[%s], iteid[%d], oteid[%d], Dest IP[%s], Tunnel IP[%s] fail", 
                    node.ifname, iteid, oteid, destIP, tunIP);
        status = GtpLinkFree(&node);
        UTLT_Assert(status == STATUS_OK, , "GtpLinkFree ifdev[%s] fail", node.ifname);
        return STATUS_ERROR;
    }

    printf("[Sudo Testing] After add tunnel, GTP Device '%s' Tunnel List\n", node.ifname);
    status = GtpTunnelList();
    if (status != STATUS_OK) {
        UTLT_Assert(0, , "GtpTunnelList ifname[%s]. fail", node.ifname);
        status = GtpLinkFree(&node);
        UTLT_Assert(status == STATUS_OK, , "GtpLinkFree ifdev[%s] fail", node.ifname);
        return STATUS_ERROR;
    }
    
    status = GtpTunnelDel(node.ifname, iteid);
    if (status != STATUS_OK) {
        UTLT_Assert(0, , "GtpTunnelDel ifname[%s], iteid[%d] fail", node.ifname, iteid);
        status = GtpLinkFree(&node);
        UTLT_Assert(status == STATUS_OK, , "GtpLinkFree ifdev[%s] fail", node.ifname);
        return STATUS_ERROR;
    }

    printf("[Sudo Testing] After delete tunnel, GTP Device '%s' Tunnel List\n", node.ifname);
    status = GtpTunnelList();
    if (status != STATUS_OK) {
        UTLT_Assert(0, , "GtpTunnelList ifname[%s]. fail", node.ifname);
        status = GtpLinkFree(&node);
        UTLT_Assert(status == STATUS_OK, , "GtpLinkFree ifdev[%s] fail", node.ifname);
        return STATUS_ERROR;
    }

    status = GtpLinkFree(&node);
    UTLT_Assert(status == STATUS_OK, return STATUS_ERROR, "GtpLinkFree ifdev[%s] fail", node.ifname);

    return STATUS_OK;
}

int main() {
    fprintf(stderr, "**********************************************************************************************\n");
    fprintf(stderr, "[Warning] This test will need privilege user to run because it will modify the network setting\n");
    fprintf(stderr, "[Warning] Please open the new terminal tab to do this test\n");
    fprintf(stderr, "**********************************************************************************************\n");

    SockPoolInit();
    Gtpv1DevPoolInit();

    testGtpv1_1();
    testGtpv1_2();

    return 0;
}
