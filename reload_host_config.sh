#!/usr/bin/env bash

INTERFACE='' # to record the dn interface name

# checks if no parameter was given as input
if [ -z "$1" ]
then
    echo "[ERRO] No parameter was given!"
    echo "[INFO] You must provide dn interface name as input parameter"
    echo "Usage:"
    echo "$0 <dn_interface>"
    echo "Example:"
    echo "$0 enp0s4"
# then check if more than two parameters were input
elif [ $# -gt 2 ]
then
    echo "[ERRO] Too many parameters!"
    echo "Usage:"
    echo "$0 <dn_interface> [-reset-firewall]"
    echo "Examples:"
    echo "$0 enp0s4"
    echo "$0 enp0s4 -reset-firewall"
else
    # if any two parameters are present, cache root credentials
    sudo -v
    if [ $? == 1 ]
    then
        echo "[ERRO] Without root permission, you cannot change iptables configuration"
        exit 1
    fi
    # if user has root permissions, then start to modify the rules

    # check if the user wants to reset the iptables firewall rules
    if [ "$1" = "-reset-firewall" ] || [ "$2" = "-reset-firewall" ]
    then
        # warn the user before deleting the rules
        echo "[WARN] Firewall reset is enabled"
        echo "[WARN] ALL iptables rules will be DELETED"
        read -p "Press ENTER to continue or Ctrl+C to abort now"

        # If yes, then delete any previous applied rules
        echo -n "[INFO] Removing all old iptables rules, if any... "
        sudo iptables -P INPUT ACCEPT
        sudo iptables -P FORWARD ACCEPT
        sudo iptables -P OUTPUT ACCEPT
        sudo iptables -t nat -F
        sudo iptables -t mangle -F
        sudo iptables -F
        sudo iptables -X
        echo "[OK]"
        # adjust the vars to accept parameters in any order
        if [ "$1" = "-reset-firewall" ]
        then
            # check if the interface parameter is present
            # (the interface name is a requirement)
            if [ -z "$2" ]
            then
                echo "[INFO] You must provide dn interface name as input parameter"
                echo "Example:"
                echo "$0 -reset-firewall enp0s4"
                exit 1
            fi
            INTERFACE=$2
        elif [ "$2" = "-reset-firewall" ]
        then
            INTERFACE=$1
        else
            echo "[ERRO] Could not set the interface name, please check your input"
        fi
    fi

    # if there isn't any other parameter, then just set the interface name
    if [ -z "$2" ]
    then
        INTERFACE=$1
    fi

    # check if interface name was correctly set
    if [ -z "${INTERFACE}" ]
    then
        echo "[ERRO] Could not set the interface name, please check your input"
        exit 2
    fi

    echo "[INFO] Using $INTERFACE as interface name"

    # then apply the new iptables firewall rules
    echo -n "[INFO] Applying iptables rules... "
    sudo iptables -t nat -A POSTROUTING -o $INTERFACE -j MASQUERADE
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
