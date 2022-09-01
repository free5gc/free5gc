NF_LIST="nrf amf smf udr pcf udm nssf ausf upf n3iwf"

for NF in ${NF_LIST}; do
    cd NFs/$NF
    git checkout main
    PATCH=("${NF^^}_PATCH")
    if [ "${!PATCH}" != "" ];
    then
        curl ${!PATCH} | git apply
        git status
    fi
    cd ../..
done
