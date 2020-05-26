#!/usr/bin/env bash

RELEASE_FOLDER=release
export GOPATH=`go env GOPATH`
export GO111MODULE="off"
if [[ ! `pwd` == */${RELEASE_FOLDER} ]]; then cd ./${RELEASE_FOLDER}; fi
RELEASE_PATH=`pwd`
SOURCE_PATH=${RELEASE_PATH}/..

usage="$(basename "$0") [-h] [-m] [-n] -- Release script for free5GC

where:
    -h  show help message
    -w  weekly release, which will clone the project from membership repository
    -m  for membership release
    -n  for open source release
        without any parameter, regarded as open source release"

while getopts 'hnmw' OPT;
do
    case ${OPT} in
        m)
            # Membership release
            RELEASE_TYPE=membership
            rm -rf ${RELEASE_PATH}/src
            mkdir -p ${RELEASE_PATH}/src/free5gc
            cd ${RELEASE_PATH}/src/free5gc
            ;;
        w)
            # Weelky release
            RELEASE_TYPE=membership
            rm -rf ${RELEASE_PATH}/src # remove old release
            git clone git@bitbucket.org:free5gc_membership/free5gc_member.git ${RELEASE_PATH}/src/free5gc
            rc=$?
            if [ $rc -ne 0 ]; then echo $rc; exit $rc; fi
            cd ${RELEASE_PATH}/src/free5gc
            ls -a -I .git -I . -I .. | xargs rm -rf # clean all files
            ;;
        n)
            # Open Source release
            RELEASE_TYPE=open_source
            rm -rf ${RELEASE_PATH}/src
            mkdir -p ${RELEASE_PATH}/src/free5gc
            cd ${RELEASE_PATH}/src/free5gc
            ;;
        h)
            echo "${usage}"
            exit
            ;;
        \?)
            printf "Illegal option: -%s\n" "${OPTARG}" >&2
            echo "${usage}" >&2
            exit 1
            ;;
    esac
done

if [[ ${OPTIND} -eq 1 ]];
then
    # patch release
    RELEASE_TYPE=open_source
    rm -rf ${RELEASE_PATH}/src
    mkdir -p ${RELEASE_PATH}/src/free5gc
    cd ${RELEASE_PATH}/src/free5gc
fi

rsync -r --exclude '.git' --exclude 'release*' --exclude 'bitbucket-pipelines.yml' --exclude '.golangci.yml' --exclude 'infra' ${SOURCE_PATH}/ .
# under src/free5gc
ROOT_PATH=`pwd`

# UPF
go get -u github.com/sirupsen/logrus

cd ${ROOT_PATH}/src/upf/
rm -rf build # if build remind
./upf_release.sh
cd ${ROOT_PATH}

# Remove asn1
rm -rf ${ROOT_PATH}/asn1c-s1ap

# golang env
ENV_SCRIPT=${ROOT_PATH}/install_env.sh
echo '#!/bin/bash
set -o xtrace

# Install go packages which are not required to switch version' > ${ENV_SCRIPT}
cat go.mod | grep -P ' v\d[(\.\d)]*' | grep -P ' v\d[(\.\d)]*-' | awk '{ print "go get -u "$1 }' >> ${ENV_SCRIPT}
echo '
# Install go packages which are required to switch version' >> ${ENV_SCRIPT}
cat go.mod | grep -P ' v\d[(\.\d)]*' | grep -v -P ' v\d[(\.\d)]*-' | sed -e 's/\+.*//g' | awk '{ print "go get -u "$1; print "cd $GOPATH/src/"$1; print "git checkout "$2}' >> ${ENV_SCRIPT}
rm -f ${ROOT_PATH}/go.mod ${ROOT_PATH}/go.sum

chmod +x ${ENV_SCRIPT}
${ENV_SCRIPT}

# .go
perl -i -pe 'BEGIN{undef $/;} s/"gofree5gc/"free5gc/g' `find ${ROOT_PATH}/ -name '*.go'`
# .conf
perl -i -pe 'BEGIN{undef $/;} s/gofree5gc/free5gc/g' `find ${ROOT_PATH}/ -name '*.conf'`

# Remove support folder
find ${ROOT_PATH} -type d -name support -mindepth 2 -exec rm -r {} +

