include $(top_srcdir)/Make_global.am

# This is _NOT_ the library release version, it's an API version.
# Please read Chapter 6 "Library interface versions" of the libtool documentation before making any modification
LIBVERSION=1:1:1


lib_LTLIBRARIES = libgtpnl.la

noinst_HEADERS = internal.h

libgtpnl_la_LIBADD = ${LIBMNL_LIBS}
libgtpnl_la_LDFLAGS = -Wl,--version-script=$(srcdir)/libgtpnl.map -version-info $(LIBVERSION)
libgtpnl_la_SOURCES = genl.c		\
		      gtp-genl.c	\
		      gtp-rtnl.c	\
		      gtp.c		\
		      libgtpnl.map
