#!/usr/bin/env bash

set -ex

osmo-clean-workspace.sh

autoreconf --install --force
./configure --enable-sanitize CFLAGS="-Werror" CPPFLAGS="-Werror"
$MAKE $PARALLEL_MAKE
$MAKE distcheck
$MAKE maintainer-clean

osmo-clean-workspace.sh
