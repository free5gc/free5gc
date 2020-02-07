---
name: Bug report
about: Creating bug report is highly welcome for improvement. If not fully convince, it can be forward to free5GC forum for further discussion.
title: "[Bugs]"
labels: ''
assignees: ''

---

**We will advise people to follow the issue template set, otherwise the issue might be disregarded.**
**free5GC mainly uses GitHub for issue tracking. Information regarding to general questions or technical support. It will be highly considered if forward to the [forum](https://forum.free5gc.org).**
**free5GC will appreciate it, if people can refer to [TS](https://github.com/free5gc/free5gc/wiki/Trouble_Shooting) and [forum](https://forum.free5gc.org) prior to bug reporting**
<!-- Remove warning (above 3 lines) while reporting the issue -->

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
