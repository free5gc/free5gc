/* Command line utility to create GTP tunnels (PDP contexts) */

/* (C) 2014 by sysmocom - s.f.m.c. GmbH
 * (C) 2016 by Pablo Neira Ayuso <pablo@netfilter.org>
 * (C) 2016 by Pablo Neira Ayuso <pablo@gnumonks.org>
 *
 * Author: Yao-Wen Chang <yaowenowo@gmail.com>
 *
 * All Rights Reserved
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <time.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <net/if.h>
#include <inttypes.h>

#include <libmnl/libmnl.h>
#include <linux/genetlink.h>

#include <linux/gtp5g.h>
#include <linux/if_link.h>
#include <libgtp5gnl/gtp5g.h>
#include <libgtp5gnl/gtp5gnl.h>

static void add_usage(const char *name)
{
    printf("%s <add|mod> <pdr|far> <gtp device> <id> [<options,...>]\n", name);
    printf("%s <del|get> <pdr|far> <gtp device> <id>\n", name);
    printf("%s list <pdr|far>\n\n", name);

    printf("PDR OPTIONS\n");
    printf("\t--pcd <precedence>\n");
    printf("\t--hdr-rm <outer-header-removal>\n");
    printf("\t--far-id <esixted-far-id>\n");
    printf("\t--ue-ipv4 <pdi-ue-ipv4>\n");
    printf("\t--f-teid <i-teid> <local-gtpu-ipv4>\n");
    printf("\t--sdf-desp <description-string>\n");
    printf("\t\tex: --sdf-desp 'permit out ip from 192.168.0.1 22,53,73 to 127.0.0.1/24'\n");
    printf("\t--sdf-tos-traff-cls <tos-traffic-class>\n");
    printf("\t--sdf-scy-param-idx <security-param-idx>\n");
    printf("\t--sdf-flow-label <flow-label>\n");
    printf("\t--sdf-id <id>\n");
    printf("\n");

    printf("OTHER OPTION BUT NOT IEs");
    printf("\t--gtpu-src-ip <gtpu-src-ip>\n");
    printf("\t\tUsed for set the source IP in GTP forwarding packet");
    printf("\t--buffer-usock-path <AF_UNIX-sock-path>\n");
    printf("\t\tUsed for sending packet which should be buffered to user space");
    printf("\n");

    printf("FAR OPTIONS\n");
    printf("\t--action <apply-action>\n");
    printf("\t--hdr-creation <description> <o-teid> <peer-ipv4> <peer-port>\n");
}

static inline int check_opts_invalid(int argc, int now)
{
    return argc <= now;
}

static struct option long_pdr_options[] =
{
    {"pcd", required_argument, NULL, 'p'},
    {"hdr-rm", required_argument, NULL, 'm'},
    {"far-id", required_argument, NULL, 'i'},
    {"ue-ipv4", required_argument, NULL, 'u'},
    {"f-teid", required_argument, NULL, 'f'},
    {"sdf-desp", required_argument, NULL, 'd'},
    {"sdf-tos-traff-cls", required_argument, NULL, 't'},
    {"sdf-scy-param-idx", required_argument, NULL, 's'},
    {"sdf-flow-label", required_argument, NULL, 'l'},
    {"sdf-id", required_argument, NULL, 'j'},

    /* Not in 3GPP spec, just used for routing */
    {"gtpu-src-ip", required_argument, NULL, 'g'},
    /* Not in 3GPP spec, just used for buffering */
    {"buffer-usock-path", required_argument, NULL, 'b'},
};

static struct gtp5g_pdr *prepare_pdr(int argc, char *argv[])
{
    int opt;
    int opt_index = 0;
    struct gtp5g_pdr *pdr;

    int ret;
    struct sockaddr_in sa;

    if (argc < 5)
        goto err;

    pdr = gtp5g_pdr_alloc();
    gtp5g_pdr_set_id(pdr, atoi(argv[4]));

