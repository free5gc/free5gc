include $(top_srcdir)/Make_global.am

@RELMAKE@

ACLOCAL_AMFLAGS = -I m4

SUBDIRS = src include tools
DIST_SUBDIRS = src include tools

pkgconfigdir = $(libdir)/pkgconfig
pkgconfig_DATA = libgtp5gnl.pc

${pkgconfig_DATA}: ${top_builddir}/config.status

@RELMAKE@
