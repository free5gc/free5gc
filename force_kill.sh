#!/bin/bash

DROP_ALL_DB=0
DB_NAME="free5gc"
DB_DROP_COLLECTION=(
    "NfProfile"
    "applicationData.influenceData.subsToNotify"
    "applicationData.subsToNotify"
    "policyData.subsToNotify"
    "exposureData.subsToNotify"
)

if [ $# -ne 0 ]; then
    while [ $# -gt 0 ]; do
        case $1 in
            -db)
                DROP_ALL_DB=1
                ;;
        esac
        shift
    done
fi

NF_LIST="nrf amf smf udr pcf udm nssf ausf n3iwf upf chf tngf"

for NF in ${NF_LIST}; do
    sudo killall -9 ${NF}
done

sudo killall tcpdump
sudo ip link del upfgtp
sudo ip link del ipsec0
XFRMI_LIST=($(ip link | grep xfrmi | awk -F'[:,@]' '{print $2}'))
for XFRMI_IF in "${XFRMI_LIST[@]}"
do
    sudo ip link del $XFRMI_IF
done
sudo rm /dev/mqueue/*
sudo rm -f /tmp/free5gc_unix_sock
sudo rm -f cert/*_*
sudo rm -f test/cert/*_*
sudo rm -f /tmp/config.json # CHF ChargingGatway FTP config

if [ $DROP_ALL_DB -eq 1 ]; then
    mongo --eval "db.dropDatabase()" "$DB_NAME"
    mongosh --eval "db.dropDatabase()" "$DB_NAME"
else
    MONGO_SCRIPT=""
    for COLLECTION in "${DB_DROP_COLLECTION[@]}"
    do
        MONGO_SCRIPT+="db.$COLLECTION.drop();"
    done
    mongo "$DB_NAME" --eval "$MONGO_SCRIPT"
    mongosh "$DB_NAME" --eval "$MONGO_SCRIPT" 
fi
