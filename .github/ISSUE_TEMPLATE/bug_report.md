---
name: Bug report
about: Creating a bug report is highly encouraged for improvement. If you are not fully convinced, it can be forwarded to free5GC forum for further discussion.
title: "[Bugs]"
labels: ''
assignees: ''

---

**Using the provided issue template is advised; otherwise, the issue may be disregarded.**
**free5GC primarily uses GitHub for issue tracking. For general questions, information, or technical support, please consider submiting them to the [forum](https://forum.free5gc.org).**
**We appreciate it if you could check the [Troubleshooting page](https://free5gc.org/guide/Troubleshooting/) and look for duplicates in the [forum](https://forum.free5gc.org) before reporting bugs.**
<!-- Please, remove the warnings (the 3 lines above) before submitting the issue -->

## Bug Decription
A clear and concise description of what the bug is.

## To Reproduce
Steps to reproduce the behavior:
1. Change config '...'
2. Code patch '...' (You may fork the project and link the modified patch diff here)
3. Run command '...'
4. Dump packet '...'

## Expected Behavior
A clear and concise description of what you expected to happen.

## Screenshots
If applicable, add screenshots to help explain the problem.

## Environment
**Please, complete the following information:**
 - free5GC Version: [e.g. v4.0.1]
 - OS: [e.g. Ubuntu 20.04 Server]
 - Kernel version: [e.g. 5.15.0-0-generic]
 - go version: [e.g. 1.21.8 linux/amd64]
 - C compiler version (optional): [e.g. gcc version 1.1.0]

## Trace Files
### Configuration Files
Provide the configuration files.

If you are unsure of what to do, please zip free5gc's `config` folder and upload it here.

### PCAP File
Dump the relevant packets and provide the PCAP file.

If you are unsure of what to do, use the following command before reproducing the bug: `sudo tcpdump -i any -w free5gc.pcap`. Then, please upload the `free5gc.pcap` file here.

### Log File
Provide the program log file.

If you are unsure of what to do, please copy the printed log from the screen and upload it here.

## System Architecture (Optional)
Please provide the draft architecture, including the scenario, use cases, installation environment (e.g. bare metal, VM, container, or k8s), etc.

## Walkthrough (Optional)
free5GC is interested on a brief description of the research conducted or its findings.

## Additional Context
It would be appreciated if you could provide any relevant additional context here.
