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

static void gtp5g_build_pdr_payload(struct nlmsghdr *nlh, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr)
{
    // Let kernel get dev easily
    if (dev->ifns >= 0)
        mnl_attr_put_u32(nlh, GTP5G_NET_NS_FD, dev->ifns);
    mnl_attr_put_u32(nlh, GTP5G_LINK, dev->ifidx);

    // Level 1 PDR
    mnl_attr_put_u16(nlh, GTP5G_PDR_ID, pdr->id);
    if (pdr->precedence)
        mnl_attr_put_u32(nlh, GTP5G_PDR_PRECEDENCE, *pdr->precedence);
    if (pdr->outer_hdr_removal)
        mnl_attr_put_u8(nlh, GTP5G_OUTER_HEADER_REMOVAL, *pdr->outer_hdr_removal);
    if (pdr->far_id)
        mnl_attr_put_u32(nlh, GTP5G_PDR_FAR_ID, *pdr->far_id);
    
    /* Not in 3GPP spec, just used for routing */
    if (pdr->role_addr_ipv4)
        mnl_attr_put_u32(nlh, GTP5G_PDR_ROLE_ADDR_IPV4, pdr->role_addr_ipv4->s_addr);

    /* Not in 3GPP spec, just used for buffering */
    if (pdr->unix_sock_path)
        mnl_attr_put_str(nlh, GTP5G_PDR_UNIX_SOCKET_PATH, pdr->unix_sock_path);

    // Level 2 PDR : PDI
    struct gtp5g_pdi *pdi = pdr->pdi;
    struct nlattr *pdi_nest, *f_teid_nest, *sdf_filter_nest, *sdf_desp_nest;
    if (pdi) {
        pdi_nest = mnl_attr_nest_start(nlh, GTP5G_PDR_PDI);
        if (pdi->ue_addr_ipv4)
            mnl_attr_put_u32(nlh, GTP5G_PDI_UE_ADDR_IPV4, pdi->ue_addr_ipv4->s_addr);

        // Level 3 : local f-teid
        struct local_f_teid *f_teid = pdi->f_teid;
        if (f_teid) {
            f_teid_nest = mnl_attr_nest_start(nlh, GTP5G_PDI_F_TEID);
            mnl_attr_put_u32(nlh, GTP5G_F_TEID_I_TEID, f_teid->teid);
            mnl_attr_put_u32(nlh, GTP5G_F_TEID_GTPU_ADDR_IPV4, f_teid->gtpu_addr_ipv4.s_addr);
            mnl_attr_nest_end(nlh, f_teid_nest);
        }

        // Level 3 : SDF Filter
        struct sdf_filter *sdf = pdi->sdf;
        if (sdf) {
            sdf_filter_nest = mnl_attr_nest_start(nlh, GTP5G_PDI_SDF_FILTER);

            // Level 4 : SDF Filter description
            struct ip_filter_rule *rule = sdf->rule;
            if (rule) {
                sdf_desp_nest = mnl_attr_nest_start(nlh, GTP5G_SDF_FILTER_FLOW_DESCRIPTION);
                mnl_attr_put_u8(nlh, GTP5G_FLOW_DESCRIPTION_ACTION, rule->action);
                mnl_attr_put_u8(nlh, GTP5G_FLOW_DESCRIPTION_DIRECTION, rule->direction);
                mnl_attr_put_u8(nlh, GTP5G_FLOW_DESCRIPTION_PROTOCOL, rule->proto);
                mnl_attr_put_u32(nlh, GTP5G_FLOW_DESCRIPTION_SRC_IPV4, rule->src.s_addr);
                mnl_attr_put_u32(nlh, GTP5G_FLOW_DESCRIPTION_SRC_MASK, rule->smask.s_addr);
                if (rule->sport_list)
                    mnl_attr_put(nlh, GTP5G_FLOW_DESCRIPTION_SRC_PORT,
                                 rule->sport_list[0] * sizeof(uint32_t) / sizeof(char),
                                 (void *) &rule->sport_list[1]);
                mnl_attr_put_u32(nlh, GTP5G_FLOW_DESCRIPTION_DEST_IPV4, rule->dest.s_addr);
                mnl_attr_put_u32(nlh, GTP5G_FLOW_DESCRIPTION_DEST_MASK, rule->dmask.s_addr);
                if (rule->dport_list)
                    mnl_attr_put(nlh, GTP5G_FLOW_DESCRIPTION_DEST_PORT,
                                 rule->dport_list[0] * sizeof(uint32_t) / sizeof(char),
                                 (void *) &rule->dport_list[1]);
                mnl_attr_nest_end(nlh, sdf_desp_nest);
            }

            if (sdf->tos_traffic_class)
                mnl_attr_put_u16(nlh, GTP5G_SDF_FILTER_TOS_TRAFFIC_CLASS, *sdf->tos_traffic_class);

            if (sdf->security_param_idx)
                mnl_attr_put_u32(nlh, GTP5G_SDF_FILTER_SECURITY_PARAMETER_INDEX, *sdf->security_param_idx);

            if (sdf->flow_label)
                mnl_attr_put_u32(nlh, GTP5G_SDF_FILTER_FLOW_LABEL, *sdf->flow_label);

            if (sdf->bi_id)
                mnl_attr_put_u32(nlh, GTP5G_SDF_FILTER_SDF_FILTER_ID, *sdf->bi_id);

            mnl_attr_nest_end(nlh, sdf_filter_nest);
        }
        mnl_attr_nest_end(nlh, pdi_nest);
    }
}

