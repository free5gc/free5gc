/* GTP specific Generic Netlink helper functions */

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
#include <string.h>
#include <unistd.h>
#include <time.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <inttypes.h>

#include <libmnl/libmnl.h>
#include <linux/genetlink.h>

#include <libgtp5gnl/gtp5g.h>
#include <libgtp5gnl/gtp5gnl.h>

#include <net/if.h>
#include <linux/gtp5g.h>
#include <linux/if_link.h>

#include "internal.h"
#include "tools.h"

static void gtp5g_build_far_payload(struct nlmsghdr *nlh, struct gtp5g_dev *dev, struct gtp5g_far *far)
{
    // Let kernel get dev easily
    if (dev->ifns >= 0)
        mnl_attr_put_u32(nlh, GTP5G_NET_NS_FD, dev->ifns);
    mnl_attr_put_u32(nlh, GTP5G_LINK, dev->ifidx);

    // Level 1 FAR
    mnl_attr_put_u32(nlh, GTP5G_FAR_ID, far->id);
    
    if (far->apply_action)
        mnl_attr_put_u8(nlh, GTP5G_FAR_APPLY_ACTION, far->apply_action);

    // Level 2 FAR : Forwarding Parameter
    struct nlattr *fwd_param_nest, *hdr_creation_nest;
    if (far->fwd_param) {
        fwd_param_nest = mnl_attr_nest_start(nlh, GTP5G_FAR_FORWARDING_PARAMETER);

        // Level 3 : Outer header Creation
        if (far->fwd_param->hdr_creation) {
            hdr_creation_nest = mnl_attr_nest_start(nlh, GTP5G_FORWARDING_PARAMETER_OUTER_HEADER_CREATION);

            mnl_attr_put_u16(nlh, GTP5G_OUTER_HEADER_CREATION_DESCRIPTION,
                             far->fwd_param->hdr_creation->desp);
            mnl_attr_put_u32(nlh, GTP5G_OUTER_HEADER_CREATION_O_TEID,
                             far->fwd_param->hdr_creation->teid);
            mnl_attr_put_u32(nlh, GTP5G_OUTER_HEADER_CREATION_PEER_ADDR_IPV4,
                             far->fwd_param->hdr_creation->peer_addr_ipv4.s_addr);
            mnl_attr_put_u16(nlh, GTP5G_OUTER_HEADER_CREATION_PORT,
                             far->fwd_param->hdr_creation->port);
            mnl_attr_nest_end(nlh, hdr_creation_nest);
        }
        mnl_attr_nest_end(nlh, fwd_param_nest);
    }
}

int gtp5g_add_far(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far)
{
    struct nlmsghdr *nlh;
    char buf[MNL_SOCKET_BUFFER_SIZE];
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return -1;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_EXCL | NLM_F_ACK, ++seq,
                               GTP5G_CMD_ADD_FAR);
    gtp5g_build_far_payload(nlh, dev, far);

    if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_add_far);

int gtp5g_mod_far(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far)
{
    struct nlmsghdr *nlh;
    char buf[MNL_SOCKET_BUFFER_SIZE];
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return -1;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_REPLACE | NLM_F_ACK, ++seq,
                               GTP5G_CMD_ADD_FAR);
    gtp5g_build_far_payload(nlh, dev, far);

    if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_mod_far);


int gtp5g_del_far(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far)
{
    char buf[MNL_SOCKET_BUFFER_SIZE];
    struct nlmsghdr *nlh;
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return -1;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_ACK, ++seq,
                               GTP5G_CMD_DEL_FAR);
    gtp5g_build_far_payload(nlh, dev, far);

    if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_del_far);

