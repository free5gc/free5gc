include $(top_srcdir)/Make_global.am

noinst_PROGRAMS = gtp5g-link		\
		  gtp5g-tunnel

gtp5g_link_SOURCES = gtp5g-link.c
gtp5g_link_LDADD = ../src/libgtp5gnl.la ${LIBMNL_LIBS}

gtp5g_tunnel_SOURCES = gtp5g-tunnel.c
gtp5g_tunnel_LDADD = ../src/libgtp5gnl.la ${LIBMNL_LIBS}
