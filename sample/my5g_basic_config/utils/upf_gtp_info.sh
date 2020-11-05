#!/bin/bash

watch -d -n 1 ip netns exec UPFns $(pwd)/../../../src/upf/lib/libgtp5gnl/tools/gtp5g-tunnel list $1

#ip netns exec UPFns $(pwd)/../src/upf/lib/libgtp5gnl/tools/gtp5g-tunnel list far