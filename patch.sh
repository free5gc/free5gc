NF_LIST="nrf amf smf udr pcf udm nssf ausf upf n3iwf"

git checkout .

for NF in ${NF_LIST}; do
    cd NFs/$NF
    git checkout .
    PATCH=("${NF^^}_PATCH")
    if [ "${!PATCH}" != "" ];
    then
        echo $NF
        git checkout main
        curl ${!PATCH} | git apply
        git status
    fi
    cd ../..
done
