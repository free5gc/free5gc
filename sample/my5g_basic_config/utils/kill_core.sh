#!/bin/bash


sudo ps -eaf  | grep "free5gc" | grep -v grep | awk '{ print $2 }' | sudo xargs kill -9