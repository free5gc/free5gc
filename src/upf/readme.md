# free5GC-UPF

## Get Started
### Prerequisites
```bash
sudo apt-get -y update
sudo apt-get -y install git gcc cmake go libmnl-dev autoconf libtool libyaml-dev
go get github.com/sirupsen/logrus
```

### Build
```bash
mkdir build
cd build
cmake ..
make -j`nproc`
```

### Test
```bash
cd build/bin
./testutlt
sudo ./testgtpv1
```

### Edit configuration file
After building from sources, edit `./build/config/upfcfg.yaml`

### Setup environment
```bash
sh -c 'echo 1 > /proc/sys/net/ipv4/ip_forward'
iptables -t nat -A POSTROUTING -o {DN_Interface_Name} -j MASQUERADE
```

### Run
```bash
cd build
sudo ./bin/free5gc-upfd
```
To show usage: `./bin/free5gc-upfd -h`


## Clean the Environment (if needed)
### Remove POSIX message queues
```bash
ls /dev/mqueue/
rm /dev/mqueue/*
```

### Remove gtp devices (using tools in libgtpnl)
```bash
cd lib/libgtpnl-1.2.1/tools
sudo ./gtp-link del {Dev-Name}
```