    while((opt = getopt_long_only(argc, argv, "", long_pdr_options, &opt_index)) != EOF) {
        switch(opt) {
            case 'p': // --pcd
                gtp5g_pdr_set_precedence(pdr, atoi(optarg));
                break;
            case 'm': // --hdr-rm
                gtp5g_pdr_set_outer_header_removal(pdr, atoi(optarg));
                break;
            case 'i': //--far-id
                gtp5g_pdr_set_far_id(pdr, atoi(optarg));
                break;
            case 'u': // --ue-ipv4
                if ((ret = inet_pton(AF_INET, optarg, &(sa.sin_addr))) != 1)
                    goto err;
                gtp5g_pdr_set_ue_addr_ipv4(pdr, &(sa.sin_addr));
                break;
            case 'f': // --f-teid {teid} {local-gtpu-ipv4}
                if (check_opts_invalid(argc, optind))
                    goto err;
                if ((ret = inet_pton(AF_INET, argv[optind], &(sa.sin_addr))) != 1)
                    goto err;
                gtp5g_pdr_set_local_f_teid(pdr, htonl(atoi(optarg)), &(sa.sin_addr));
                break;
            case 'd': // --sdf-desp {description string}
                gtp5g_pdr_set_sdf_filter_description(pdr, optarg);
                break;
            case 't': // --sdf-tos-traff-cls
                gtp5g_pdr_set_tos_traffic_class(pdr, atoi(optarg));
                break;
            case 's': // --sdf-scy-param-idx
                gtp5g_pdr_set_security_param_idx(pdr, atoi(optarg));
                break;
            case 'l': // --sdf-flow-label
                gtp5g_pdr_set_flow_label(pdr, atoi(optarg));
                break;
            case 'j': // --sdf-id
                gtp5g_pdr_set_sdf_filter_id(pdr, atoi(optarg));
                break;
            case 'g': // --gtpu-src-ip
                if ((ret = inet_pton(AF_INET, optarg, &(sa.sin_addr))) != 1)
                    goto err;
                gtp5g_pdr_set_role_addr_ipv4(pdr, &(sa.sin_addr));
                break;
            case 'b': // --buffer-usock-path
                gtp5g_pdr_set_unix_sock_path(pdr, optarg);
                break;
        }
    }
    return pdr;

err:
    printf("Invalid opts !!\n");
    return NULL;
}

static int add_pdr(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_pdr *pdr;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    pdr = prepare_pdr(argc, argv);

    gtp5g_add_pdr(genl_id, nl, dev, pdr);

    gtp5g_pdr_free(pdr);
    gtp5g_dev_free(dev);

    return 0;
}

static int mod_pdr(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_pdr *pdr;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 6) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    pdr = prepare_pdr(argc, argv);

    gtp5g_mod_pdr(genl_id, nl, dev, pdr);

    gtp5g_pdr_free(pdr);
    gtp5g_dev_free(dev);

    return 0;
}


static int del_pdr(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_pdr *pdr;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    pdr = gtp5g_pdr_alloc();
    gtp5g_pdr_set_id(pdr, atoi(argv[++optidx]));

    gtp5g_del_pdr(genl_id, nl, dev, pdr);

    gtp5g_pdr_free(pdr);
    gtp5g_dev_free(dev);

    return 0;
}

static int list_pdr(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    return gtp5g_list_pdr(genl_id, nl);
}

static int get_pdr(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_pdr *pdr, *rt_pdr;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    pdr = gtp5g_pdr_alloc();
    gtp5g_pdr_set_id(pdr, atoi(argv[++optidx]));

    rt_pdr = gtp5g_pdr_find_by_id(genl_id, nl, dev, pdr);
    if (!rt_pdr)
        goto FREE;

    gtp5g_print_pdr(rt_pdr);
    gtp5g_pdr_free(rt_pdr);

FREE:
    gtp5g_pdr_free(pdr);
    gtp5g_dev_free(dev);

    return 0;
}

static struct option long_far_options[] =
{
    {"action", required_argument, NULL, 'a'},
    {"hdr-creation", required_argument, NULL, 'h'},
};

static struct gtp5g_far *prepare_far(int argc, char *argv[])
{
    int opt;
    int opt_index = 0;
    struct gtp5g_far *far;

    int ret;
    struct sockaddr_in sa;

    if (argc < 5)
        goto err;

    far = gtp5g_far_alloc();
    gtp5g_far_set_id(far, atoi(argv[4]));

    while((opt = getopt_long_only(argc, argv, "", long_far_options, &opt_index)) != EOF) {
        switch(opt) {
            case 'a': // --action
                gtp5g_far_set_apply_action(far, atoi(optarg));
                break;
            case 'h': // --hdr-creation {desp} {o-teid} {peer-ipv4} {peer-port}
                if (check_opts_invalid(argc, optind + 2))
                    goto err;
                if ((ret = inet_pton(AF_INET, argv[optind + 1], &(sa.sin_addr))) != 1)
                    goto err;
                gtp5g_far_set_outer_header_creation(far, atoi(optarg), htonl(atoi(argv[optind])),
                                                    &(sa.sin_addr), htons(atoi(argv[optind + 2])));
        }
    }
    return far;

err:
    printf("Invalid opts !!\n");
    return NULL;
}

