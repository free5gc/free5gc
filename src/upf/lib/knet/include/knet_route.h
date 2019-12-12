#ifndef __ROUTE_H__
#define __ROUTE_H__

#include <sys/socket.h>
#include <net/if.h>
#include <netinet/in.h>
#include <linux/rtnetlink.h>

#include "utlt_debug.h"
#include "utlt_list.h"

#define RT_TABLE_ANY    -1

typedef struct {
    int available;
    union {
        struct in_addr  v4;
        struct in6_addr v6;
    } addr;
} InAddr46;

// Use INET6_ADDRSTRLEN for both v4 and v6 ifnames
typedef struct {
    ListNode            node;

    int                 tableId;
    // family: protocol family, e.g. AF_INET, AF_INET6
    int                 family;
    unsigned char       tos;
    /* rtProtocol: routing protocol
     *  RTPROT_UNSPEC    = 0
     *  RTPROT_REDIRECT  = 1
     *  RTPROT_KERNEL    = 2  : route installed by kernel
     *  RTPROT_BOOT      = 3  : route installed during boot
     *  RTPROT_STATIC    = 4  : route installed by administrator
     *  RTPROT_DHCP      = 16 : DHCP client
     *
     * Values >= RTPROT_STATIC are not interpreted by kernel, they are
     * just user-defined.
     */
    unsigned char       rtProtocol;
    /* scope:
     *     RT_SCOPE_UNIVERSE = 0   : everywhere in the universe
     *
     *      ... user defined values ...
     *
     *     RT_SCOPE_SITE    = 200
     *     RT_SCOPE_LINK    = 253 : destination attached to link
     *     RT_SCOPE_HOST    = 254 : local address
     *     RT_SCOPE_NOWHERE = 255 : not existing destination
     */
    unsigned char       scope;
    /* type:
     *    RTN_UNSPEC        = 0
     *    RTN_UNICAST       = 1
     *    RTN_LOCAL         = 2
     *    RTN_BROADCAST     = 3
     *    RTN_ANYCAST       = 4
     *    RTN_MULTICAST     = 5
     *    RTN_BLACKHOLE     = 6
     *    RTN_UNREACHABLE   = 7
     *    RTN_PROHIBIT      = 8
     *    RTN_THROW         = 9
     *    RTN_NAT           = 10
     *    RTN_XRESOLVE      = 11
     *    __RTN_MAX         = 12
     */
    unsigned char       type;
    /* rtmFlags:
     *    RTM_F_NOTIFY     = 0x100: notify user of route change
     *    RTM_F_CLONED     = 0x200: this route is cloned
     *    RTM_F_EQUALIZE   = 0x400: Multipath equalizer: NI
     *    RTM_F_PREFIX     = 0x800: Prefix addresses
     */
    unsigned int        rtmFlags;

    InAddr46            dstIp;
    // dstCIDR: eg. 24 or 32 for IPv4
    int                 dstCIDR;
    InAddr46            gateway;
    char                outIfName[IF_NAMESIZE];
    uint32_t            priority;

    /**
     * perfSrcIp: Preferred source address in cases where more than 
     *            one source address could be used.
     */ 
    InAddr46            perfSrcIp;
} RouteEntryNode;

/**
 * @param  family: AF_INET or AF_INET6
 * @param  table:  id
 *    RT_TABLE_UNSPEC    = 0
 *    
 *    ... user defined values ...
 *    
 *    RT_TABLE_COMPAT    = 252
 *    RT_TABLE_DEFAULT   = 253
 *    RT_TABLE_MAIN      = 254
 *    RT_TABLE_LOCAL     = 255
 *    RT_TABLE_MAX       = 0xFFFFFFFF
 *    RT_TABLE_ANY
 * 
 * @return List of RouteEntry or NULL for error
 */
ListNode *KnetGetRoutes(int family, int table);

void KnetRtListFree(ListNode* rtEntries);

/**
 * @param *gatewayAddr ip string or NULL
 */
Status KnetAddRoute(const char *ifname, const char *dstIP, uint32_t cidrPrefix, 
                 const char *gatewayAddr, uint32_t priority);

/**
 * @param *gatewayAddr ip string or NULL
 */
Status KnetDelRoute(const char *ifname, const char *dstIP, uint32_t cidrPrefix, 
                 const char *gatewayAddr, uint32_t priority);

/**
 * @param  *routeEntries: List of RouteEntry
 */
Status KnetPrintRoutes(ListNode *routeEntries);

const char *inet6Ntoa(struct in6_addr in6);

#endif /* __ROUTE_H__ */
