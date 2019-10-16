#ifndef __3GPP_TYPE_H__
#define __3GPP_TYPE_H__

#include <stdint.h>
#include "utlt_network.h"
#include "utlt_lib.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

// IP protocol, Ref: https://stackoverflow.com/questions/1076714/max-length-for-client-ip-address/1076755
#define IPV4_LEN                    4
#define IPV6_LEN                    16
#define IPV4V6_LEN                  20

#define S1AP_SCTP_PORT              36412
#define NGAP_SCTP_PORT              38412
#define GTPV2_C_UDP_PORT            2123
#define GTPV1_U_UDP_PORT            2152
#define PFCP_UDP_PORT               8805

#define SCTP_S1AP_PPID              18
#define SCTP_X2AP_PPID              27
#define SCTP_NGAP_PPID              60

#define MAX_NUM_OF_RAN              128
#define MAX_NUM_OF_5G_UE            128

#define MAX_POOL_OF_5G_UE           (MAX_NUM_OF_RAN * MAX_NUM_OF_5G_UE)

#define MAX_NUM_OF_ENB              128
#define MAX_NUM_OF_UE               128
#define MAX_NUM_OF_SESS             4
#define MAX_NUM_OF_BEARER           4
#define MAX_NUM_OF_TUNNEL           3   /* Num of Tunnel per Bearer */
#define MAX_NUM_OF_PF               16  /* Num of Packet Filter per Bearer */

#define MAX_POOL_OF_UE              (MAX_NUM_OF_ENB * MAX_NUM_OF_UE)
#define MAX_POOL_OF_SESS            (MAX_POOL_OF_UE * MAX_NUM_OF_SESS)
#define MAX_POOL_OF_BEARER          (MAX_POOL_OF_SESS * MAX_NUM_OF_BEARER)
#define MAX_POOL_OF_TUNNEL          (MAX_POOL_OF_BEARER * MAX_NUM_OF_TUNNEL)
#define MAX_POOL_OF_PF              (MAX_POOL_OF_BEARER * MAX_NUM_OF_PF)
#define MAX_POOL_OF_DIAMETER_SESS   (MAX_POOL_OF_UE * MAX_NUM_OF_SESS)

#define MAX_NUM_OF_HOSTNAME         16
#define MAX_NUM_OF_PCC_RULE         8 
#define MAX_NUM_OF_FLOW             8   /* Num of Flow per PCC Rule */
#define MAX_NUM_OF_PACKET_FILTER    16  /* Num of Packet Filter per Bearer */

#define MAX_SDU_LEN                 8192
#define PLMN_ID_LEN                 3
#define AMF_REGION_ID_LEN           2
#define AMF_SET_ID_LEN              4
#define AMF_POINTER_LEN             4
#define SST_LEN                     1
#define SD_LEN                      3

#define BCD_TO_BUFFER_LEN(x)        (((x)+1)/2)
#define MAX_IMSI_BCD_LEN            15
#define MAX_IMSI_LEN                BCD_TO_BUFFER_LEN(MAX_IMSI_BCD_LEN)

#define RAND_LEN                    16
#define AUTN_LEN                    16
#define AUTS_LEN                    14
#define MAX_RES_LEN                 16

#define MAX_APN_LEN                 100
#define MAX_PCO_LEN                 251

#define MAX_SBI_MESSAGE_LEN         4096

typedef struct _plmnId_t {
    uint8_t octet[3];
} __attribute__ ((packed)) plmnId_t;

typedef struct _Tai {
    plmnId_t planId;
    uint16_t tac;
} Tai;

typedef struct _Paa {
/* 8.34 PDN Type  */
#define GTP_PDN_TYPE_IPV4                               1
#define GTP_PDN_TYPE_IPV6                               2
#define GTP_PDN_TYPE_IPV4V6                             3
#define GTP_PDN_TYPE_NON_IP                             4
#define SBI_PDN_TYPE_IPV4                               1
#define SBI_PDN_TYPE_IPV6                               2
#define SBI_PDN_TYPE_IPV4V6                             3
#define SBI_PDN_TYPE_NON_IP                             4
    ENDIAN2(uint8_t spare:5;,
            uint8_t pdnType:3;)
        union {
            /* GTP_PDN_TYPE_IPV4 */
            struct in_addr      addr4;
            /* GTP_PDN_TYPE_IPV6 */
            struct in6_addr     addr6;
            /* GTP_PDN_TYPE_BOTH */
            struct {
                struct in_addr  addr4;
                struct in6_addr addr6;
            } dualStack;
        };
} Paa;

typedef struct _Bitrate {
    /* bits per seconds */
    uint64_t    downlink;
    uint64_t    uplink;
} Bitrate;

typedef struct _Qos {
#define PDN_QCI_1                                       1
#define PDN_QCI_2                                       2
#define PDN_QCI_3                                       3
#define PDN_QCI_4                                       4
#define PDN_QCI_5                                       5
#define PDN_QCI_6                                       6
#define PDN_QCI_7                                       7
#define PDN_QCI_8                                       8
#define PDN_QCI_9                                       9
    uint8_t     qci;

    struct {
        uint8_t     priorityLevel;
        uint8_t     preEmptionCapability;
        uint8_t     preEmptionVulnerability;
    } arp;

    Bitrate     mbr;
    Bitrate     gbr;
} Qos;

typedef struct _Pdn {
    uint32_t    contextId;
    uint8_t     apn[MAX_APN_LEN+1];
#define HSS_PDN_TYPE_IPV4                       0
#define HSS_PDN_TYPE_IPV6                       1
#define HSS_PDN_TYPE_IPV4V6                     2
#define HSS_PDN_TYPE_IPV4_OR_IPV6               3
    int8_t      pdnType;

    Qos         qos;
    Bitrate     ambr;

    Paa         paa;
    Ip          upfIp;
} Pdn;

void SetPlmnId(plmnId_t *plmnId, uint16_t mcc, uint16_t mnc, uint8_t mncLen);
uint16_t GetMcc(plmnId_t *plmnId);
uint16_t GetMnc(plmnId_t *plmnId);
uint16_t GetMncLen(plmnId_t *plmnId);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __3GPP_TYPE_H__ */