static int add_far(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_far *far;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    far = prepare_far(argc, argv);

    gtp5g_add_far(genl_id, nl, dev, far);

    gtp5g_far_free(far);
    gtp5g_dev_free(dev);

    return 0;
}

static int mod_far(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_far *far;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    far = prepare_far(argc, argv);

    gtp5g_mod_far(genl_id, nl, dev, far);

    gtp5g_far_free(far);
    gtp5g_dev_free(dev);

    return 0;
}


static int del_far(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_far *far;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    far = gtp5g_far_alloc();
    gtp5g_far_set_id(far, atoi(argv[++optidx]));

    gtp5g_del_far(genl_id, nl, dev, far);

    gtp5g_far_free(far);
    gtp5g_dev_free(dev);

    return 0;
}

static int list_far(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    return gtp5g_list_far(genl_id, nl);
}

static int get_far(int argc, char *argv[], int genl_id, struct mnl_socket *nl)
{
    struct gtp5g_dev *dev;
    uint32_t ifidx;
    struct gtp5g_far *far, *rt_far;
    int optidx;

    // TODO: Need to modify argc in release version
    if (argc < 5) {
        add_usage(argv[0]);
        return EXIT_FAILURE;
    }
    optidx = 3;

    ifidx = if_nametoindex(argv[optidx]);
    if (ifidx == 0) {
        fprintf(stderr, "wrong 5G GTP interface %s\n", argv[optidx]);
        return EXIT_FAILURE;
    }
    dev = gtp5g_dev_alloc();
    gtp5g_dev_set_ifidx(dev, ifidx);

    far = gtp5g_far_alloc();
    gtp5g_far_set_id(far, atoi(argv[++optidx]));

    rt_far = gtp5g_far_find_by_id(genl_id, nl, dev, far);
    if (!rt_far)
        goto FREE;

    gtp5g_print_far(rt_far);
    gtp5g_far_free(rt_far);

FREE:
    gtp5g_far_free(far);
    gtp5g_dev_free(dev);

    return 0;
}

int main(int argc, char *argv[])
{
    struct mnl_socket *nl;
    int32_t genl_id;
    int ret = 0;

    if (argc < 3) {
        printf("%s <add|mod|delete|list|get> <pdr|far> [<options,...>]\n", argv[0]);
        exit(EXIT_FAILURE);
    }

    nl = genl_socket_open();
    if (nl == NULL) {
        perror("mnl_socket_open");
        exit(EXIT_FAILURE);
    }

    genl_id = genl_lookup_family(nl, "gtp5g");
    if (genl_id < 0) {
        printf("not found gtp genl family\n");
        exit(EXIT_FAILURE);
    }

    if (strncmp(argv[1], "add", strlen(argv[1])) == 0) {
        if (strncmp(argv[2], "pdr", strlen(argv[2])) == 0)
            ret = add_pdr(argc, argv, genl_id, nl);
        if (strncmp(argv[2], "far", strlen(argv[2])) == 0)
            ret = add_far(argc, argv, genl_id, nl);
    }
    else if (strncmp(argv[1], "mod", strlen(argv[1])) == 0) {
        if (strncmp(argv[2], "pdr", strlen(argv[2])) == 0)
            ret = mod_pdr(argc, argv, genl_id, nl);
        if (strncmp(argv[2], "far", strlen(argv[2])) == 0)
            ret = mod_far(argc, argv, genl_id, nl);
    }
    else if (strncmp(argv[1], "delete", strlen(argv[1])) == 0) {
        if (strncmp(argv[2], "pdr", strlen(argv[2])) == 0)
            ret = del_pdr(argc, argv, genl_id, nl);
        if (strncmp(argv[2], "far", strlen(argv[2])) == 0)
            ret = del_far(argc, argv, genl_id, nl);
    }
    else if (strncmp(argv[1], "list", strlen(argv[1])) == 0) {
        if (strncmp(argv[2], "pdr", strlen(argv[2])) == 0)
            ret = list_pdr(argc, argv, genl_id, nl);
        if (strncmp(argv[2], "far", strlen(argv[2])) == 0)
            ret = list_far(argc, argv, genl_id, nl);
    }
    else if (strncmp(argv[1], "get", strlen(argv[1])) == 0) {
        if (strncmp(argv[2], "pdr", strlen(argv[2])) == 0)
            ret = get_pdr(argc, argv, genl_id, nl);
        if (strncmp(argv[2], "far", strlen(argv[2])) == 0)
            ret = get_far(argc, argv, genl_id, nl);
    }
    else {
        printf("Unknown command `%s'\n", argv[1]);
        exit(EXIT_FAILURE);
    }

    mnl_socket_close(nl);

    return ret;
}
