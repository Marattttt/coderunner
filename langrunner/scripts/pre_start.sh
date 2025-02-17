#!/bin/bash

set -e 
iptables -A OUTPUT -d 0.0.0.0/0 -j DROP
