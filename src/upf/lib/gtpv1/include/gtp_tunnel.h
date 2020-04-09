#ifndef __GTP_TUNNEL_H__
#define __GTP_TUNNEL_H__

#include "utlt_debug.h"

#include "libgtp5gnl/gtp5gnl.h"

Status GtpTunnelAddPdr(const char *ifname, struct gtp5g_pdr *pdr);
Status GtpTunnelModPdr(const char *ifname, struct gtp5g_pdr *pdr);
Status GtpTunnelDelPdr(const char *ifname, uint16_t id);
struct gtp5g_pdr *GtpTunnelFindPdrById(const char *ifname, uint16_t id);

Status GtpTunnelAddFar(const char *ifname, struct gtp5g_far *far);
Status GtpTunnelModFar(const char *ifname, struct gtp5g_far *far);
Status GtpTunnelDelFar(const char *ifname, uint32_t id);
struct gtp5g_far *GtpTunnelFindFarById(const char *ifname, uint32_t id);

#endif /* __GTP_TUNNEL_H__ */