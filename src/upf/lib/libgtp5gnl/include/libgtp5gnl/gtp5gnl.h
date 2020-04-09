#ifndef _LIBGTP5GNL_H_
#define _LIBGTP5GNL_H_

#include <stdint.h>

struct mnl_socket;
struct nlmsghdr;

struct mnl_socket *genl_socket_open(void);
void genl_socket_close(struct mnl_socket *nl);
struct nlmsghdr *genl_nlmsg_build_hdr(char *buf, uint16_t type, uint16_t flags,
				      uint32_t seq, uint8_t cmd);
int genl_socket_talk(struct mnl_socket *nl, struct nlmsghdr *nlh, uint32_t seq,
		     int (*cb)(const struct nlmsghdr *nlh, void *data),
		     void *data);
int genl_lookup_family(struct mnl_socket *nl, const char *family);

struct in_addr;

int gtp_dev_create(int dest_ns, const char *gtp_ifname, int fd);
int gtp_dev_create_ran(int dest_ns, const char *gtp_ifname, int fd);
int gtp_dev_config(const char *iface, struct in_addr *net, uint32_t prefix);
int gtp_dev_destroy(const char *gtp_ifname);

struct gtp5g_dev;
struct gtp5g_pdr;
struct gtp5g_far;

int gtp5g_add_pdr(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr);
int gtp5g_add_far(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far);

int gtp5g_mod_pdr(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr);
int gtp5g_mod_far(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far);

int gtp5g_del_pdr(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr);
int gtp5g_del_far(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far);

int gtp5g_list_pdr(int genl_id, struct mnl_socket *nl);
int gtp5g_list_far(int genl_id, struct mnl_socket *nl);

void gtp5g_print_pdr(struct gtp5g_pdr *pdr);
void gtp5g_print_far(struct gtp5g_far *far);

struct gtp5g_pdr *gtp5g_pdr_find_by_id(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_pdr *pdr);
struct gtp5g_far *gtp5g_far_find_by_id(int genl_id, struct mnl_socket *nl, struct gtp5g_dev *dev, struct gtp5g_far *far);

#endif
