/* GTP specific Generic Netlink helper functions */

/* (C) 2014 by sysmocom - s.f.m.c. GmbH
 * Author: Pablo Neira Ayuso <pablo@gnumonks.org>
 *
 * All Rights Reserved
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as
 * published by the Free Software Foundation; either version 2.1 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <inttypes.h>

#include <libmnl/libmnl.h>
#include <linux/genetlink.h>

#include <libgtpnl/gtp.h>
#include <libgtpnl/gtpnl.h>

#include <net/if.h>
#include <linux/gtp.h>
#include <linux/if_link.h>

#include "internal.h"

static void gtp_build_payload(struct nlmsghdr *nlh, struct gtp_tunnel *t)
{
	mnl_attr_put_u32(nlh, GTPA_VERSION, t->gtp_version);
	if (t->ifns >= 0)
		mnl_attr_put_u32(nlh, GTPA_NET_NS_FD, t->ifns);
	mnl_attr_put_u32(nlh, GTPA_LINK, t->ifidx);
	if (t->sgsn_addr.s_addr)
		mnl_attr_put_u32(nlh, GTPA_PEER_ADDRESS, t->sgsn_addr.s_addr);
	if (t->ms_addr.s_addr)
		mnl_attr_put_u32(nlh, GTPA_MS_ADDRESS, t->ms_addr.s_addr);
	if (t->gtp_version == GTP_V0) {
		mnl_attr_put_u64(nlh, GTPA_TID, t->u.v0.tid);
		mnl_attr_put_u16(nlh, GTPA_FLOW, t->u.v0.flowid);
	} else if (t->gtp_version == GTP_V1) {
		mnl_attr_put_u32(nlh, GTPA_I_TEI, t->u.v1.i_tei);
		mnl_attr_put_u32(nlh, GTPA_O_TEI, t->u.v1.o_tei);
	}
}

int gtp_add_tunnel(int genl_id, struct mnl_socket *nl, struct gtp_tunnel *t)
{
	struct nlmsghdr *nlh;
	char buf[MNL_SOCKET_BUFFER_SIZE];
	uint32_t seq = time(NULL);

	if (t->gtp_version > GTP_V1) {
		fprintf(stderr, "wrong GTP version %u, use v0 or v1\n",
			t->gtp_version);
		return -1;
	}

	nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_EXCL | NLM_F_ACK, ++seq,
				   GTP_CMD_NEWPDP);
	gtp_build_payload(nlh, t);

	if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0)
		perror("genl_socket_talk");

	return 0;
}
EXPORT_SYMBOL(gtp_add_tunnel);

int gtp_del_tunnel(int genl_id, struct mnl_socket *nl, struct gtp_tunnel *t)
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	struct nlmsghdr *nlh;
	uint32_t seq = time(NULL);

	nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_ACK, ++seq,
				   GTP_CMD_DELPDP);
	gtp_build_payload(nlh, t);

	if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0)
		perror("genl_socket_talk");

	return 0;
}
EXPORT_SYMBOL(gtp_del_tunnel);

struct gtp_pdp {
	uint32_t	version;
	union {
		struct {
			uint64_t tid;
		} v0;
		struct {
			uint32_t i_tei;
			uint32_t o_tei;
		} v1;
	} u;
	struct in_addr	sgsn_addr;
	struct in_addr	ms_addr;
};

static int genl_gtp_validate_cb(const struct nlattr *attr, void *data)
{
	const struct nlattr **tb = data;
	int type = mnl_attr_get_type(attr);

	if (mnl_attr_type_valid(attr, GTPA_MAX) < 0)
		return MNL_CB_OK;

	switch(type) {
	case GTPA_TID:
		if (mnl_attr_validate(attr, MNL_TYPE_U64) < 0) {
			perror("mnl_attr_validate");
			return MNL_CB_ERROR;
		}
		break;
	case GTPA_O_TEI:
	case GTPA_I_TEI:
	case GTPA_PEER_ADDRESS:
	case GTPA_MS_ADDRESS:
	case GTPA_VERSION:
		if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0) {
			perror("mnl_attr_validate");
			return MNL_CB_ERROR;
		}
		break;
	default:
		break;
	}
	tb[type] = attr;
	return MNL_CB_OK;
}

static int genl_gtp_attr_cb(const struct nlmsghdr *nlh, void *data)
{
	struct nlattr *tb[GTPA_MAX + 1] = {};
	char buf[INET_ADDRSTRLEN];
	struct gtp_pdp pdp = {};
	struct genlmsghdr *genl;

	mnl_attr_parse(nlh, sizeof(*genl), genl_gtp_validate_cb, tb);
	if (tb[GTPA_TID])
		pdp.u.v0.tid = mnl_attr_get_u64(tb[GTPA_TID]);
	if (tb[GTPA_I_TEI])
		pdp.u.v1.i_tei = mnl_attr_get_u32(tb[GTPA_I_TEI]);
	if (tb[GTPA_O_TEI])
		pdp.u.v1.o_tei = mnl_attr_get_u32(tb[GTPA_O_TEI]);
	if (tb[GTPA_PEER_ADDRESS]) {
		pdp.sgsn_addr.s_addr =
			mnl_attr_get_u32(tb[GTPA_PEER_ADDRESS]);
	}
	if (tb[GTPA_MS_ADDRESS]) {
		pdp.ms_addr.s_addr = mnl_attr_get_u32(tb[GTPA_MS_ADDRESS]);
	}
	if (tb[GTPA_VERSION]) {
		pdp.version = mnl_attr_get_u32(tb[GTPA_VERSION]);
	}

	printf("version %u ", pdp.version);
	if (pdp.version == GTP_V0) {
		inet_ntop(AF_INET, &pdp.ms_addr, buf, sizeof(buf));
		printf("tid %"PRIu64" ms_addr %s ",
		       pdp.u.v0.tid, buf);
	} else if (pdp.version == GTP_V1) {
		inet_ntop(AF_INET, &pdp.ms_addr, buf, sizeof(buf));
		printf("tei %u/%u ms_addr %s ", pdp.u.v1.i_tei,
		       pdp.u.v1.o_tei, buf);
	}
	inet_ntop(AF_INET, &pdp.sgsn_addr, buf, sizeof(buf));
	printf("sgsn_addr %s\n", buf);

	return MNL_CB_OK;
}

int gtp_list_tunnel(int genl_id, struct mnl_socket *nl)
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	struct nlmsghdr *nlh;
	uint32_t seq = time(NULL);

	nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_DUMP, 0,
				   GTP_CMD_GETPDP);

	if (genl_socket_talk(nl, nlh, seq, genl_gtp_attr_cb, NULL) < 0) {
		perror("genl_socket_talk");
		return 0;
	}

	return 0;
}
EXPORT_SYMBOL(gtp_list_tunnel);
