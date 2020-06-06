---
name: Bug report
about: Create a report to help us improve. If you do not ensure, you can first discuss
  in forum.
title: "[Bugs]"
labels: ''
assignees: ''

---

**If you don't follow the issue template, your issue may be closed.**
**Please note this is an issue tracker, not a support forum. For general questions, please use [forum](https://forum.free5gc.org).**
**Also, please refer to [TS](https://github.com/free5gc/free5gc/wiki/Trouble_Shooting) and [forum](https://forum.free5gc.org) before bug reporting.**
<!-- Remove above line after reporting the issue -->

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
If applicable, add screenshots to help explain your problem.

## Environment (please complete the following information):
 - free5GC Version [e.g. v3.0.2]
 - OS: [e.g. Ubuntu 18.04 Server]
 - Kernel version [e.g. 5.0.0-23-generic]
- go version [e.g. 1.12.9 linux/amd64]
 - c compiler version (Option) [e.g. gcc version 7.5.0]

## Trace File
### Configuration File
Provide your configuration here.

If you don't know what to do, you can zip the `config` folder and upload the zip file here.

### PCAP File
Dump the packet and provide the pcap file here.

If you don't know what to do, you can use `sudo tcpdump -i any -w free5gc.pcap` before you run the bug reproduce. After that, upload the free5gc.pcap file here.

### Log File
Provide the program log file here.

If you do not know what to do, copy the log print on the screen and put it here.

## System architecture (Option)
Give the draft of your architecture. What's your scenario, bare metal use case or using k8s, etc.

## Walkthrough (Option)
What you have find or research for short.

## Additional context
Add any other context about the problem here.
