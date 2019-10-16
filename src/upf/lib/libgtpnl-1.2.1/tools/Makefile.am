include $(top_srcdir)/Make_global.am

noinst_PROGRAMS = gtp-link		\
		  gtp-tunnel

gtp_link_SOURCES = gtp-link.c
gtp_link_LDADD = ../src/libgtpnl.la ${LIBMNL_LIBS}

gtp_tunnel_SOURCES = gtp-tunnel.c
gtp_tunnel_LDADD = ../src/libgtpnl.la ${LIBMNL_LIBS}
