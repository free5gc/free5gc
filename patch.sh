NF_LIST="nrf amf smf udr pcf udm nssf ausf upf n3iwf"

for NF in ${NF_LIST}; do
    cd NFs/$NF
    git checkout . # remove uncommited changes
    git checkout main # switch to main branch
    PATCH=("${NF^^}_PATCH")
    if [ "${!PATCH}" != "" ];
    then
        echo $NF # print NF name
        curl ${!PATCH} | git apply # apply patch
        git status
    fi
    cd ../..
done
