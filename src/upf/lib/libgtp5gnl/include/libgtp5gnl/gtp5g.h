#ifndef _LIBGTP5G_H_
#define _LIBGTP5G_H_

#include <stdint.h>
#include <netinet/ip.h>

struct gtp5g_dev;
struct gtp5g_pdr;
struct gtp5g_far;

struct gtp5g_dev *gtp5g_dev_alloc(void);
struct gtp5g_pdr *gtp5g_pdr_alloc(void);
struct gtp5g_far *gtp5g_far_alloc(void);

void gtp5g_dev_free(struct gtp5g_dev *dev);
void gtp5g_pdr_free(struct gtp5g_pdr *pdr);
void gtp5g_far_free(struct gtp5g_far *far);

void gtp5g_dev_set_ifns(struct gtp5g_dev *dev, int ifns);
void gtp5g_dev_set_ifidx(struct gtp5g_dev *dev, uint32_t ifidx);

void gtp5g_pdr_set_id(struct gtp5g_pdr *pdr, uint16_t id);
void gtp5g_pdr_set_precedence(struct gtp5g_pdr *pdr, uint32_t precedence);
void gtp5g_pdr_set_far_id(struct gtp5g_pdr *pdr, uint32_t far_id);
void gtp5g_pdr_set_outer_header_removal(struct gtp5g_pdr *pdr, uint8_t outer_hdr_removal);
void gtp5g_pdr_set_ue_addr_ipv4(struct gtp5g_pdr *pdr, struct in_addr *ue_addr_ipv4);
void gtp5g_pdr_set_local_f_teid(struct gtp5g_pdr *pdr, uint32_t teid, struct in_addr *gtpu_addr_ipv4);

/* Not in 3GPP spec, just used for routing */
void gtp5g_pdr_set_role_addr_ipv4(struct gtp5g_pdr *pdr, struct in_addr *role_addr_ipv4);
/* Not in 3GPP spec, just used for buffering */
void gtp5g_pdr_set_unix_sock_path(struct gtp5g_pdr *pdr, const char *unix_sock_path);

void gtp5g_pdr_set_sdf_filter_description(struct gtp5g_pdr *pdr, const char *rule_str);
void gtp5g_pdr_set_tos_traffic_class(struct gtp5g_pdr *pdr, uint16_t tos_traffic_class);
void gtp5g_pdr_set_security_param_idx(struct gtp5g_pdr *pdr, uint32_t security_param_idx);
void gtp5g_pdr_set_flow_label(struct gtp5g_pdr *pdr, uint32_t flow_label);
void gtp5g_pdr_set_sdf_filter_id(struct gtp5g_pdr *pdr, uint32_t id);

uint16_t *gtp5g_pdr_get_id(struct gtp5g_pdr *pdr);
uint32_t *gtp5g_pdr_get_precedence(struct gtp5g_pdr *pdr);
uint32_t *gtp5g_pdr_get_far_id(struct gtp5g_pdr *pdr);
uint8_t  *gtp5g_pdr_get_outer_header_removal(struct gtp5g_pdr *pdr);
struct in_addr *gtp5g_pdr_get_ue_addr_ipv4(struct gtp5g_pdr *pdr);
uint32_t *gtp5g_pdr_get_local_f_teid_teid(struct gtp5g_pdr *pdr);
struct in_addr *gtp5g_pdr_get_local_f_teid_gtpu_addr_ipv4(struct gtp5g_pdr *pdr);

// TODO: Get sdf_filter value

// TODO: Get FAR value
void gtp5g_far_set_id(struct gtp5g_far *far, uint32_t id);
void gtp5g_far_set_apply_action(struct gtp5g_far *far, uint8_t apply_action);
void gtp5g_far_set_outer_header_creation(struct gtp5g_far *far, uint16_t desp, uint32_t teid, struct in_addr *peer_addr_ipv4, uint16_t port);

uint32_t *gtp5g_far_get_id(struct gtp5g_far *far);
uint8_t *gtp5g_far_get_apply_action(struct gtp5g_far *far);
uint16_t *gtp5g_far_get_outer_header_creation_description(struct gtp5g_far *far);
uint32_t *gtp5g_far_get_outer_header_creation_teid(struct gtp5g_far *far);
struct in_addr *gtp5g_far_get_outer_header_creation_peer_addr_ipv4(struct gtp5g_far *far);
uint16_t *gtp5g_far_get_outer_header_creation_port(struct gtp5g_far *far);

int *gtp5g_far_get_related_pdr_num(struct gtp5g_far *far);
uint16_t *gtp5g_far_get_related_pdr_list(struct gtp5g_far *far);

#endif