int gtp5g_add_pdr(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr)
{
    struct nlmsghdr *nlh;
    char buf[MNL_SOCKET_BUFFER_SIZE];
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return -1;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_EXCL | NLM_F_ACK, ++seq,
                               GTP5G_CMD_ADD_PDR);

    // Add mandatory IEs here
    if (!pdr->precedence) {
        perror("Add PDR must have precedence");
        return -1;
    }

    gtp5g_build_pdr_payload(nlh, dev, pdr);

    if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_add_pdr);

int gtp5g_mod_pdr(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr)
{
    struct nlmsghdr *nlh;
    char buf[MNL_SOCKET_BUFFER_SIZE];
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return -1;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_REPLACE | NLM_F_ACK, ++seq,
                               GTP5G_CMD_ADD_PDR);
    gtp5g_build_pdr_payload(nlh, dev, pdr);

    if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_mod_pdr);

int gtp5g_del_pdr(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr)
{
    char buf[MNL_SOCKET_BUFFER_SIZE];
    struct nlmsghdr *nlh;
    uint32_t seq = time(NULL);
    
    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return -1;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_ACK, ++seq,
                               GTP5G_CMD_DEL_PDR);
    gtp5g_build_pdr_payload(nlh, dev, pdr);

    if (genl_socket_talk(nl, nlh, seq, NULL, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_del_pdr);

static int genl_gtp5g_pdr_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_PDR_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_PDR_ID:
            if (mnl_attr_validate(attr, MNL_TYPE_U16) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_PDR_PRECEDENCE:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_PDR_PDI:
            if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_OUTER_HEADER_REMOVAL:
            if (mnl_attr_validate(attr, MNL_TYPE_U8) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_PDR_FAR_ID:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;

        /* Not in 3GPP spec, just used for routing */
        case GTP5G_PDR_ROLE_ADDR_IPV4:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("PDR mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_pdi_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_PDI_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_PDI_UE_ADDR_IPV4:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_PDI_F_TEID:
             if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_PDI_SDF_FILTER:
             if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("PDI mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_f_teid_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_F_TEID_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_F_TEID_I_TEID:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_F_TEID_GTPU_ADDR_IPV4:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("F-TEID mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_sdf_filter_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_SDF_FILTER_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_SDF_FILTER_FLOW_DESCRIPTION:
             if (mnl_attr_validate(attr, MNL_TYPE_NESTED) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_SDF_FILTER_TOS_TRAFFIC_CLASS:
            if (mnl_attr_validate(attr, MNL_TYPE_U16) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_SDF_FILTER_SECURITY_PARAMETER_INDEX:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_SDF_FILTER_FLOW_LABEL:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_SDF_FILTER_SDF_FILTER_ID:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("SDF Filter mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_flow_description_validate_cb(const struct nlattr *attr, void *data)
{
    const struct nlattr **tb = data;
    int type = mnl_attr_get_type(attr);

    if (mnl_attr_type_valid(attr, GTP5G_FLOW_DESCRIPTION_ATTR_MAX) < 0)
        return MNL_CB_OK;

    switch(type) {
        case GTP5G_FLOW_DESCRIPTION_ACTION:
            if (mnl_attr_validate(attr, MNL_TYPE_U8) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FLOW_DESCRIPTION_DIRECTION:
            if (mnl_attr_validate(attr, MNL_TYPE_U8) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FLOW_DESCRIPTION_PROTOCOL:
            if (mnl_attr_validate(attr, MNL_TYPE_U8) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FLOW_DESCRIPTION_SRC_IPV4:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FLOW_DESCRIPTION_SRC_MASK:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FLOW_DESCRIPTION_DEST_IPV4:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
        case GTP5G_FLOW_DESCRIPTION_DEST_MASK:
            if (mnl_attr_validate(attr, MNL_TYPE_U32) < 0)
                goto VALIDATE_FAIL;
            break;
    default:
        break;
    }

    tb[type] = attr;
    return MNL_CB_OK;

VALIDATE_FAIL:
    perror("Flow Description mnl_attr_validate");
    return MNL_CB_ERROR;
}

static int genl_gtp5g_attr_list_cb(const struct nlmsghdr *nlh, void *data)
{
    struct nlattr *pdr_tb[GTP5G_PDR_ATTR_MAX + 1] = {};
    struct nlattr *pdi_tb[GTP5G_PDI_ATTR_MAX + 1] = {};
    struct nlattr *f_teid_tb[GTP5G_F_TEID_ATTR_MAX + 1] = {};
    struct nlattr *sdf_tb[GTP5G_SDF_FILTER_ATTR_MAX + 1] = {};
    struct nlattr *rule_tb[GTP5G_FLOW_DESCRIPTION_ATTR_MAX + 1] = {};

    char buf[INET_ADDRSTRLEN];
    struct genlmsghdr *genl;

    const char *indent_str = "  ";
    struct in_addr *ipv4 = (void *)buf;
    int mask;

    mnl_attr_parse(nlh, sizeof(*genl), genl_gtp5g_pdr_validate_cb, pdr_tb);
    if (pdr_tb[GTP5G_PDR_ID])
        printf("[PDR No.%u Info]\n", mnl_attr_get_u16(pdr_tb[GTP5G_PDR_ID]));
    if (pdr_tb[GTP5G_PDR_PRECEDENCE])
        printf("%s- Precedence: %u\n", indent_str, mnl_attr_get_u32(pdr_tb[GTP5G_PDR_PRECEDENCE]));
    if (pdr_tb[GTP5G_OUTER_HEADER_REMOVAL])
        printf("%s- Outer Header Removal: %u\n", indent_str, mnl_attr_get_u8(pdr_tb[GTP5G_OUTER_HEADER_REMOVAL]));

    if (pdr_tb[GTP5G_PDR_PDI]) {
        mnl_attr_parse_nested(pdr_tb[GTP5G_PDR_PDI], genl_gtp5g_pdi_validate_cb, pdi_tb);

        printf("%s[PDI Info]\n", indent_str);
        if (pdi_tb[GTP5G_PDI_UE_ADDR_IPV4]) {
            ipv4->s_addr = mnl_attr_get_u32(pdi_tb[GTP5G_PDI_UE_ADDR_IPV4]);
            inet_ntop(AF_INET, ipv4, buf, sizeof(buf));
            printf("%s%s- UE IPv4: %s\n", indent_str, indent_str, buf);
        }

        if (pdi_tb[GTP5G_PDI_F_TEID]) {
            mnl_attr_parse_nested(pdi_tb[GTP5G_PDI_F_TEID], genl_gtp5g_f_teid_validate_cb, f_teid_tb);

            printf("%s%s[Local F-Teid Info]\n", indent_str, indent_str);
            if (f_teid_tb[GTP5G_F_TEID_I_TEID])
                printf("%s%s%s- In Teid: %u\n", indent_str, indent_str, indent_str,
                       ntohl(mnl_attr_get_u32(f_teid_tb[GTP5G_F_TEID_I_TEID])));

            if (f_teid_tb[GTP5G_F_TEID_GTPU_ADDR_IPV4]) {
                ipv4->s_addr = mnl_attr_get_u32(f_teid_tb[GTP5G_F_TEID_GTPU_ADDR_IPV4]);
                inet_ntop(AF_INET, ipv4, buf, sizeof(buf));
                printf("%s%s%s- UPF IPv4: %s\n", indent_str, indent_str, indent_str, buf);
            }
        }

        if (pdi_tb[GTP5G_PDI_SDF_FILTER]) {
            mnl_attr_parse_nested(pdi_tb[GTP5G_PDI_SDF_FILTER], genl_gtp5g_sdf_filter_validate_cb, sdf_tb);

            printf("%s%s[SDF Filter Info]\n", indent_str, indent_str);

            if (sdf_tb[GTP5G_SDF_FILTER_FLOW_DESCRIPTION]) {
                mnl_attr_parse_nested(sdf_tb[GTP5G_SDF_FILTER_FLOW_DESCRIPTION],
                                      genl_gtp5g_flow_description_validate_cb, rule_tb);
                printf("%s%s%s- Flow Description:", indent_str, indent_str, indent_str);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_ACTION]) {
                    switch (mnl_attr_get_u8(rule_tb[GTP5G_FLOW_DESCRIPTION_ACTION])) {
                        case GTP5G_SDF_FILTER_PERMIT:
                            printf(" permit");
                            break;
                        default:
                            printf(" unknown_action");
                    }
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DIRECTION]) {
                    switch (mnl_attr_get_u8(rule_tb[GTP5G_FLOW_DESCRIPTION_DIRECTION])) {
                        case GTP5G_SDF_FILTER_IN:
                            printf(" in");
                            break;
                        case GTP5G_SDF_FILTER_OUT:
                            printf(" out");
                            break;
                        default:
                            printf(" unknown_direction");
                    }
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_PROTOCOL]) {
                    int proto = mnl_attr_get_u8(rule_tb[GTP5G_FLOW_DESCRIPTION_PROTOCOL]);
                    if (proto == 0xff)
                        printf(" ip");
                    else
                        printf(" %d", proto);
                }
                
                printf(" from ");

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_IPV4]) {
                    ipv4->s_addr = mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_IPV4]);
                    if (ipv4->s_addr == 0)
                        printf("any");
                    else {
                        inet_ntop(AF_INET, ipv4, buf, sizeof(buf));
                        printf("%s", buf);
                    }
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_MASK]) {
                    mask = netmask_to_decimal(mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_MASK]));
                    if (mask < 32)
                        printf("/%d", mask);
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]) {
                    printf(" ");
                    port_list_from_kernel_print(mnl_attr_get_payload(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]),
                        mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]) / (sizeof(uint32_t) / sizeof(char)));
                }

                printf(" to ");

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_IPV4]) {
                    ipv4->s_addr = mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_IPV4]);
                    if (ipv4->s_addr == 0)
                        printf("assigned");
                    else {
                        inet_ntop(AF_INET, ipv4, buf, sizeof(buf));
                        printf("%s", buf);
                    }
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_MASK]) {
                    mask = netmask_to_decimal(mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_MASK]));
                    if (mask < 32)
                        printf("/%d", mask);
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]) {
                    printf(" ");
                    port_list_from_kernel_print(mnl_attr_get_payload(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]),
                        mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]) / (sizeof(uint32_t) / sizeof(char)));
                }

                printf("\n");
            }

            if (sdf_tb[GTP5G_SDF_FILTER_TOS_TRAFFIC_CLASS])
                printf("%s%s%s- ToS Traffic Class: %u\n", indent_str, indent_str, indent_str,
                       mnl_attr_get_u16(sdf_tb[GTP5G_SDF_FILTER_TOS_TRAFFIC_CLASS]));

            if (sdf_tb[GTP5G_SDF_FILTER_SECURITY_PARAMETER_INDEX])
                printf("%s%s%s- Security Parameter Index: %u\n", indent_str, indent_str, indent_str,
                       mnl_attr_get_u32(sdf_tb[GTP5G_SDF_FILTER_SECURITY_PARAMETER_INDEX]));

            if (sdf_tb[GTP5G_SDF_FILTER_FLOW_LABEL])
                printf("%s%s%s- Flow Label: %u\n", indent_str, indent_str, indent_str,
                       mnl_attr_get_u32(sdf_tb[GTP5G_SDF_FILTER_FLOW_LABEL]));

            if (sdf_tb[GTP5G_SDF_FILTER_SDF_FILTER_ID])
                printf("%s%s%s- SDF Filter ID: %u\n", indent_str, indent_str, indent_str,
                       mnl_attr_get_u32(sdf_tb[GTP5G_SDF_FILTER_SDF_FILTER_ID]));
        }
    }

    if (pdr_tb[GTP5G_PDR_FAR_ID])
        printf("%s- FAR ID: %u\n", indent_str, mnl_attr_get_u32(pdr_tb[GTP5G_PDR_FAR_ID]));

    /* Not in 3GPP spec, just used for routing */
    if (pdr_tb[GTP5G_PDR_ROLE_ADDR_IPV4]) {
        ipv4->s_addr = mnl_attr_get_u32(pdr_tb[GTP5G_PDR_ROLE_ADDR_IPV4]);
        inet_ntop(AF_INET, ipv4, buf, sizeof(buf));
        printf("%s- GTP-U IPv4: %s (For routing)\n", indent_str, buf);
    }

    return MNL_CB_OK;
}

