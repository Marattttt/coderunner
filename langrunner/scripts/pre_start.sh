#!/bin/bash

set -e 
iptables -A OUTPUT -j DROP
