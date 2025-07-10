---
name: Bug report
about: Creating a bug report is highly encouraged for improvement. If you are not fully convinced, it can be forwarded to the free5GC forum for further discussion.
title: "[Bugs]"
labels: ''
assignees: ''

---

**Using the provided issue template is advised; otherwise, the issue may be disregarded.**
**free5GC primarily uses GitHub for issue tracking. For general questions, information, or technical support, please consider submiting them to the [forum](https://forum.free5gc.org).**
**We appreciate it if you could check the [Troubleshooting page](https://free5gc.org/guide/Troubleshooting/) and look for duplicates in the [forum](https://forum.free5gc.org) before reporting bugs.**
<!-- Please, remove the warnings (the 3 lines above) before submitting the issue -->

## Describe the bug
A clear and concise description of what the bug is.

## To Reproduce
Steps to reproduce the behavior:
1. Change config '...'
2. Code patch '...' (You can fork the project and give us the patch diff you have modified)
3. Run command '...'
4. Dump packet '...'

## Expected behavior
A clear and concise description of what you expected to happen.

## Screenshots
If applicable, add screenshots to help explain the problem.

## Environment (please complete the following information):
 - free5GC Version: [e.g. v3.0.100]
 - OS: [e.g. Ubuntu 200.04 Server]
 - Kernel version: [e.g. 5.200.0-0-generic]
 - go version: [e.g. 1.10.0 linux/amd64]
 - c compiler version (Option): [e.g. gcc version 1.1.0]

## Trace File
### Configuration File
Provide the configuration file here.

If not clear of what to do, the `config` folder can be zip and upload it here.

### PCAP File
Dump the packet and provide the pcap file here.

If not clear of what to do, this command can be used `sudo tcpdump -i any -w free5gc.pcap` prior to running bug reproduce. Then upload the pcap file `free5gc.pcap`.

### Log File
Provide the program log file here.

If not clear of what to do, copy the printed log on the screen and upload it here.

## System architecture (Option)
Please provide the draft architecture, including the scenario, use cases, installation environment(bare metal, vm, container, or k8s), etc.

## Walkthrough (Option)
free5GC will be interested on the research or finding in brief.

## Additional context
It will be appreciated if other context can be added here.