int gtp5g_list_pdr(int genl_id, struct mnl_socket *nl)
{
    char buf[MNL_SOCKET_BUFFER_SIZE];
    struct nlmsghdr *nlh;
    uint32_t seq = time(NULL);

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_DUMP, 0,
                               GTP5G_CMD_GET_PDR);

    if (genl_socket_talk(nl, nlh, seq, genl_gtp5g_attr_list_cb, NULL) < 0) {
        perror("genl_socket_talk");
        return -1;
    }

    return 0;
}
EXPORT_SYMBOL(gtp5g_list_pdr);

void gtp5g_print_pdr(struct gtp5g_pdr *pdr)
{
    struct gtp5g_pdi *pdi;
    struct local_f_teid *f_teid;
    struct sdf_filter *sdf;
    struct ip_filter_rule *rule;

    const char *indent_str = "  ";
    char buf[INET_ADDRSTRLEN];
    int mask;

    if (!pdr) {
        perror("PDR is NULL");
        return;
    }

    printf("[PDR No.%u Info]\n", pdr->id);
    printf("%s- Precedence: %u\n", indent_str, *pdr->precedence);
    if (pdr->outer_hdr_removal)
        printf("%s- Outer Header Removal: %u\n", indent_str, *pdr->outer_hdr_removal);
    
    if (pdr->pdi) {
        pdi = pdr->pdi;
        printf("%s[PDI Info]\n", indent_str);

        if (pdi->ue_addr_ipv4) {
            inet_ntop(AF_INET, pdi->ue_addr_ipv4, buf, sizeof(buf));
            printf("%s%s- UE IPv4: %s\n", indent_str, indent_str, buf);
        }

        if (pdi->f_teid) {
            f_teid = pdi->f_teid;
            printf("%s%s[Local F-Teid Info]\n", indent_str, indent_str);

            printf("%s%s%s- In Teid: %u\n", indent_str, indent_str, indent_str, ntohl(f_teid->teid));

            inet_ntop(AF_INET, &f_teid->gtpu_addr_ipv4, buf, sizeof(buf));
            printf("%s%s%s- UPF IPv4: %s\n", indent_str, indent_str, indent_str, buf);
        }

        if (pdi->sdf) {
            sdf = pdi->sdf;
            printf("%s%s[SDF Filter Info]\n", indent_str, indent_str);

            if (sdf->rule) {
                rule = sdf->rule;
                printf("%s%s%s- Flow Description:", indent_str, indent_str, indent_str);

                switch (rule->action) {
                    case GTP5G_SDF_FILTER_PERMIT:
                        printf(" permit");
                        break;
                    default:
                        printf(" unknown_action");
                }

                switch (rule->direction) {
                    case GTP5G_SDF_FILTER_IN:
                        printf(" in");
                        break;
                    case GTP5G_SDF_FILTER_OUT:
                        printf(" out");
                        break;
                    default:
                        printf(" unknown_direction");
                }

                if (rule->proto == 0xff)
                    printf(" ip");
                else
                    printf(" %d", rule->proto);

                printf(" from ");

                if (!rule->src.s_addr)
                    printf("any");
                else {
                    inet_ntop(AF_INET, &rule->src, buf, sizeof(buf));
                    printf("%s", buf);
                }

                mask = netmask_to_decimal(rule->smask.s_addr);
                if (mask < 32)
                    printf("/%d", mask);

                if (rule->sport_num && rule->sport_list) {
                    printf(" ");
                    port_list_from_kernel_print(rule->sport_list, rule->sport_num);
                }

                printf(" to ");

                if (!rule->dest.s_addr)
                    printf("assigned");
                else {
                    inet_ntop(AF_INET, &rule->dest, buf, sizeof(buf));
                    printf("%s", buf);
                }

                mask = netmask_to_decimal(rule->dmask.s_addr);
                if (mask < 32)
                    printf("/%d", mask);

                if (rule->dport_num && rule->dport_list) {
                    printf(" ");
                    port_list_from_kernel_print(rule->dport_list, rule->dport_num);
                }

                printf("\n");
            }

            if (sdf->tos_traffic_class)
                printf("%s%s%s- ToS Traffic Class: %u\n", indent_str, indent_str, indent_str,
                    *sdf->tos_traffic_class);

            if (sdf->security_param_idx)
                printf("%s%s%s- Security Parameter Index: %u\n", indent_str, indent_str, indent_str,
                    *sdf->security_param_idx);

            if (sdf->flow_label)
                printf("%s%s%s- Flow Label: %u\n", indent_str, indent_str, indent_str,
                    *sdf->flow_label);

            if (sdf->bi_id)
                printf("%s%s%s- SDF Filter ID: %u\n", indent_str, indent_str, indent_str,
                    *sdf->bi_id);
        }
    }

    if (pdr->far_id)
        printf("%s- FAR ID: %u\n", indent_str, *pdr->far_id);

    /* Not in 3GPP spec, just used for routing */
    if (pdr->role_addr_ipv4) {
        inet_ntop(AF_INET, pdr->role_addr_ipv4, buf, sizeof(buf));
        printf("%s- GTP-U IPv4: %s (For routing)\n", indent_str, buf);
    }
}
EXPORT_SYMBOL(gtp5g_print_pdr);

