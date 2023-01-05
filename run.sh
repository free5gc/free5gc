#!/usr/bin/env bash

LOG_PATH="./log/"
LOG_NAME="free5gc.log"
TODAY=$(date +"%Y%m%d_%H%M%S")
PCAP_MODE=0
N3IWF_ENABLE=0

PID_LIST=()
echo $$ > run.pid

if [ $# -ne 0 ]; then
    while [ $# -gt 0 ]; do
        case $1 in
            -p)
                shift
                case $1 in
                    -*)
                        continue ;;
                    *)
                        if [ "$1" != "" ];
                        then
                            LOG_PATH=$1
                        fi
                esac ;;
            -cp)
                PCAP_MODE=$((${PCAP_MODE} | 0x01))
                ;;
            -dp)
                PCAP_MODE=$((${PCAP_MODE} | 0x02))
                ;;
            -n3iwf)
                N3IWF_ENABLE=1
                ;;
        esac
        shift
    done
fi

function terminate()
{
    rm run.pid
    echo "Receive SIGINT, terminating..."
    if [ $N3IWF_ENABLE -ne 0 ]; then
        sudo ip xfrm state > ${LOG_PATH}NWu_SA_state.log
        sudo ip xfrm state flush
        sudo ip xfrm policy flush
        sudo ip link del ipsec0
        XFRMI_LIST=($(ip link | grep xfrmi | awk -F'[:,@]' '{print $2}'))
        for XFRMI_IF in "${XFRMI_LIST[@]}"
        do
            sudo ip link del $XFRMI_IF
        done
    fi
    
    for ((i=${#PID_LIST[@]}-1;i>=0;i--)); do
        sudo kill -SIGTERM ${PID_LIST[i]}
    done
    sleep 2
    wait ${PID_LIST}
    exit 0
}

trap terminate SIGINT

LOG_PATH=${LOG_PATH%/}"/"${TODAY}"/"
echo "log path: $LOG_PATH"

if [ ! -d ${LOG_PATH} ]; then
    mkdir -p ${LOG_PATH}
fi

if [ $PCAP_MODE -ne 0 ]; then
    PCAP=${LOG_PATH}free5gc.pcap
    case $PCAP_MODE in
        1)  # -cp
            if [ $N3IWF_ENABLE -ne 0 ]; then
                sudo tcpdump -i any 'sctp port 38412 || tcp port 8000 || udp port 8805 || udp port 500 || udp port 4500' -w ${PCAP} &
            else
                sudo tcpdump -i any 'sctp port 38412 || tcp port 8000 || udp port 8805' -w ${PCAP} &
            fi
            ;;
        2)  # -dp
            if [ $N3IWF_ENABLE -ne 0 ]; then
                sudo tcpdump -i any 'udp port 2152 || ip proto 50' -w ${PCAP} &
            else
                sudo tcpdump -i any 'udp port 2152' -w ${PCAP} &
            fi
            ;;
        3)  # include -cp -dp
            if [ $N3IWF_ENABLE -ne 0 ]; then
                sudo tcpdump -i any 'sctp port 38412 || tcp port 8000 || udp port 8805 || udp port 500 || udp port 4500 || udp port 2152 || ip proto 50' -w ${PCAP} &
            else
                sudo tcpdump -i any 'sctp port 38412 || tcp port 8000 || udp port 8805 || udp port 2152' -w ${PCAP} &
            fi
            ;;
    esac

    SUDO_TCPDUMP_PID=$!
    sleep 0.1
    TCPDUMP_PID=$(pgrep -P $SUDO_TCPDUMP_PID)
    PID_LIST+=($SUDO_TCPDUMP_PID $TCPDUMP_PID)
fi

sudo -E ./bin/upf -c ./config/upfcfg.yaml -l ${LOG_PATH}upf.log -lc ${LOG_PATH}${LOG_NAME} &
SUDO_UPF_PID=$!
sleep 0.1
UPF_PID=$(pgrep -P $SUDO_UPF_PID)
PID_LIST+=($SUDO_UPF_PID $UPF_PID)

mongo --eval "db.NfProfile.drop();db.applicationData.influenceData.subsToNotify.drop();db.applicationData.subsToNotify.drop();db.policyData.subsToNotify.drop();db.exposureData.subsToNotify.drop()" free5gc
mongosh --eval "db.NfProfile.drop();db.applicationData.influenceData.subsToNotify.drop();db.applicationData.subsToNotify.drop();db.policyData.subsToNotify.drop();db.exposureData.subsToNotify.drop()" free5gc

sleep 0.1

NF_LIST="nrf amf smf udr pcf udm nssf ausf"

export GIN_MODE=release

for NF in ${NF_LIST}; do
    ./bin/${NF} -c ./config/${NF}cfg.yaml -l ${LOG_PATH}${NF}.log -lc ${LOG_PATH}${LOG_NAME} &
    PID_LIST+=($!)
    sleep 0.1
done

if [ $N3IWF_ENABLE -ne 0 ]; then
    sudo ./bin/n3iwf -c ./config/n3iwfcfg.yaml -l ${LOG_PATH}n3iwf.log -lc ${LOG_PATH}${LOG_NAME} &
    SUDO_N3IWF_PID=$!
    sleep 1
    N3IWF_PID=$(pgrep -P $SUDO_N3IWF_PID)
    PID_LIST+=($SUDO_N3IWF_PID $N3IWF_PID)
fi

wait ${PID_LIST}
exit 0
