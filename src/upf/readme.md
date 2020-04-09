# free5GC-UPF

## Get Started
### Prerequisites
Libraries used in UPF
```bash
sudo apt-get -y update
sudo apt-get -y install git gcc cmake go libmnl-dev autoconf libtool libyaml-dev
go get github.com/sirupsen/logrus
```

Linux kernel module 5G GTP-U (Linux kernel version = 5.0.0-23-generic)
```bash
git clone https://github.com/PrinzOwO/gtp5g.git
cd gtp5g
make
sudo make install
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

### Remove gtp devices (using tools in libgtp5gnl)
```bash
cd lib/libgtp5gnl/tools
sudo ./gtp5g-link del {Dev-Name}
```