static int genl_gtp5g_attr_cb(const struct nlmsghdr *nlh, void *data)
{
    struct nlattr *pdr_tb[GTP5G_PDR_ATTR_MAX + 1] = {};
    struct nlattr *pdi_tb[GTP5G_PDI_ATTR_MAX + 1] = {};
    struct nlattr *f_teid_tb[GTP5G_F_TEID_ATTR_MAX + 1] = {};
    struct nlattr *sdf_tb[GTP5G_SDF_FILTER_ATTR_MAX + 1] = {};
    struct nlattr *rule_tb[GTP5G_FLOW_DESCRIPTION_ATTR_MAX + 1] = {};

    struct genlmsghdr *genl;
    struct gtp5g_pdr *pdr;
    struct ip_filter_rule *rule;
    struct in_addr ipv4;
    const char *pstr;

    mnl_attr_parse(nlh, sizeof(*genl), genl_gtp5g_pdr_validate_cb, pdr_tb);
    pdr = *(struct gtp5g_pdr **) data = gtp5g_pdr_alloc();

    if (pdr_tb[GTP5G_PDR_ID])
        gtp5g_pdr_set_id(pdr, mnl_attr_get_u16(pdr_tb[GTP5G_PDR_ID]));

    if (pdr_tb[GTP5G_PDR_PRECEDENCE])
        gtp5g_pdr_set_precedence(pdr, mnl_attr_get_u32(pdr_tb[GTP5G_PDR_PRECEDENCE]));

    if (pdr_tb[GTP5G_OUTER_HEADER_REMOVAL])
        gtp5g_pdr_set_outer_header_removal(pdr, mnl_attr_get_u8(pdr_tb[GTP5G_OUTER_HEADER_REMOVAL]));

    if (pdr_tb[GTP5G_PDR_PDI]) {
        mnl_attr_parse_nested(pdr_tb[GTP5G_PDR_PDI], genl_gtp5g_pdi_validate_cb, pdi_tb);

        if (pdi_tb[GTP5G_PDI_UE_ADDR_IPV4]) {
            ipv4.s_addr = mnl_attr_get_u32(pdi_tb[GTP5G_PDI_UE_ADDR_IPV4]);
            gtp5g_pdr_set_ue_addr_ipv4(pdr, &ipv4);
        }

        if (pdi_tb[GTP5G_PDI_F_TEID]) {
            mnl_attr_parse_nested(pdi_tb[GTP5G_PDI_F_TEID], genl_gtp5g_f_teid_validate_cb, f_teid_tb);

            ipv4.s_addr = mnl_attr_get_u32(f_teid_tb[GTP5G_F_TEID_GTPU_ADDR_IPV4]);
            gtp5g_pdr_set_local_f_teid(pdr, mnl_attr_get_u32(f_teid_tb[GTP5G_F_TEID_I_TEID]), &ipv4);
        }

        if (pdi_tb[GTP5G_PDI_SDF_FILTER]) {
            mnl_attr_parse_nested(pdi_tb[GTP5G_PDI_SDF_FILTER], genl_gtp5g_sdf_filter_validate_cb, sdf_tb);

            if (sdf_tb[GTP5G_SDF_FILTER_FLOW_DESCRIPTION]) {
                mnl_attr_parse_nested(sdf_tb[GTP5G_SDF_FILTER_FLOW_DESCRIPTION],
                                      genl_gtp5g_flow_description_validate_cb, rule_tb);

                if (!pdr->pdi->sdf)
                    pdr->pdi->sdf = calloc(1, sizeof(*pdr->pdi->sdf));

                rule = pdr->pdi->sdf->rule;

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_ACTION])
                    rule->action = mnl_attr_get_u8(rule_tb[GTP5G_FLOW_DESCRIPTION_ACTION]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DIRECTION])
                    rule->direction = mnl_attr_get_u8(rule_tb[GTP5G_FLOW_DESCRIPTION_DIRECTION]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_PROTOCOL])
                    rule->proto = mnl_attr_get_u8(rule_tb[GTP5G_FLOW_DESCRIPTION_PROTOCOL]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_IPV4])
                    rule->src.s_addr = mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_IPV4]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_MASK])
                    rule->smask.s_addr = mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_MASK]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]) {
                    rule->sport_num = mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]) / (sizeof(uint32_t) / sizeof(char));
                    rule->sport_list = calloc(1, mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]));
                    memcpy(rule->sport_list, mnl_attr_get_payload(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]), mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_SRC_PORT]));
                }

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_IPV4])
                    rule->dest.s_addr = mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_IPV4]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_MASK])
                    rule->dest.s_addr = mnl_attr_get_u32(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_MASK]);

                if (rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]) {
                    rule->dport_num = mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]) / (sizeof(uint32_t) / sizeof(char));
                    rule->dport_list = calloc(1, mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]));
                    memcpy(rule->dport_list, mnl_attr_get_payload(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]), mnl_attr_get_payload_len(rule_tb[GTP5G_FLOW_DESCRIPTION_DEST_PORT]));
                }
            }

            if (sdf_tb[GTP5G_SDF_FILTER_TOS_TRAFFIC_CLASS])
                gtp5g_pdr_set_tos_traffic_class(pdr, mnl_attr_get_u16(sdf_tb[GTP5G_SDF_FILTER_TOS_TRAFFIC_CLASS]));

            if (sdf_tb[GTP5G_SDF_FILTER_SECURITY_PARAMETER_INDEX])
                gtp5g_pdr_set_security_param_idx(pdr, mnl_attr_get_u32(sdf_tb[GTP5G_SDF_FILTER_SECURITY_PARAMETER_INDEX]));

            if (sdf_tb[GTP5G_SDF_FILTER_FLOW_LABEL])
                gtp5g_pdr_set_flow_label(pdr, mnl_attr_get_u32(sdf_tb[GTP5G_SDF_FILTER_FLOW_LABEL]));

            if (sdf_tb[GTP5G_SDF_FILTER_SDF_FILTER_ID])
                gtp5g_pdr_set_sdf_filter_id(pdr, mnl_attr_get_u32(sdf_tb[GTP5G_SDF_FILTER_SDF_FILTER_ID]));
        }
    }

    if (pdr_tb[GTP5G_PDR_FAR_ID])
        gtp5g_pdr_set_far_id(pdr, mnl_attr_get_u32(pdr_tb[GTP5G_PDR_FAR_ID]));

    /* Not in 3GPP spec, just used for routing */
    if (pdr_tb[GTP5G_PDR_ROLE_ADDR_IPV4]) {
        ipv4.s_addr = mnl_attr_get_u32(pdr_tb[GTP5G_PDR_ROLE_ADDR_IPV4]);
        gtp5g_pdr_set_role_addr_ipv4(pdr, &ipv4);
    }

    /* Not in 3GPP spec, just used for buffering */
    if (pdr_tb[GTP5G_PDR_UNIX_SOCKET_PATH]) {
        pstr = mnl_attr_get_str(pdr_tb[GTP5G_PDR_UNIX_SOCKET_PATH]);
        gtp5g_pdr_set_unix_sock_path(pdr, pstr);
    }

    return MNL_CB_OK;
}

struct gtp5g_pdr *gtp5g_pdr_find_by_id(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr)
{
    char buf[MNL_SOCKET_BUFFER_SIZE];
    struct nlmsghdr *nlh;
    struct gtp5g_pdr *rt_pdr = NULL;
    uint32_t seq = time(NULL);

    if (!dev) {
        fprintf(stderr, "5G GTP device is NULL\n");
        return NULL;
    }

    nlh = genl_nlmsg_build_hdr(buf, genl_id, NLM_F_ACK, ++seq,
                               GTP5G_CMD_GET_PDR);
    gtp5g_build_pdr_payload(nlh, dev, pdr);

    if (genl_socket_talk(nl, nlh, seq, genl_gtp5g_attr_cb, &rt_pdr) < 0) {
        perror("genl_socket_talk");
        return NULL;
    }

    return rt_pdr;
}
EXPORT_SYMBOL(gtp5g_pdr_find_by_id);
