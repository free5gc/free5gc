#!/usr/bin/env bash

LOG_PATH="./log/"
SSLKEY_LOG_FOLDER="sslkey"
NF_LOG_FOLDER="nf"
LIB_LOG_FOLDER="lib"
LOG_NAME="free5gc.log"
TODAY=$(date +"%Y%m%d_%H%M%S")
PCAP_MODE=0
N3IWF_ENABLE=0

PID_LIST=()

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

LOG_PATH_TIME=${LOG_PATH%/}"/"${TODAY}"/"
echo "log path: $LOG_PATH"

if [ ! -d ${LOG_PATH_TIME} ]; then
    mkdir -p ${LOG_PATH_TIME}
fi

if [ $PCAP_MODE -ne 0 ]; then
    PCAP=${LOG_PATH_TIME}free5gc.pcap
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

    PID_LIST+=($!)
    sleep 0.1
fi

sudo -E ./NFs/upf/build/bin/free5gc-upfd -c ./config/upfcfg.yaml -l ${LOG_PATH}${NF_LOG_FOLDER}/upf.log -g ${LOG_PATH}${LOG_NAME} &
PID_LIST+=($!)

sleep 1

NF_LIST="nrf amf smf udr pcf udm nssf ausf"

export GIN_MODE=release

for NF in ${NF_LIST}; do
    ./bin/${NF} &
    PID_LIST+=($!)
    sleep 0.1
done

if [ $N3IWF_ENABLE -ne 0 ]; then
    N3IWF_IKE_BIND_ADDRESS="127.0.0.21"
    N3IWF_UE_ADDR="10.0.0.1/24"
    sudo ip link add name ipsec0 type vti local ${N3IWF_IKE_BIND_ADDRESS} remote 0.0.0.0 key 5
    sudo ip addr add ${N3IWF_UE_ADDR} dev ipsec0
    sudo ip link set ipsec0 up
    sleep 1

    sudo ./bin/n3iwf &
    SUDO_N3IWF_PID=$!
    sleep 1
    N3IWF_PID=$(pgrep -P $SUDO_N3IWF_PID)
    PID_LIST+=($SUDO_N3IWF_PID $N3IWF_PID)
fi

function terminate()
{
    moveLog

    if [ $N3IWF_ENABLE -ne 0 ]; then
        sudo ip xfrm state > ${LOG_PATH_TIME}NWu_SA_state.log
        sudo ip xfrm state flush
        sudo ip xfrm policy flush
        sudo ip link del ipsec0
    fi

    sudo kill -SIGTERM ${PID_LIST[${#PID_LIST[@]}-2]} ${PID_LIST[${#PID_LIST[@]}-1]}
    sleep 2
}

function moveLog()
{
    SSLKEY_LOG_PATH=${LOG_PATH_TIME}${SSLKEY_LOG_FOLDER}
    if [ ! -d ${SSLKEY_LOG_PATH} ]; then
        mkdir -p ${SSLKEY_LOG_PATH}
    fi
    for KEYLOG in $(ls *sslkey.log); do
        mv ${KEYLOG} ${SSLKEY_LOG_PATH}
    done

    mv -f ${LOG_PATH}${NF_LOG_FOLDER}  ${LOG_PATH_TIME}
    mv -f ${LOG_PATH}${LIB_LOG_FOLDER}  ${LOG_PATH_TIME}
    mv ${LOG_PATH}${LOG_NAME}  ${LOG_PATH_TIME}

    sudo kill -SIGTERM ${PID_LIST[${#PID_LIST[@]}-2]} ${PID_LIST[${#PID_LIST[@]}-1]}

    sleep 2
}

trap terminate SIGINT
wait ${PID_LIST}
