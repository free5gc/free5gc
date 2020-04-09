#!/bin/bash
SCRIPT_DIR=`pwd`/script

Usage() {
    echo "usage: $0 < SimpleUPTest | ULCLTest1 | RANSetup | UPFSetup | I-UPFSetup | A-UPFSetup | Clean >"
}

if [ $# -ne 1 ]; then
    Usage
    exit -1
fi

# Check OS
if [ -f /etc/os-release ]; then
    # freedesktop.org and systemd
    . /etc/os-release
    OS=$NAME
    VER=$VERSION_ID
else
    # Fall back to uname, e.g. "Linux <version>", also works for BSD, etc.
    OS=$(uname -s)
    VER=$(uname -r)
    echo "This Linux version is too old: $OS:$VER, we don't support!"
    exit 1
fi

sudo -v
if [ $? == 1 ]
then
    echo "Without root permission, you cannot run the test due to our test is using namespace"
    exit 1
fi

while getopts 'o' OPT;
do
    case $OPT in
        o) DUMP_NS=True;;
    esac
done
shift $(($OPTIND - 1))

cd ${SCRIPT_DIR}
if [ $1 == "SimpleUPTest" ]; then
    ./ns_ran_upf.sh
elif [ $1 == "ULCLTest1" ]; then
    ./ns_ran1_iupf1_aupf1.sh
elif [ $1 == "RANSetup" ]; then
    ./ran.sh
elif [ $1 == "UPFSetup" ]; then
    ./upf.sh
elif [ $1 == "I-UPFSetup" ]; then
    ./iupf.sh
elif [ $1 == "A-UPFSetup" ]; then
    ./aupf.sh
elif [ $1 == "Clean" ]; then
    ./cleanup.sh
else
    Usage
    exit -1
fi

echo "Done."
