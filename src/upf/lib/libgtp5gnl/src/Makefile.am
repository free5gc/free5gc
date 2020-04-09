include $(top_srcdir)/Make_global.am

# This is _NOT_ the library release version, it's an API version.
# Please read Chapter 6 "Library interface versions" of the libtool documentation before making any modification
LIBVERSION=1:1:1


lib_LTLIBRARIES = libgtp5gnl.la

noinst_HEADERS = internal.h tools.h

libgtp5gnl_la_LIBADD = ${LIBMNL_LIBS}
libgtp5gnl_la_LDFLAGS = -Wl,--version-script=$(srcdir)/libgtp5gnl.map -version-info $(LIBVERSION)
libgtp5gnl_la_SOURCES = genl.c		\
		      gtp5g-genl-pdr.c  \
		      gtp5g-genl-far.c  \
		      gtp5g-rtnl.c	\
		      gtp5g.c		\
		      libgtp5gnl.map