# Build go package
GOS=`go env GOHOSTOS`
GOARCH=`go env GOARCH`
LIB_SRC=${ROOT_PATH}/lib
LIB_DST=${ROOT_PATH}/pkg/${GOS}_${GOARCH}/free5gc/lib
NON_STATIC=('CommonConsumerTestData/AMF/TestAmf' 'CommonConsumerTestData/AMF/TestComm' 'CommonConsumerTestData/AUSF/TestUEAuth' 'CommonConsumerTestData/NSSF/TestNSSelection' 'CommonConsumerTestData/PCF/TestAMPolicy' 'CommonConsumerTestData/PCF/TestBDTPolicy' 'CommonConsumerTestData/PCF/TestPolicyAuthorization' 'CommonConsumerTestData/PCF/TestSMPolicy' 'CommonConsumerTestData/SMF/TestPDUSession' 'CommonConsumerTestData/UDM/TestGenAuthData' 'CommonConsumerTestData/UDR/TestRegistrationProcedure' 'Namf_Communication' 'Namf_EventExposure' 'Namf_Location' 'Namf_MT' 'Nausf_SoRProtection' 'Nausf_UEAuthentication' 'Nausf_UPUProtection' 'Nnrf_AccessToken' 'Nnrf_NFDiscovery' 'Nnrf_NFManagement' 'Nnssf_NSSAIAvailability' 'Nnssf_NSSelection' 'Npcf_AMPolicy' 'Npcf_BDTPolicyControl' 'Npcf_PolicyAuthorization' 'Npcf_SMPolicyControl' 'Npcf_UEPolicy' 'Nsmf_PDUSession' 'Nudm_EventExposure' 'Nudm_ParameterProvision' 'Nudm_SubscriberDataManagement' 'Nudm_UEAuthentication' 'Nudm_UEContextManagement' 'Nudr_DataRepository' 'openapi' 'openapi/common' 'openapi/models' 'util_3gpp')
cd ${LIB_SRC}
ALL_LIB=($(go list ./... 2>/dev/null | sed 's/.*free5gc\/lib\///g'))
cd ${ROOT_PATH}
STATIC=()

for ((i = 0; i < ${#ALL_LIB[@]}; i++))
do
    if [[ ! ${NON_STATIC[@]} =~ ${ALL_LIB[$i]} && ! ${ALL_LIB[$i]} =~ ^.*logger ]]; then
        STATIC+=("${ALL_LIB[$i]}")
    fi
done

mkdir -p ${LIB_DST}
export GOPATH="${GOPATH}:${ROOT_PATH}/../.."
for ((i = 0; i < ${#STATIC[@]}; i++))
do
    cd ${LIB_SRC}/${STATIC[$i]}
    go build -o ${LIB_DST}`pwd | perl -pe 's/.*lib//g'`.a -x
done

cd ${ROOT_PATH}
tar zcvf free5gc_libs.tar.gz pkg
rm -rf ${ROOT_PATH}/pkg

# Remove source
for ((i = 0; i < ${#STATIC[@]}; i++))
do
    if [[ ${STATIC[$i]} == *'/'* ]]; then
        continue
    fi
    cd ${LIB_SRC}/${STATIC[$i]}
    perl -i -pe 'BEGIN{undef $/;} s/func.*\K\{[\s\S]*?\n\}/\{\}/g' `find . -path ./logger -prune -o -name '*.go' -print`
    perl -i -pe 'BEGIN{undef $/;} s/package /\/\/go:binary\-only\-package\n\npackage /g' `find . -path ./logger -prune -o -name '*.go' -print`
done
cd ${ROOT_PATH}

if [[ ${RELEASE_TYPE} == 'open_source' ]]; then
    find ${ROOT_PATH} -type f -name *_debug.sh | xargs rm -f
    find ${ROOT_PATH} -type f -name *_debug.go | xargs rm -f
    perl -i -pe 'BEGIN{undef $/;} s/\/\/[ \t]*\+build[ \t]+!debug[\s]+//g' `find ${ROOT_PATH}/ -name '*.go'`
    #find ${ROOT_PATH} -type f -name 'logger.go' | xargs sed -i '/+build/,/package/{//!d};/+build/d'
elif [[ ${RELEASE_TYPE} == 'membership' ]]; then
    find ${ROOT_PATH} -type f -name 'logger_debug.go' | sed 's/_debug//' | xargs rm -f
    perl -i -pe 'BEGIN{undef $/;} s/\/\/[ \t]*\+build[ \t]+debug[\s]+//g' `find ${ROOT_PATH}/ -name 'logger_debug.go'`
    #find ${ROOT_PATH} -type f -name 'logger_debug.go' | xargs sed -i '/+build/,/package/{//!d};/+build/d'
fi

perl -i -pe 'BEGIN{undef $/;} s/GO111MODULE=on/GO111MODULE=off/' ${ROOT_PATH}/test*.sh

