# free5GC

## Hardware Tested
There are no gNB and UE for standalone 5GC available in the market yet.

## Minimum Requirement
- Software
    - OS: Ubuntu 18.04 or later versions
    - gcc 7.3.0
    - Go 1.12.9 linux/amd64
    - QEMU emulator 2.11.1
    - kernel version 5.0.0-23-generic (MUST for UPF)
    
**Note: Please use Ubuntu 18.04 or later versions and go 1.12.9 linux/amd64** 


You can use `go version` to check your current Go version.
```bash
- Hardware
    - CPU: Intel i5 processor
    - RAM: 4GB
    - Hard drive: 160G
    - NIC card: 1Gbps ethernet card

- Hardware recommended
    - CPU: Intel i7 processor
    - RAM: 8GB
    - Hard drive: 160G
    - NIC card: 10Gbps ethernet card
```


## Installation
### A. Pre-requisite

0. Required kernel version `5.0.0-23-generic`. This request is from the module
   [gtp5g](https://github.com/PrinzOwO/gtp5g) that we has used. Any more details
   please check [here](https://github.com/PrinzOwO/gtp5g)
   ```bash
   # Check kernel version
   $ uname -r
   5.0.0-23-generic
   ```

1. Require go language
    * If another version of Go is installed
        - Please remove the previous Go version
            - ```sudo rm -rf /usr/local/go```
        - Install Go 1.12.9
            ```bash
            wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
            sudo tar -C /usr/local -zxvf go1.12.9.linux-amd64.tar.gz
            ```
    * Clean installation
        - Install Go 1.12.9
             ```bash
            wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
            sudo tar -C /usr/local -zxvf go1.12.9.linux-amd64.tar.gz
            mkdir -p ~/go/{bin,pkg,src}
            echo 'export GOPATH=$HOME/go' >> ~/.bashrc
            echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
            echo 'export PATH=$PATH:$GOPATH/bin:$GOROOT/bin' >> ~/.bashrc
            echo 'export GO111MODULE=off' >> ~/.bashrc
            source ~/.bashrc
            ```

2. Required packages for control plane
    ```bash
    sudo apt -y update
    sudo apt -y install mongodb wget git
    sudo systemctl start mongodb
    ```

3. Required packages for user plane
    ```bash
    sudo apt -y update
    sudo apt -y install git gcc cmake autoconf libtool pkg-config libmnl-dev libyaml-dev
    go get -u github.com/sirupsen/logrus
    ```

4. Network Setting
    ```bash
    sudo sysctl -w net.ipv4.ip_forward=1
    sudo iptables -t nat -A POSTROUTING -o ${DN_INTERFACE} -j MASQUERADE
    ```

### B. Install Control Plane Entities
    
1. Clone free5GC project in `$GOPATH/src`
    ```bash
    cd $GOPATH/src
    git clone https://github.com/free5gc/free5gc.git
    cd free5gc
    git submodule update --init
    ```

    **In step 2, the folder name should remain free5gc. Please do not modify it or the compilation would fail.**
2. Run the script to install dependent packages
    ```bash
    cd $GOPATH/src/free5gc
    chmod +x ./install_env.sh
    ./install_env.sh
    
    Please ignore error messages during the package dependencies installation process.
    ```

3. Compile network function services in `$GOPATH/src/free5gc` individually, e.g. AMF (redo this step for each NF), or
    ```bash
    cd $GOPATH/src/free5gc
    go build -o bin/amf -x src/amf/amf.go
    ```
    **To build all network functions in one command**
    ```bash
    ./build.sh
    ```


### C. Install User Plane Function (UPF)
    
1. Please check Linux kernel version if it is `5.0.0-23-generic`
    ```bash
    uname -r
    ```


    Get Linux kernel module 5G GTP-U
    ```bash
    git clone https://github.com/PrinzOwO/gtp5g.git
    cd gtp5g
    make
    sudo make install
    ```
2. Build from sources
    ```bash
    cd $GOPATH/src/free5gc/src/upf
    mkdir build
    cd build
    cmake ..
    make -j`nproc`
    ```
    
**Note: Config is located at** `$GOPATH/src/free5gc/src/upf/build/config/upfcfg.yaml
   `

## Run

### A. Run Core Network 
Option 1. Run network function service individually, e.g. AMF (redo this for each NF), or
```bash
cd $GOPATH/src/free5gc
./bin/amf
```

**Note: For N3IWF needs specific configuration in section B** 

Option 2. Run whole core network with command
```
./run.sh
```

### B. Run N3IWF (Individually)
To run N3IWF, make sure the machine is equipped with three network interfaces. (one is for connecting AMF, another is for connecting UPF, the other is for IKE daemon)

We need to configure each interface with a suitable IP address.

We have to create an interface for IPSec traffic:
```bash
# replace <...> to suitable value
sudo ip link add ipsec0 type vti local <IKEBindAddress> remote 0.0.0.0 key <IPSecInterfaceMark>
```
Assign an address to this interface, then bring it up:
```bash
# replace <...> to suitable value
sudo ip address add <IPSecInterfaceAddress/CIDRPrefix> dev ipsec0
sudo ip link set dev ipsec0 up
```

Run N3IWF (root privilege is required):
```bash
cd $GOPATH/src/free5gc/
sudo ./bin/n3iwf
```

## Test
Start Wireshark to capture any interface with `pfcp||icmp||gtp` filter and run the tests below to simulate the procedures:
```bash
cd $GOPATH/src/free5gc
chmod +x ./test.sh
```
a. TestRegistration
```bash
(In directory: $GOPATH/src/free5gc)
./test.sh TestRegistration
```
b. TestServiceRequest
```bash
./test.sh TestServiceRequest
```
c. TestXnHandover
```bash
./test.sh TestXnHandover
```
d. TestDeregistration
```bash
./test.sh TestDeregistration
```
e. TestPDUSessionReleaseRequest
```bash
./test.sh TestPDUSessionReleaseRequest
```

f. TestPaging
```!
./test.sh TestPaging
```

g. TestN2Handover
```!
./test.sh TestN2Handover
```

h. TestNon3GPP
```bash
./test.sh TestNon3GPP
```

i. TestULCL
```bash
./test_ulcl.sh -om 3 TestRegistration
```

*For more details, you can reference to our [wiki](https://github.com/free5gc/free5gc/wiki)*

## Release Note
### v3.0.2
+ refactor:
    + (all) refactor coding style for NFs written in golang, including folder name, package name, and file name
    + (upf) refactor the method of UPF initialization
    + (all) merge NF config models to one file each
    + (openapi) move openapi clients to repo "openapi"
+ feature:
    + (all) logger support output to file
    + (openapi) Add Convert to convert interface
    + (all) support h2c mode for SBI, use h2c as default mode
    + (openapi) add serialize deserialize function
+ bugfix:
    + (smf) fix duplicated pdu session handling in SMF
    + (nrf) fix subscribe time decode issue
    + (udm) op and opc decision rule
    + (smf) sm context release occur panic
    + (smf) fix A-UPF use NodeID not UPIP in DL
    + (amf) add ie nil check when handling handover request acknowledge
    + (amf) add missing ie sourceToTargetTransparentContaier to ngap message handoverRequest
    + (smf) fix ulcl workaround in release v3.0.0
    + (milenage) fix f1 function bug
    + (amf) fix generate Kamf P0 parameter

### v3.0.1
+ project:
    + Change the way we manage project. Using git submodule to manage hole
      project to let each NF and library has its own version control
    + Open source our library
+ Add document of SMF ULCL limitation
+ hotfix:
    + fix NRF return nil error (issue#12)
    + fix OPc crash error (issue#21)
    + update webconsole pakcage version to prevent security issue
    + SMF fix pdu session release procedure
    + fix pdu session release procedure test
+ SMF:
    + SMF support NF deregistration

### v3.0.0
+ AMF
    + Support SMF selection at PDU session establishment
    + Fix SUCI handling procedure
+ SMF
    + Feature
        + ULCL by config
        + Authorized QoS
    + Bugfix
        + PDU Session Establishment PDUAddress Information
        + PDU Session Establishment N1 Message
        + SMContext Release Procedure
+ UPF:
    + ULCL feature
    + support SDF Filter
    + support N9 interface
+ OAM
    + Get Registered UE Context API
    + OAM web UI to display Registered UE Context
+ N3IWF
    + Support Registration procedure for untrusted non-3GPP access
    + Support UE Requested PDU Session Establishment via Untrusted non-3GPP Access
+ UDM
    + SUCI to SUPI de-concealment
    + Notification 
        + Callback notification to NF ( in SDM service)
        + UDM initiated deregistration notification to NF ( in UECM service)

### v2.0.2
+ Add debug mode on NFs
+ Auto add Linux routing when UPF runs
+ Add AMF consumer for AM policy
+ Add SM policy
+ Allow security NIA0 and NEA0
+ Add handover feature
+ Add webui
+ Update license
+ Bugfix for incorrect DNN
+ Bugfix for NFs registering to NRF

### v2.0.1
+ Global
    + Update license and readme
    + Add Paging feature
    + Bugfix for AN release issue
    + Add URL for SBI in NFs' config

+ AMF
    + Add Paging feature
    + Bugfix for SCTP PPID to 60
    + Bugfix for UE release in testing
    + Bugfix for too fast send UP data in testing
    + Bugfix for sync with defaultc config in testing

+ SMF
    + Add Paging feature
    + Create PDR with FAR ID
    + Bugfix for selecting DNN fail handler

+ UPF
    + Sync config default address with Go NFs
    + Remove GTP tunnel by removing PDR/FAR
    + Bugfix for PFCP association setup
    + Bugfix for new PDR/FAR creating
    + Bugfix for PFCP session report
    + Bugfix for getting from PDR
    + Bugfix for log format and update logger version

+ PCF
    + Bugfix for lost field and method

