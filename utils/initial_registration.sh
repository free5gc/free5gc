#!/bin/bash

sudo ip netns exec UEns curl --insecure --location --request POST 'https://127.0.0.1:10000/registration/' \
--header 'Content-Type: application/json' \
--data-raw '{
    "authenticationMethod": "5G_AKA",
    "supiOrSuci": "2089300007487",
    "K": "5122250214c33e723a5dd523fc145fc0",
    "opcType": "OP",
	"opc": "c9e8763286b5b9ffbdf56e1297d0887b",
	"plmnId": "",
	"servingNetworkName": "",
    "n3IWFIpAddress": "192.168.127.1",
    "SNssai": {
        "Sst": 1,
        "Sd": "010203"
    }
}'