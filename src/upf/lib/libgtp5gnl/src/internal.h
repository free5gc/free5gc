#ifndef INTERNAL_H
#define INTERNAL_H 1

#include "config.h"
#ifdef HAVE_VISIBILITY_HIDDEN
#	define __visible	__attribute__((visibility("default")))
#	define EXPORT_SYMBOL(x)	typeof(x) (x) __visible
#else
#	define EXPORT_SYMBOL
#endif

#include <stdint.h>
#include <netinet/in.h>

struct gtp5g_dev {
    int ifns;
    uint32_t ifidx;
};

struct gtp5g_outer_header_creation {
    uint16_t desp;
    uint32_t teid; // o_teid
    struct in_addr peer_addr_ipv4;
    uint16_t port;
};

struct gtp5g_forwarding_parameter {
//    uint8_t dest_int;
//    char *network_instance;
    
    struct gtp5g_outer_header_creation *hdr_creation;
};

struct gtp5g_far {
    uint32_t id;
    uint8_t apply_action;
    
    struct gtp5g_forwarding_parameter *fwd_param;

    /* Not IEs in 3GPP Spec, for other purpose */
    int related_pdr_num;
    uint16_t *related_pdr_list;
};

struct local_f_teid {
    uint32_t teid; // i_teid
    struct in_addr gtpu_addr_ipv4;
};

struct ip_filter_rule {
    uint8_t action;                   // permit only
    uint8_t direction;                  // in/out
    uint8_t proto;                      // number or "ip" which is not used for matching
    struct in_addr src, smask;          // ip addr or "any" -> 0.0.0.0
    struct in_addr dest, dmask;         // ip addr or "assigned" -> 0.0.0.0
    int sport_num;
    uint32_t *sport_list;               // one value, range or not existed -> [0, 0]
    int dport_num;
    uint32_t *dport_list;               // one value, range or not existed -> [0, 0]
};

struct sdf_filter {
    struct ip_filter_rule *rule;
    uint16_t *tos_traffic_class;
    uint32_t *security_param_idx;
    uint32_t *flow_label;               // exactly 3 Octets
    uint32_t *bi_id;
};

struct gtp5g_pdi {
//    uint8_t src_int;
//    char *network_instance;
    struct in_addr *ue_addr_ipv4;

/* Local F-TEID */
    struct local_f_teid *f_teid;
    struct sdf_filter *sdf;
};

struct gtp5g_pdr {
    uint16_t id;
    uint32_t *precedence;
    struct gtp5g_pdi *pdi;

    uint8_t *outer_hdr_removal;
    uint32_t *far_id;

    /* Not in 3GPP spec, just used for routing */
    struct in_addr *role_addr_ipv4;

    /* Not in 3GPP spec, just used for buffering */
    char *unix_sock_path;
};

#endif