static int genl_gtp5g_far_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_FAR_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_FAR_ID:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FAR_APPLY_ACTION:
            if (mnl_attr_validate(attr, MNL_TYPE_U8) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FAR_FORWARDING_PARAMETER:
            if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("FAR mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_forwarding_parameter_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_FORWARDING_PARAMETER_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_FORWARDING_PARAMETER_OUTER_HEADER_CREATION:
             if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("Forwarding Parameter mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_outer_header_creation_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_OUTER_HEADER_CREATION_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_OUTER_HEADER_CREATION_DESCRIPTION:
            if (mnl_attr_validate(attr, MNL_TYPE_U16) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_OUTER_HEADER_CREATION_O_TEID:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_OUTER_HEADER_CREATION_PEER_ADDR_IPV4:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_OUTER_HEADER_CREATION_PORT:
            if (mnl_attr_validate(attr, MNL_TYPE_U16) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("Outer Header Creation mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_attr_list_cb(const struct nlmsghdr *nlh, void *data)
{
    struct nlattr *far_tb[GTP5G_FAR_ATTR_MAX + 1] = {};
    struct nlattr *fwd_param_tb[GTP5G_FORWARDING_PARAMETER_ATTR_MAX + 1] = {};
    struct nlattr *hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_ATTR_MAX + 1] = {};
    char buf[INET_ADDRSTRLEN];
    struct genlmsghdr *genl;

    const char *indent_str = "  ";
    struct in_addr *ipv4 = (void *)buf;

    mnl_attr_parse(nlh, sizeof(*genl), genl_gtp5g_far_validate_cb, far_tb);
    if (far_tb[GTP5G_FAR_ID])
        printf("[FAR No.%u Info]\n", mnl_attr_get_u32(far_tb[GTP5G_FAR_ID]));
    if (far_tb[GTP5G_FAR_APPLY_ACTION])
        printf("%s- Apply Action: %u\n", indent_str, mnl_attr_get_u8(far_tb[GTP5G_FAR_APPLY_ACTION]));

    if (far_tb[GTP5G_FAR_FORWARDING_PARAMETER]) {
        mnl_attr_parse_nested(far_tb[GTP5G_FAR_FORWARDING_PARAMETER], genl_gtp5g_forwarding_parameter_validate_cb, fwd_param_tb);

        printf("%s[Forwarding Parameter Info]\n", indent_str);
        if (fwd_param_tb[GTP5G_FORWARDING_PARAMETER_OUTER_HEADER_CREATION]) {
            mnl_attr_parse_nested(fwd_param_tb[GTP5G_FORWARDING_PARAMETER_OUTER_HEADER_CREATION],
                                  genl_gtp5g_outer_header_creation_validate_cb, hdr_creation_tb);

            printf("%s%s[Outer Header Creation Info]\n", indent_str, indent_str);
            if (hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_DESCRIPTION])
                printf("%s%s%s- Description: %u\n", indent_str, indent_str, indent_str,
                       mnl_attr_get_u16(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_DESCRIPTION]));

            if (hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_O_TEID])
                printf("%s%s%s- Out Teid: %u\n", indent_str, indent_str, indent_str,
                       ntohl(mnl_attr_get_u32(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_O_TEID])));

            if (hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_PEER_ADDR_IPV4]) {
                ipv4->s_addr = mnl_attr_get_u32(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_PEER_ADDR_IPV4]);
                inet_ntop(AF_INET, ipv4, buf, sizeof(buf));
                printf("%s%s%s- RAN IPv4: %s\n", indent_str, indent_str, indent_str, buf);
            }

            if (hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_PORT])
                printf("%s%s%s- Port: %u\n", indent_str, indent_str, indent_str,
                       ntohs(mnl_attr_get_u16(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_PORT])));
        }
    }

    if (far_tb[GTP5G_FAR_RELATED_TO_PDR]) {
        printf("%s- Related PDR ID: ", indent_str);
        u16_id_list_from_kernel_space_print(mnl_attr_get_payload(far_tb[GTP5G_FAR_RELATED_TO_PDR]),
            mnl_attr_get_payload_len(far_tb[GTP5G_FAR_RELATED_TO_PDR]) / (sizeof(uint16_t) / sizeof(char)));
        printf(" (Not a real IE)\n");
    }

    return MNL_CB_OK;
}

