#ifndef __NETWORK_H__
#define __NETWORK_H__

#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <netinet/ip.h>
#include <sys/epoll.h>

#include "utlt_debug.h"
#include "utlt_lib.h"
#include "utlt_list.h"
#include "utlt_pool.h"

#ifdef __cplusplus
extern "C" {
#endif /* __cplusplus */

#define MAX_NUM_OF_EVENT 8192

typedef struct _Sock Sock;
typedef struct _SockAddr SockAddr;
typedef Status (*SockHandler)(Sock *sockPtr, void *data);

struct _SockAddr {
#define _family sa.sa_family
#define _port s4.sin_port
    union {
        struct sockaddr         sa;
        struct sockaddr_in      s4;
        struct sockaddr_in6     s6;
        struct sockaddr_storage ss;
        struct sockaddr_un      su; // Used for unix socket
    };
    // User define fields
    SockAddr *next; // this is for dual-stack / multi IP situation used
};

struct _Sock {
    int fd;
    SockAddr localAddr;
    SockAddr remoteAddr;
    int rflag; // This flag is used to flags in recvfrom
    int wflag; // This flag is used to flags in sendto

    SockHandler handler;
    void *data;

    int epollMode;
};

typedef int (*SockHandler)(Sock *sock, void *data);

// IP type
typedef struct _Ip {
    union {
        struct in_addr      addr4;
        struct in6_addr     addr6;
        struct {
            struct in_addr  addr4;
            struct in6_addr addr6;
        } dualStack;
    };
    uint32_t                len;
    ENDIAN3(uint8_t         ipv4:1;,
            uint8_t         ipv6:1;,
            uint8_t         reserved:6;)
} __attribute__ ((packed)) Ip;

// List of socket node
typedef struct _SockNode {
    ListNode node;
    char ip[INET6_ADDRSTRLEN];
    Sock* sock;
} SockNode;

// Socket
Status SockPoolInit();
Status SockPoolFinal();
Sock *SockCreate(int domain, int type, int protocol);
Status SockFree(Sock *sock);
Status SockListFree(ListNode *list);
SockNode *SockNodeListAdd(ListNode *list, const char *ip);
Status SockNodeListFree(ListNode *list);

Status SockBind(Sock *sock, SockAddr *sockAddr);
Status SockConnect(Sock *sock, SockAddr *sockAddr);
#define SockSetOpt(__sock, __level, __optname, __optval) \
    setsockopt(__sock->fd, __level, __optname, __optval, sizeof(__optval))
Status SockListen(Sock *sock, int queueLen);
Status SockAccept(Sock *lSock, Sock *acSock);
#define SockSetReadFlag(__sock, __flag) (__sock)->rflag = __flag
#define SockSetWriteFlag(__sock, __flag) (__sock)->wflag = __flag
#define SockSetEpollMode(__sock, __flag) ((__sock)->epollMode = __flag)

int SockRead(Sock *sock, void *buffer, int size);
Status SockWrite(Sock *sock, const void *buffer, int size);
int SockRecvFrom(Sock *sock, void *buffer, int size);
Status SockSendTo(Sock *sock, void *buffer, int size);

Status SockRegister(Sock *sock, SockHandler handler, void *data);
Status SockUnregister(Sock *sock);

// Socket Address
Status SockSetAddr(SockAddr *sockAddr, int domain, const char *addr, int port);
Status SockAddrSet(SockAddr *sockAddr, int domain, const char *addr, int port);
Status SockAddrFilter(SockAddr **sa_list, int family);
Status SockAddrCopy(SockAddr **dst, const SockAddr **src);
Status SockAddrSort(SockAddr **sa_list, int family);
Status SockAddrFillScopeIdInLocal(SockAddr *saList);
Status SockAddrFreeAll(SockAddr *saList);
int SockAddrLen(const void *sockAddr);

// Utilities
const char *UTLT_InetNtop(const void *sockAddr);
#define SockAddrIPPrint(__sockAddr) UTLT_InetNtop(__sockAddr)
#define GetIP(__sockAddr) SockAddrIPPrint(__sockAddr)

const int SockAddrPortPrint(const void *sockAddr);
#define GetPort(__sockAddr) SockAddrPortPrint(__sockAddr)

Status GetAddrFromHost(char *addr, const char *host, size_t maxAddrLen); // host: hostname or ip

// UDP
Sock *UdpSockCreate(int domain);
#define UdpFree(__sock) SockFree(__sock)

Status UdpSockSetAddr(SockAddr *sockAddr, int domain, const char *addr, int port);
Sock *UdpServerCreate(int domain, const char *addr, int port);
Sock *UdpClientCreate(int domain, const char *addr, int port);
#define UdpRecvFrom(__sock, __buffer, __size) \
        SockRecvFrom(__sock, __buffer, __size)
#define UdpSendTo(__sock, __buffer, __size) \
        SockSendTo(__sock, __buffer, __size)

// Unix Socket (AF_UNIX)
Sock *UnixSockCreate(int type);
Status UnixFree(Sock *sock);
Status UnixSockSetAddr(SockAddr *sockAddr, const char *path);
Sock *UnixServerCreate(int type, const char *path);
#define UnixRecv(__sock, __buffer, __size) \
        SockRead(__sock, __buffer, __size)
#define UnixSend(__sock, __buffer, __size) \
        SockWrite(__sock, __buffer, __size)

// Epoll
int EpollCreate();
Status EpollRegisterEvent(int epfd, Sock *sock);
Status EpollModifyEvent(int epfd, Sock *sock);
Status EpollDeregisterEvent(int epfd, Sock *sock);
int EpollWait(int epfd, struct epoll_event *epEvent, int timeout);

#ifdef __cplusplus
}
#endif /* __cplusplus */

#endif /* __NETWORK_H__ */
