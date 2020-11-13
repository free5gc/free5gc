#!/bin/bash

ue_addr=${ue_addr:-"127.0.0.1"}
ue_port=${ue_port:-"10000"}
scheme=${scheme:-"https"}
auth_method=${auth_method:-"5G_AKA"}
n3iwf_address=${n3iwf_address:-"192.168.127.1"}
supi_or_suci=${supi_or_suci:-"2089300007487"}
k=${k:-"5122250214c33e723a5dd523fc145fc0"}
opc_type=${opc_type:-"OP"}
opc=${opc:-"c9e8763286b5b9ffbdf56e1297d0887b"}
ike_bind_addr=${ike_bind_addr:-"192.168.127.2"}

while [ $# -gt 0 ]; do
   if [[ $1 == *"--"* ]]; then
        param="${1/--/}"
        declare $param="$2"
   fi
  shift
done

sudo ip netns exec UEns curl --insecure --location --request POST "$scheme://$ue_addr:$ue_port/registration/" \
--header 'Content-Type: application/json' \
--data-raw "{
    \"authenticationMethod\": \"$auth_method\",
    \"supiOrSuci\": \"$supi_or_suci\",
    \"K\": \"$k\",
    \"opcType\": \"$opc_type\",
	\"opc\": \"$opc\",
	\"plmnId\": \"\",
	\"servingNetworkName\": \"\",
    \"n3IWFIpAddress\": \"$n3iwf_address\",
    \"ikeBindAddress\": \"$ike_bind_addr\",
    \"SNssai\": {
        \"Sst\": 1,
        \"Sd\": \"010203\"
    }
}"