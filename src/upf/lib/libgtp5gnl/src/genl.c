/* Generic Netlink helper functions */

/* (C) 2014 by sysmocom - s.f.m.c. GmbH
 * (C) 2017 by Pablo Neira Ayuso <pablo@gnumonks.org>
 * Author: Yao-Wen Chang <yaowenowo@gmail.com>
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
#include <unistd.h>
#include <time.h>

#include <libmnl/libmnl.h>
#include <linux/genetlink.h>

#include <libgtp5gnl/gtp5gnl.h>

#include "internal.h"

struct nlmsghdr *
genl_nlmsg_build_hdr(char *buf, uint16_t type, uint16_t flags, uint32_t seq,
		     uint8_t cmd)
{
	struct nlmsghdr *nlh;
	struct genlmsghdr *genl;

	nlh = mnl_nlmsg_put_header(buf);
	nlh->nlmsg_type = type;
	nlh->nlmsg_flags = NLM_F_REQUEST | flags;
	nlh->nlmsg_seq = seq;

	genl = mnl_nlmsg_put_extra_header(nlh, sizeof(struct genlmsghdr));
	genl->cmd = cmd;
	genl->version = 0;

	return nlh;
}
EXPORT_SYMBOL(genl_nlmsg_build_hdr);

static int genl_ctrl_validate_cb(const struct nlattr *attr, void *data)
{
	const struct nlattr **tb = data;
	int type = mnl_attr_get_type(attr);

	if (mnl_attr_type_valid(attr, CTRL_ATTR_MAX) < 0)
		return MNL_CB_OK;

	switch(type) {
	case CTRL_ATTR_FAMILY_ID:
		if (mnl_attr_validate(attr, MNL_TYPE_U16) < 0) {
			perror("mnl_attr_validate");
			return MNL_CB_ERROR;
		}
		break;
	}
	tb[type] = attr;
	return MNL_CB_OK;
}

static int genl_ctrl_cb(const struct nlmsghdr *nlh, void *data)
{
	struct nlattr *tb[CTRL_ATTR_MAX + 1] = {};
	struct genlmsghdr *genl = mnl_nlmsg_get_payload(nlh);
	int32_t *genl_id = data;

	mnl_attr_parse(nlh, sizeof(*genl), genl_ctrl_validate_cb, tb);
	if (tb[CTRL_ATTR_FAMILY_ID])
		*genl_id = mnl_attr_get_u16(tb[CTRL_ATTR_FAMILY_ID]);
	else
		*genl_id = -1;

	return MNL_CB_OK;
}

struct mnl_socket *genl_socket_open(void)
{
	struct mnl_socket *nl;

	nl = mnl_socket_open(NETLINK_GENERIC);
	if (nl == NULL) {
		perror("mnl_socket_open");
		return NULL;
	}

	if (mnl_socket_bind(nl, 0, MNL_SOCKET_AUTOPID) < 0) {
		perror("mnl_socket_bind");
		return NULL;
	}

	return nl;
}
EXPORT_SYMBOL(genl_socket_open);

void genl_socket_close(struct mnl_socket *nl)
{
	mnl_socket_close(nl);
}
EXPORT_SYMBOL(genl_socket_close);

int genl_socket_talk(struct mnl_socket *nl, struct nlmsghdr *nlh, uint32_t seq,
		     int (*cb)(const struct nlmsghdr *nlh, void *data),
		     void *data)
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	int ret;

	if (mnl_socket_sendto(nl, nlh, nlh->nlmsg_len) < 0) {
		perror("mnl_socket_send");
		return -1;
	}

	ret = mnl_socket_recvfrom(nl, buf, sizeof(buf));
	while (ret > 0) {
		ret = mnl_cb_run(buf, ret, seq, mnl_socket_get_portid(nl),
				 cb, data);
		if (ret <= 0)
			break;
		ret = mnl_socket_recvfrom(nl, buf, sizeof(buf));
	}

	return ret;
}
EXPORT_SYMBOL(genl_socket_talk);

static struct nlmsghdr *
genl_nlmsg_build_lookup(char *buf, const char *subsys_name)
{
	struct nlmsghdr *nlh;
	struct genlmsghdr *genl;

	nlh = mnl_nlmsg_put_header(buf);
	nlh->nlmsg_type	= GENL_ID_CTRL;
	nlh->nlmsg_flags = NLM_F_REQUEST | NLM_F_ACK;
	nlh->nlmsg_seq = time(NULL);

	genl = mnl_nlmsg_put_extra_header(nlh, sizeof(struct genlmsghdr));
	genl->cmd = CTRL_CMD_GETFAMILY;
	genl->version = 1;

	mnl_attr_put_u16(nlh, CTRL_ATTR_FAMILY_ID, GENL_ID_CTRL);
	mnl_attr_put_strz(nlh, CTRL_ATTR_FAMILY_NAME, subsys_name);

	return nlh;
}

int genl_lookup_family(struct mnl_socket *nl, const char *family)
{
	int32_t genl_id;
	char buf[MNL_SOCKET_BUFFER_SIZE];
	struct nlmsghdr *nlh = genl_nlmsg_build_lookup(buf, family);
	int err;

	err = genl_socket_talk(nl, nlh, nlh->nlmsg_seq, genl_ctrl_cb, &genl_id);
	if (err < 0)
		return -1;

	return genl_id;
}
EXPORT_SYMBOL(genl_lookup_family);
