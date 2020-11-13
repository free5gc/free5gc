#!/usr/bin/env bash
#
# Script not ready yet
#

CFG_DIR="$(pwd)/.."
PID_LIST=()

dirs=$(echo $CFG_DIR | tr "//" "\n" | tac)
arr=($dirs)
CORE_CFG_PATH="${arr[4]}/${arr[3]}/${arr[2]}/my5G-core.conf"

#echo $CORE_CFG_PATH
#for d in $dirs
#do
#    echo "> [$d]"
#done
#
#exit 0

cd ../../../src/upf/build
mv -f ./config/upfcfg.yaml ./config/upfcfg.yaml.old
cp ../config/upfcfg.example.my5Gcore-basic-config.yaml ./config/upfcfg.yaml
sudo -v
sudo -E ip netns exec UPFns ./bin/free5gc-upfd " -free5gccfg $CORE_CFG_PATH" &
PID_LIST+=($!)

sleep 1

cd ../../..


mongo free5gc --eval "db.NfProfile.drop()"
mongo free5gc --eval "db.urilist.drop()"

./bin/webconsole -free5gccfg "$CORE_CFG_PATH" -webuicfg "$CFG_DIR/webuicfg.conf" &
PID_LIST+=($!)

NF_LIST="nrf amf smf udr pcf udm nssf ausf"

export GIN_MODE=release

for NF in ${NF_LIST}; do
    ./bin/${NF} -free5gccfg "$CORE_CFG_PATH" $(echo " -${NF}cfg $CFG_DIR/${NF}cfg.conf") &
    PID_LIST+=($!)
done

sudo ./bin/n3iwf -free5gccfg "$CORE_CFG_PATH" -n3iwfcfg "$CFG_DIR/n3iwfcfg.conf" &
SUDO_N3IWF_PID=$!
sleep 1
N3IWF_PID=$(pgrep -P $SUDO_N3IWF_PID)
PID_LIST+=($SUDO_N3IWF_PID $N3IWF_PID)

function terminate()
{
    # kill amf first
    while $(sudo kill -SIGINT ${PID_LIST[2]} 2>/dev/null); do
        sleep 2
    done

    for ((idx=${#PID_LIST[@]}-1;idx>=0;idx--)); do
        sudo kill -SIGKILL ${PID_LIST[$idx]}
    done
}

trap terminate SIGINT
wait ${PID_LIST}
