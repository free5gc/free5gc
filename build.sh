#!/usr/bin/env bash

NF_LIST="amf ausf nrf nssf pcf smf udm udr n3iwf"

for NF in ${NF_LIST}; do 
    echo "Start build ${NF}...."
    go build -o bin/${NF} -x src/${NF}/${NF}.go
done

echo "Start build UPF...."

cd src/upf # cwd is upf
rm -rf build
mkdir -p build
cd ./build
cmake ..
make -j$(nproc)
