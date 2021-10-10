#!/bin/bash

yum install -y git
git clone https://github.com/TOMOFUMI-KONDO/compare-http.git /home/ec2-user/compare-http
chown -R ec2-user:ec2-user /home/ec2-user/compare-http

# Raise UDP receive buffer size for QUIC
# https://github.com/lucas-clemente/quic-go/wiki/UDP-Receive-Buffer-Size
sysctl -w net.core.rmem_max=2500000
