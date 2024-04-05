#!/usr/bin/env bash

# checks if no parameter was given as input
if [ -z "$1" ]
then
    echo "[ERRO] No parameter was given!"
    echo "[INFO] You must provide dn interface name as input parameter"
    echo "Usage:"
    echo "$0 <dn_interface>"
    echo "Example:"
    echo "$0 enp0s4"
else
    # if the parameter is present, cache root credentials
    sudo -v
    if [ $? == 1 ]
    then
        echo "[ERRO] Without root permission, you cannot change iptables configuration"
        exit 1
    fi
    # if user has root permissions, then start to modify the rules
    echo "[INFO] Using $1 as interface name"

    # first, delete any previous applied rules
    echo -n "[INFO] Removing all old iptables rules, if any... "
    sudo iptables -P INPUT ACCEPT
    sudo iptables -P FORWARD ACCEPT
    sudo iptables -P OUTPUT ACCEPT
    sudo iptables -t nat -F
    sudo iptables -t mangle -F
    sudo iptables -F
    sudo iptables -X
    echo "[OK]"

    # then apply the new ones
    echo -n "[INFO] Applying iptables rules... "
    sudo iptables -t nat -A POSTROUTING -o $1 -j MASQUERADE
    sudo iptables -I FORWARD 1 -j ACCEPT
    echo "[OK]"
    echo -n "[INFO] Setting kernel net.ipv4.ip_forward flag... "
    sudo sysctl -w net.ipv4.ip_forward=1 >/dev/null
    echo "[OK]"
    echo -n "[INFO] Stopping ufw firewall... "
    sudo systemctl stop ufw
    echo "[OK]"

    echo "[INFO] Configuration applied successfully"
fi
