#ifndef _LIBGTP_H_
#define _LIBGTP_H_

#include <stdint.h>
#include <netinet/ip.h>

struct gtp_tunnel;

struct gtp_tunnel *gtp_tunnel_alloc(void);
void gtp_tunnel_free(struct gtp_tunnel *t);

void gtp_tunnel_set_ifns(struct gtp_tunnel *t, int ifns);
void gtp_tunnel_set_ifidx(struct gtp_tunnel *t, uint32_t ifidx);
void gtp_tunnel_set_ms_ip4(struct gtp_tunnel *t, struct in_addr *ms_addr);
void gtp_tunnel_set_sgsn_ip4(struct gtp_tunnel *t, struct in_addr *sgsn_addr);
void gtp_tunnel_set_version(struct gtp_tunnel *t, uint32_t version);
void gtp_tunnel_set_tid(struct gtp_tunnel *t, uint64_t tid);
void gtp_tunnel_set_i_tei(struct gtp_tunnel *t, uint32_t i_tei);
void gtp_tunnel_set_o_tei(struct gtp_tunnel *t, uint32_t o_tei);
void gtp_tunnel_set_flowid(struct gtp_tunnel *t, uint16_t flowid);

const int gtp_tunnel_get_ifns(struct gtp_tunnel *t);
const uint32_t gtp_tunnel_get_ifidx(struct gtp_tunnel *t);
const struct in_addr *gtp_tunnel_get_ms_ip4(struct gtp_tunnel *t);
const struct in_addr *gtp_tunnel_get_sgsn_ip4(struct gtp_tunnel *t);
int gtp_tunnel_get_version(struct gtp_tunnel *t);
uint64_t gtp_tunnel_get_tid(struct gtp_tunnel *t);
uint32_t gtp_tunnel_get_i_tei(struct gtp_tunnel *t);
uint32_t gtp_tunnel_get_o_tei(struct gtp_tunnel *t);
uint16_t gtp_tunnel_get_flowid(struct gtp_tunnel *t);

#endif
