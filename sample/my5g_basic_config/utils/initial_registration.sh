#!/bin/bash

UE_ADDR="127.0.0.1"
UE_PORT="10000"
SCHEME="https"
AUTH_METHOD="5G_AKA"
N3IWF_ADDR="192.168.127.1"
SUPI_OR_SUCI="2089300007487"
K="5122250214c33e723a5dd523fc145fc0"
OPC_TYPE="OP"
OPC="c9e8763286b5b9ffbdf56e1297d0887b"
IKE_BIND_ADDR="192.168.127.2"

sudo ip netns exec UEns curl --insecure --location --request POST "$SCHEME://$UE_ADDR:$UE_PORT/registration/" \
--header 'Content-Type: application/json' \
--data-raw "{
    \"authenticationMethod\": \"$AUTH_METHOD\",
    \"supiOrSuci\": \"$SUPI_OR_SUCI\",
    \"K\": \"$K\",
    \"opcType\": \"$OPC_TYPE\",
	\"opc\": \"$OPC\",
	\"plmnId\": \"\",
	\"servingNetworkName\": \"\",
    \"n3IWFIpAddress\": \"$N3IWF_ADDR\",
    \"ikeBindAddress\": \"$IKE_BIND_ADDR\",
    \"SNssai\": {
        \"Sst\": 1,
        \"Sd\": \"010203\"
    }
}"