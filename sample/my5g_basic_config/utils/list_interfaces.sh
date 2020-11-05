#!/bin/bash

# Arguments: $1: Interface ('grep'-regexp).
# Print the interfaces and their types

# Static list of types (from `ip link help`):
TYPES=(bond bond_slave bridge dummy gre gretap ifb ip6gre ip6gretap ip6tnl ipip ipoib ipvlan macvlan macvtap nlmon sit vcan veth vlan vti vxlan tun tap)

iface="$1"

for type in "${TYPES[@]}"; do
  ip link show type "${type}" | grep -E '^[0-9]+:' | cut -d ':' -f 2 | sed 's|^[[:space:]]*||' | while read _if; do
    echo "${_if}:${type}"
  done | grep "^${iface}"
done