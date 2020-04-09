/* GTP specific RTNetlink helper functions */

/* (C) 2014 by sysmocom - s.f.m.c. GmbH
 * (C) 2017 by Pablo Neira Ayuso <pablo@gnumonks.org>
 * Author: Yao-Wen Chang <yaowenowo@gmail.com>
 *
 * All Rights Reserved
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as
 * published by the Free Software Foundation; either version 2 of the
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
#include <string.h>
#include <time.h>

#include <libmnl/libmnl.h>
#include <net/if.h>
#include <linux/if_link.h>
#include <linux/rtnetlink.h>
#include <linux/gtp5g.h>
#include <linux/if_link.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <time.h>

#include <libgtp5gnl/gtp5gnl.h>

#include "internal.h"

static struct nlmsghdr *
gtp_put_nlmsg(char *buf, uint16_t type, uint16_t nl_flags, uint32_t seq)
{
	struct nlmsghdr *nlh;

	nlh = mnl_nlmsg_put_header(buf);
	nlh->nlmsg_type	= type;
	nlh->nlmsg_flags = NLM_F_REQUEST | nl_flags;
	nlh->nlmsg_seq = seq;

	return nlh;
}

static struct mnl_socket *rtnl_open(void)
{
	struct mnl_socket *nl;

	nl = mnl_socket_open(NETLINK_ROUTE);
	if (nl == NULL) {
		perror("mnl_socket_open");
		return NULL;
	}

	if (mnl_socket_bind(nl, 0, MNL_SOCKET_AUTOPID) < 0) {
		perror("mnl_socket_bind");
		goto err;
	}

	return nl;
err:
	mnl_socket_close(nl);
	return NULL;
}

static int rtnl_talk(struct mnl_socket *nl, struct nlmsghdr *nlh)
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	int ret;

	ret = mnl_socket_sendto(nl, nlh, nlh->nlmsg_len);
	if (ret < 0)
		return ret;

	ret = mnl_socket_recvfrom(nl, buf, sizeof(buf));
	if (ret < 0)
		return ret;

	return mnl_cb_run(buf, ret, nlh->nlmsg_seq, mnl_socket_get_portid(nl),
			  NULL, NULL);
}

static int gtp_dev_talk(struct nlmsghdr *nlh, uint32_t seq)
{
	struct mnl_socket *nl;
	int ret;

	nl = rtnl_open();
	if (nl == NULL)
		return -1;

	ret = rtnl_talk(nl, nlh);

	mnl_socket_close(nl);
	return ret;
}

static int _gtp_dev_create(int dest_ns, const char *gtp_ifname, int fd,
			   enum ifla_gtp5g_role role)
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	struct nlmsghdr *nlh;
	struct ifinfomsg *ifm;
	unsigned int seq = time(NULL);
	struct nlattr *nest, *nest2;

	nlh = gtp_put_nlmsg(buf, RTM_NEWLINK,
			    NLM_F_CREATE | NLM_F_EXCL | NLM_F_ACK, seq);
	ifm = mnl_nlmsg_put_extra_header(nlh, sizeof(*ifm));
	ifm->ifi_family = AF_INET;
	ifm->ifi_change |= IFF_UP;
	ifm->ifi_flags |= IFF_UP;

	if (dest_ns >= 0)
		mnl_attr_put_u32(nlh, IFLA_NET_NS_FD, dest_ns);
	mnl_attr_put_str(nlh, IFLA_IFNAME, gtp_ifname);
	nest = mnl_attr_nest_start(nlh, IFLA_LINKINFO);
	mnl_attr_put_str(nlh, IFLA_INFO_KIND, "gtp5g");
	nest2 = mnl_attr_nest_start(nlh, IFLA_INFO_DATA);
	mnl_attr_put_u32(nlh, IFLA_GTP5G_FD1, fd);
	mnl_attr_put_u32(nlh, IFLA_GTP5G_PDR_HASHSIZE, 131072);
	if (role != GTP5G_ROLE_UPF)
		mnl_attr_put_u32(nlh, IFLA_GTP5G_ROLE, role);
	mnl_attr_nest_end(nlh, nest2);
	mnl_attr_nest_end(nlh, nest);

	return gtp_dev_talk(nlh, seq);
}

int gtp_dev_create(int dest_ns, const char *gtp_ifname, int fd)
{
	return _gtp_dev_create(dest_ns, gtp_ifname, fd, GTP5G_ROLE_UPF);
}
EXPORT_SYMBOL(gtp_dev_create);

int gtp_dev_create_ran(int dest_ns, const char *gtp_ifname, int fd)
{
	return _gtp_dev_create(dest_ns, gtp_ifname, fd, GTP5G_ROLE_RAN);
}
EXPORT_SYMBOL(gtp_dev_create_ran);

int gtp_dev_destroy(const char *gtp_ifname)
{
	char buf[MNL_SOCKET_BUFFER_SIZE];
	struct nlmsghdr *nlh;
	struct ifinfomsg *ifm;
	unsigned int seq = time(NULL);

	nlh = gtp_put_nlmsg(buf, RTM_DELLINK, NLM_F_ACK, seq);
	ifm = mnl_nlmsg_put_extra_header(nlh, sizeof(*ifm));
	ifm->ifi_family = AF_INET;
	ifm->ifi_change |= IFF_UP;
	ifm->ifi_flags &= ~IFF_UP;
	ifm->ifi_index = if_nametoindex(gtp_ifname);

	return gtp_dev_talk(nlh, seq);
}
EXPORT_SYMBOL(gtp_dev_destroy);

int gtp_dev_config(const char *ifname, struct in_addr *dst, uint32_t prefix)
{
	struct mnl_socket *nl;
	char buf[MNL_SOCKET_BUFFER_SIZE];
	struct nlmsghdr *nlh;
	struct rtmsg *rtm;
	int iface, ret;

	iface = if_nametoindex(ifname);
	if (iface == 0) {
		perror("if_nametoindex");
		return -1;
	}

	nlh = mnl_nlmsg_put_header(buf);
	nlh->nlmsg_type	= RTM_NEWROUTE;
	nlh->nlmsg_flags = NLM_F_REQUEST | NLM_F_CREATE | NLM_F_ACK;
	nlh->nlmsg_seq = time(NULL);

	rtm = mnl_nlmsg_put_extra_header(nlh, sizeof(struct rtmsg));
	rtm->rtm_family = AF_INET;
	rtm->rtm_dst_len = prefix;
	rtm->rtm_src_len = 0;
	rtm->rtm_tos = 0;
	rtm->rtm_protocol = RTPROT_STATIC;
	rtm->rtm_table = RT_TABLE_MAIN;
	rtm->rtm_type = RTN_UNICAST;
	rtm->rtm_scope = RT_SCOPE_UNIVERSE;
	rtm->rtm_flags = 0;

	mnl_attr_put_u32(nlh, RTA_DST, dst->s_addr);
	mnl_attr_put_u32(nlh, RTA_OIF, iface);

	nl = rtnl_open();
	if (nl == NULL)
		return -1;

	ret = rtnl_talk(nl, nlh);

	mnl_socket_close(nl);

	return ret;
}
EXPORT_SYMBOL(gtp_dev_config);
