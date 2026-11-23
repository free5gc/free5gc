# Ansible

## Prerequisites: `ansible`

```bash
sudo apt update
sudo apt install -y software-properties-common
sudo add-apt-repository --yes --update ppa:ansible/ansible
sudo apt install -y ansible
```

## Kernel module prerequisites

The playbook installs **gtp5g v0.10.3**, matching `quick-setup.sh`. Use a compatible kernel and install its matching headers before running the playbook; `/lib/modules/$(uname -r)/build` must exist.

Unmodified gtp5g v0.10.3 fails to build on Ubuntu kernel **7.0.0-43-generic** because of a UDP tunnel API change. The native `quick-setup.sh` flow was validated on Ubuntu **26.04.1** with kernel **6.18.64-061864-generic**, gtp5g **v0.10.3**, and MongoDB **9.0.8**. That validation covers native NF startup and UE-to-DN traffic; the full Ansible/Helm deployment on Ubuntu 26.04 has not been validated.