int gtp5g_list_far(int genl_id, struct mnl_socket *nl)
{
    char buf[MNL_SOCKET_BUFFER_SIZE];
    struct nlmsghdr *nlh;
    uint32_t seq = time(NULL);

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_DUMP, 0,
                               GTP5G_CMD_GET_FAR);

    if (genl_socket_talk(nl, nlh, seq, genl_gtp5g_attr_list_cb, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_list_far);

void gtp5g_print_far(struct gtp5g_far *far)
{
    struct gtp5g_forwarding_parameter *fwd_param;
    struct gtp5g_outer_header_creation *hdr_creation;

    const char *indent_str = "  ";
    char buf[INET_ADDRSTRLEN];

    if (!far) {
        perror("FAR is NULL");
        return;
    }

    printf("[FAR No.%u Info]\n", far->id);
    printf("%s- Apply Action: %u\n", indent_str, far->apply_action);
    if (far->fwd_param) {
        fwd_param = far->fwd_param;
        printf("%s[Forwarding Parameter Info]\n", indent_str);

        if (fwd_param->hdr_creation) {
            hdr_creation = fwd_param->hdr_creation;
            printf("%s%s[Outer Header Creation Info]\n", indent_str, indent_str);

            printf("%s%s%s- Description: %u\n", indent_str, indent_str, indent_str,
                hdr_creation->desp);
            
            printf("%s%s%s- Out Teid: %u\n", indent_str, indent_str, indent_str,
                ntohl(hdr_creation->teid));

            inet_ntop(AF_INET, &hdr_creation->peer_addr_ipv4, buf, sizeof(buf));
            printf("%s%s%s- RAN IPv4: %s\n", indent_str, indent_str, indent_str, buf);

            printf("%s%s%s- Port: %u\n", indent_str, indent_str, indent_str,
                ntohs(hdr_creation->port));
        }
    }

    if (far->related_pdr_num && far->related_pdr_list) {
        printf("%s- Related PDR ID: ", indent_str);
        u16_id_list_from_kernel_space_print(far->related_pdr_list, far->related_pdr_num);
        printf(" (Not a real IE)\n");
    }

}
EXPORT_SYMBOL(gtp5g_print_far);

static int genl_gtp5g_attr_cb(const struct nlmsghdr *nlh, void *data)
{
    struct nlattr *far_tb[GTP5G_FAR_ATTR_MAX + 1] = {};
    struct nlattr *fwd_param_tb[GTP5G_FORWARDING_PARAMETER_ATTR_MAX + 1] = {};
    struct nlattr *hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_ATTR_MAX + 1] = {};

    struct genlmsghdr *genl;
    struct gtp5g_far *far;
    struct in_addr ipv4;

    mnl_attr_parse(nlh, sizeof(*genl), genl_gtp5g_far_validate_cb, far_tb);
    far = *(struct gtp5g_far **) data = gtp5g_far_alloc();

    if (far_tb[GTP5G_FAR_ID])
        gtp5g_far_set_id(far, mnl_attr_get_u32(far_tb[GTP5G_FAR_ID]));
    if (far_tb[GTP5G_FAR_APPLY_ACTION])
        gtp5g_far_set_apply_action(far, mnl_attr_get_u8(far_tb[GTP5G_FAR_APPLY_ACTION]));

    if (far_tb[GTP5G_FAR_FORWARDING_PARAMETER]) {
        mnl_attr_parse_nested(far_tb[GTP5G_FAR_FORWARDING_PARAMETER], genl_gtp5g_forwarding_parameter_validate_cb, fwd_param_tb);

        if (fwd_param_tb[GTP5G_FORWARDING_PARAMETER_OUTER_HEADER_CREATION]) {
            mnl_attr_parse_nested(fwd_param_tb[GTP5G_FORWARDING_PARAMETER_OUTER_HEADER_CREATION],
                                  genl_gtp5g_outer_header_creation_validate_cb, hdr_creation_tb);

            ipv4.s_addr = mnl_attr_get_u32(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_PEER_ADDR_IPV4]);
            gtp5g_far_set_outer_header_creation(far,
                mnl_attr_get_u16(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_DESCRIPTION]),
                mnl_attr_get_u32(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_O_TEID]),
                &ipv4, mnl_attr_get_u16(hdr_creation_tb[GTP5G_OUTER_HEADER_CREATION_PORT]));
        }
    }

    if (far_tb[GTP5G_FAR_RELATED_TO_PDR]) {
        far->related_pdr_num = mnl_attr_get_payload_len(far_tb[GTP5G_FAR_RELATED_TO_PDR]) / (sizeof(uint16_t) / sizeof(char));
        far->related_pdr_list = calloc(1, mnl_attr_get_payload_len(far_tb[GTP5G_FAR_RELATED_TO_PDR]));
        memcpy(far->related_pdr_list, mnl_attr_get_payload(far_tb[GTP5G_FAR_RELATED_TO_PDR]), mnl_attr_get_payload_len(far_tb[GTP5G_FAR_RELATED_TO_PDR]));
    }

    return MNL_CB_OK;
}

struct gtp5g_far *gtp5g_far_find_by_id(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far)
{
    struct nlmsghdr *nlh;
    char buf[MNL_SOCKET_BUFFER_SIZE];
    struct gtp5g_far *rt_far;
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return NULL;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_EXCL | NLM_F_ACK, ++seq,
                               GTP5G_CMD_GET_FAR);
    gtp5g_build_far_payload(nlh, dev, far);

    if (genl_socket_talk(nl, nlh, seq, genl_gtp5g_attr_cb, &rt_far) < 0) {
        perror("genl_socket_talk");
        return NULL;
    }

    return rt_far;
}
EXPORT_SYMBOL(gtp5g_far_find_by_id);
