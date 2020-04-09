#include <netinet/in.h>

static inline int netmask_to_decimal(uint32_t ip_n) {
    uint32_t  ip_h = ntohl(ip_n), ret = 32;
    while ((ip_h & 1) == 0 && ip_h) {
        ret--;
        ip_h >>= 1;
    }
    return ret;
}

static inline uint32_t decimal_to_netmask(uint32_t mask) {
    uint32_t ret = 0, one_block;
    for (int i = 0; i < 4 && mask > 0; i++, mask -= 8) {
        one_block = (mask >= 8 ? 8 : mask);
        ret += (((1 << one_block) - 1) << (8 - one_block)) << (i * 8);
    }
    return ret;
}

static inline uint32_t *port_list_create(char *port_list) {
    uint32_t *ret = calloc(0xff, sizeof(uint32_t));
    uint32_t port1, port2, cnt = 0;

    char *tok_ptr = strtok(port_list, ","), *chr_ptr;
    while (tok_ptr != NULL)  {
        chr_ptr = strchr(tok_ptr, '-');
        if (chr_ptr) {
            *chr_ptr = '\0'; port1 = atoi(tok_ptr); port2 = atoi(chr_ptr + 1);
            if (port1 <= port2)
                ret[++cnt] = port1 + (port2 << 16);
            else
                ret[++cnt] = port2 + (port1 << 16);
        }
        else {
            port1 = atoi(tok_ptr);
            ret[++cnt] = port1 + (port1 << 16);
        }
        tok_ptr = strtok(NULL, ",");
    }
    ret[0] = cnt;

    return ret;
}

static inline void port_list_from_kernel_print(uint32_t *port_list, int len) {
    if (!port_list || len <= 0)
        return;

    int port1, port2;
    for (int i = 0; i < len; i++) {
        port1 = (port_list[i] & 0xFFFF);
        port2 = (port_list[i] >> 16);
        
        if (i)
            printf(",");
        
        if (port1 == port2)
            printf("%d", port1);
        else
            printf("%d-%d", port1, port2);
    }
}

static inline void u16_id_list_from_kernel_space_print(uint16_t *u16_id_list, int len) {
    if (!u16_id_list || len <= 0)
        return;

    printf("%d", u16_id_list[0]);
    for (int i = 1; i < len; i++)
        printf(", %d", u16_id_list[i]);
}