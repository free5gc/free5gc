#!/bin/bash

# use some fixed parameters for testing (milenage test set 19)
cd ../../../
go build -o ./bin/ue -x ./src/ue/ue.go
sudo ip netns exec UEns ./bin/ue  #--ping=1 #\
#                                  --plmnid=010203 \
#                                  --supi=2089300007487 \
#                                  --op=c9e8763286b5b9ffbdf56e1297d0887b  \
#                                  --k=5122250214c33e723a5dd523fc145fc0 \
#                                  --opc=981d464c7c52eb6e5036234984ad0bcf \
#                                  --n3iwfip=192.168.127.1 \
#                                  --ueip=192.168.127.2 \
#                                  --mongourl=127.0.0.1:27017 \
#                                  --createue