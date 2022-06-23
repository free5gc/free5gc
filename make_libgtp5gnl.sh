#!/bin/sh
cd ..
rm -rf libgtp5gnl
git clone https://github.com/free5gc/libgtp5gnl.git
cd libgtp5gnl
autoreconf -iv
./configure --prefix=`pwd`
make
cd -