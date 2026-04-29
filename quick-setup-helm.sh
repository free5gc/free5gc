#!/bin/bash

ANSIBLE_PATH="ansible-helm"

if ! command -v ansible >/dev/null 2>&1; then
    echo "Ansible is not installed. Please install Ansible and try again. Use:"
    echo " sudo apt update"
    echo " sudo apt install -y software-properties-common"
    echo " sudo add-apt-repository --yes --update ppa:ansible/ansible"
    echo " sudo apt install -y ansible"
    exit 1
fi

read -p "Enter SSH user: " USER
ansible-playbook -i $ANSIBLE_PATH/inventory.ini $ANSIBLE_PATH/install.yaml -u "$USER" --ask-pass --ask-become-pass